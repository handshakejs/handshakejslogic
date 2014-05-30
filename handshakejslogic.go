package handshakejslogic

import (
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/scottmotte/redisurlparser"
)

var (
	conn redis.Conn
)

func Setup(redis_url string) {
	fmt.Println(redis_url)

	ru, err := redisurlparser.Parse(redis_url)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	fmt.Println(ru)

	//fmt.Println("adfjdkfjdkfjdkfjdkfjdkfjdkfjdkf")
	//fmt.Println(ru)
	//_, err = redis.Dial("tcp", ru.Host+":"+ru.Port)
	//if err != nil {
	//	panic(err)
	//}
	//if _, err := conn.Do("AUTH", ru.Password); err != nil {
	//	conn.Close()
	//	panic(err)
	//}
}

func AppsCreate(app map[string]interface{}) (map[string]interface{}, error) {

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
		return app, err
	}

	//err = addAppToApps(conn, app["app_name"].(string))
	//if err != nil {
	//	return nil, err
	//}

	//err = addAppToKey(conn, key, app)
	//if err != nil {
	//	return nil, err
	//}

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
