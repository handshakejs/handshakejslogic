package handshakejslogic

import (
	"bytes"
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/scottmotte/redisurlparser"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_10                     = "10"
	AUTHCODE_LIFE_IN_MS_DEFAULT = "120000"
	AUTHCODE_LENGTH_DEFAULT     = "4"
)

var (
	conn                redis.Conn
	AUTHCODE_LIFE_IN_MS string
	AUTHCODE_LENGTH     string
)

type Options struct {
	AuthcodeLifeInMs string
	AuthcodeLength   string
}

type LogicError struct {
	Code    string
	Field   string
	Message string
}

func Setup(redis_url string, options *Options) {
	if options.AuthcodeLifeInMs == "" {
		AUTHCODE_LIFE_IN_MS = AUTHCODE_LIFE_IN_MS_DEFAULT
	} else {
		AUTHCODE_LIFE_IN_MS = options.AuthcodeLifeInMs
	}
	if options.AuthcodeLength == "" {
		AUTHCODE_LENGTH = AUTHCODE_LENGTH_DEFAULT
	} else {
		AUTHCODE_LENGTH = options.AuthcodeLength
	}

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
	err := validateAppDoesNotExist(key)
	if err != nil {
		logic_error := &LogicError{"not_unique", "app_name", "app_name must be unique"}
		return app, logic_error
	}
	err = addAppToApps(conn, app["app_name"].(string))
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}
	err = saveApp(conn, key, app)
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return app, nil
}

func IdentitiesConfirm(identity map[string]interface{}) (map[string]interface{}, *LogicError) {
	app_name, logic_error := checkAppNamePresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["app_name"] = app_name

	email, logic_error := checkEmailPresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["email"] = email

	authcode, logic_error := checkAuthcodePresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["authcode"] = authcode

	app_name_key := "apps/" + identity["app_name"].(string)
	key := app_name_key + "/identities/" + identity["email"].(string)

	err := validateAppExists(app_name_key)
	if err != nil {
		logic_error := &LogicError{"not_found", "app_name", "app_name could not be found"}
		return identity, logic_error
	}
	err = validateIdentityExists(key)
	if err != nil {
		logic_error := &LogicError{"not_found", "email", "email could not be found"}
		return identity, logic_error
	}

	var r struct {
		Email             string `redis:"email"`
		AppName           string `redis:"app_name"`
		Authcode          string `redis:"authcode"`
		AuthcodeExpiredAt string `redis:"authcode_expired_at"`
	}
	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}
	err = redis.ScanStruct(values, &r)
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}

	email = r.Email
	res_authcode := r.Authcode
	res_authcode_expired_at := r.AuthcodeExpiredAt

	current_ms_epoch_time := (time.Now().Unix() * 1000)
	res_authcode_expired_at_int64, err := strconv.ParseInt(res_authcode_expired_at, 10, 64)
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}

	if len(res_authcode) > 0 && res_authcode == authcode {
		if res_authcode_expired_at_int64 < current_ms_epoch_time {
			logic_error := &LogicError{"expired", "authcode", "authcode has expired. request another one."}
			return identity, logic_error
		}

		app_salt, err := redis.String(conn.Do("HGET", app_name_key, "salt"))
		if err != nil {
			logic_error := &LogicError{"unknown", "", err.Error()}
			return identity, logic_error
		}
		hash := pbkdf2.Key([]byte(email), []byte(app_salt), 1000, 16, sha1.New)
		identity["hash"] = hex.EncodeToString(hash)

		return identity, nil
	} else {
		logic_error := &LogicError{"incorrect", "authcode", "the authcode was incorrect"}
		return identity, logic_error
	}
}

func IdentitiesCreate(identity map[string]interface{}) (map[string]interface{}, *LogicError) {
	app_name, logic_error := checkAppNamePresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["app_name"] = app_name

	email, logic_error := checkEmailPresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["email"] = email

	app_name_key := "apps/" + identity["app_name"].(string)
	key := app_name_key + "/identities/" + identity["email"].(string)

	err := validateAppExists(app_name_key)
	if err != nil {
		logic_error := &LogicError{"not_found", "app_name", "app_name could not be found"}
		return identity, logic_error
	}
	err = addIdentityToIdentities(conn, app_name_key, identity["email"].(string))
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}
	err = saveIdentity(conn, key, identity)
	if err != nil {
		logic_error := &LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return identity, nil
}

func saveApp(conn redis.Conn, key string, app map[string]interface{}) error {
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

func validateAppDoesNotExist(key string) error {
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

func validateAppExists(key string) error {
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		return err
	}
	if res.(int64) != 1 {
		err = errors.New("That app_name does not exist.")
		return err
	}

	return nil
}

func validateIdentityExists(key string) error {
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		return err
	}
	if res.(int64) != 1 {
		err = errors.New("That identity does not exist.")
		return err
	}

	return nil
}
func addIdentityToIdentities(conn redis.Conn, app_name_key string, email string) error {
	_, err := conn.Do("SADD", app_name_key+"/identities", email)
	if err != nil {
		return err
	}

	return nil
}

func saveIdentity(conn redis.Conn, key string, identity map[string]interface{}) error {
	base_10, err := strconv.Atoi(BASE_10)
	if err != nil {
		return err
	}

	authcode_life_in_ms, err := strconv.Atoi(AUTHCODE_LIFE_IN_MS)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	authcode, err := randomAuthCode()
	if err != nil {
		return err
	}
	identity["authcode"] = authcode
	unixtime := (time.Now().Unix() * 1000) + int64(authcode_life_in_ms)
	identity["authcode_expired_at"] = strconv.FormatInt(unixtime, base_10)

	args := []interface{}{key}
	for k, v := range identity {
		args = append(args, k, v)
	}
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		return err
	}

	return nil
}

func randomAuthCode() (string, error) {
	base_10, err := strconv.Atoi(BASE_10)
	if err != nil {
		return "", err
	}
	authcode_length, err := strconv.Atoi(AUTHCODE_LENGTH)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	var buffer bytes.Buffer

	for i := 1; i <= authcode_length; i++ {
		random_number := int64(rand.Intn(10))
		number_as_string := strconv.FormatInt(random_number, base_10)
		buffer.WriteString(number_as_string)
	}

	return buffer.String(), nil
}

func checkAppNamePresent(identity map[string]interface{}) (string, *LogicError) {
	var app_name string
	if str, ok := identity["app_name"].(string); ok {
		app_name = strings.Replace(str, " ", "", -1)
	} else {
		app_name = ""
	}
	if app_name == "" {
		logic_error := &LogicError{"required", "app_name", "app_name cannot be blank"}
		return app_name, logic_error
	}

	return app_name, nil
}

func checkEmailPresent(identity map[string]interface{}) (string, *LogicError) {
	var email string
	if str, ok := identity["email"].(string); ok {
		email = strings.Replace(str, " ", "", -1)
	} else {
		email = ""
	}
	if email == "" {
		logic_error := &LogicError{"required", "email", "email cannot be blank"}
		return email, logic_error
	}

	return email, nil
}

func checkAuthcodePresent(identity map[string]interface{}) (string, *LogicError) {
	var authcode string
	if str, ok := identity["authcode"].(string); ok {
		authcode = strings.Replace(str, " ", "", -1)
	} else {
		authcode = ""
	}
	if authcode == "" {
		logic_error := &LogicError{"required", "authcode", "authcode cannot be blank"}
		return authcode, logic_error
	}

	return authcode, nil
}

func Conn() redis.Conn {
	return conn
}
