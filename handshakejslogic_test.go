package handshakejslogic_test

import (
	//"fmt"
	handshakejslogic "github.com/handshakejs/handshakejslogic"
	"github.com/stvp/tempredis"
	"testing"
)

const (
	EMAIL    = "app0@mailinator.com"
	APP_NAME = "app0"
	SALT     = "1234"
)

func TestAppsCreate(t *testing.T) {
	config := tempredis.Config{
		"port":      "11001",
		"databases": "1",
	}

	tempredis.Temp(config, func(err error) {
		if err != nil {
			panic(err)
		}

		redis_url := "127.0.0.1:11001"
		handshakejslogic.Setup(redis_url)

		//app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME}

		//result, err := handshakejslogic.AppsCreate(app)
		//if err != nil {
		//	t.Errorf("Error", err)
		//}
		//if result["email"] != EMAIL {
		//	t.Errorf("Incorrect email " + result["email"].(string))
		//}
		//if result["app_name"] != APP_NAME {
		//	t.Errorf("Incorrect app_name " + result["app_name"].(string))
		//}
		//if result["salt"] == nil {
		//	t.Errorf("Salt is nil and should not be.")
		//}
	})
}

//func TestAppsCreateCustomSalt(t *testing.T) {
//	app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": SALT}
//
//	result, err := handshakejslogic.AppsCreate(app)
//	if err != nil {
//		t.Errorf("Error")
//	}
//
//	if result["salt"] != SALT {
//		t.Errorf("Salt did not equal " + SALT)
//	}
//}
//
//func TestAppsCreateCustomBlankSalt(t *testing.T) {
//	app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": ""}
//
//	result, err := handshakejslogic.AppsCreate(app)
//	if err != nil {
//		t.Errorf("Error")
//	}
//
//	if result["salt"] == "" {
//		t.Errorf("It should generate a salt if blank.")
//	}
//}
