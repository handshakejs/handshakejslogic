package handshakejslogic

import (
	"errors"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/scottmotte/redisurlparser"
	"log"
	"strings"
)

var (
	conn redis.Conn
)

type App struct {
	AppName string
	Email   string
	Salt    string
}

type LogicError struct {
	Code    string
	Field   string
	Message string
}

func Setup(redis_url string) {
	ru, err := redisurlparser.Parse(redis_url)
	if err != nil {
		log.Fatal(err)
	}

	conn, err = redis.Dial("tcp", ru.Host+":"+ru.Port)
	if err != nil {
		log.Fatal(err)
	}

	if ru.Password != "" {
		if _, err := conn.Do("AUTH", ru.Password); err != nil {
			conn.Close()
			log.Fatal(err)
		}
	}
}

func AppsCreate(app map[string]interface{}) (map[string]interface{}, *LogicError) {
	var app_name string
	if str, ok := app["app_name"].(string); ok {
		app_name = strings.Replace(str, " ", "", -1)
	} else {
		app_name = ""
	}
	if app_name == "" {
		logic_error := &LogicError{"required", "app_name", "app_name cannot be blank"}
		return app, logic_error
	}
	app["app_name"] = app_name

	generated_salt := uniuri.NewLen(20)
	if app["salt"] == nil {
		app["salt"] = generated_salt
	}
	if app["salt"].(string) == "" {
		app["salt"] = generated_salt
	}

	key := "apps/" + app["app_name"].(string)
	err := checkIfAppExists(key)
	if err != nil {
		logic_error := &LogicError{"not_unique", "app_name", "app_name must be unique"}
		return app, logic_error
	}
	err = addAppToApps(conn, app["app_name"].(string))
	if err != nil {
		logic_error := &LogicError{"unkown", "", err.Error()}
		return nil, logic_error
	}
	err = addAppToKey(conn, key, app)
	if err != nil {
		logic_error := &LogicError{"unkown", "", err.Error()}
		return nil, logic_error
	}

	return app, nil
}

func addAppToKey(conn redis.Conn, key string, app map[string]interface{}) error {
	args := []interface{}{key}
	for k, v := range app {
		args = append(args, k, v)
	}
	_, err := conn.Do("HMSET", args...)
	if err != nil {
		return err
	}

	return nil
}

func addAppToApps(conn redis.Conn, app_name string) error {
	_, err := conn.Do("SADD", "apps", app_name)
	if err != nil {
		return err
	}

	return nil
}

func checkIfAppExists(key string) error {
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		return err
	}
	if res.(int64) == 1 {
		err = errors.New("That app_name already exists.")
		return err
	}

	return nil
}
