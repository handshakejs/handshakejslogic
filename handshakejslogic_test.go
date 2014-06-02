package handshakejslogic_test

import (
	"../handshakejslogic"
	"github.com/stvp/tempredis"
	"testing"
)

const (
	APP_NAME  = "app0"
	EMAIL     = "app0@mailinator.com"
	SALT      = "1234"
	REDIS_URL = "redis://127.0.0.1:11001"
)

func tempredisConfig() tempredis.Config {
	config := tempredis.Config{
		"port":      "11001",
		"databases": "1",
	}
	return config
}

func TestAppsCreate(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			panic(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME}

		handshakejslogic.Setup(REDIS_URL)
		result, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["email"] != EMAIL {
			t.Errorf("Incorrect email " + result["email"].(string))
		}
		if result["app_name"] != APP_NAME {
			t.Errorf("Incorrect app_name " + result["app_name"].(string))
		}
		if result["salt"] == nil {
			t.Errorf("Salt is nil and should not be.")
		}
	})
}

func TestAppsCreateCustomSalt(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			panic(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": SALT}

		handshakejslogic.Setup(REDIS_URL)
		result, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["salt"] != SALT {
			t.Errorf("Salt did not equal " + SALT)
		}
	})
}

func TestAppsCreateCustomBlankSalt(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			panic(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": ""}

		handshakejslogic.Setup(REDIS_URL)
		result, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["salt"] == nil || result["salt"].(string) == "" {
			t.Errorf("It should generate a salt if blank.")
		}
	})
}

//func TestAppsCreateBlankAppName(t *testing.T) {
//	tempredis.Temp(tempredisConfig(), func(err error) {
//		if err != nil {
//			panic(err)
//		}
//
//		app := map[string]interface{}{"email": EMAIL, "app_name": ""}
//
//		handshakejslogic.Setup(REDIS_URL)
//		_, logic_error := handshakejslogic.AppsCreate(app)
//		if logic_error.Code != "required" {
//			t.Errorf("Error", err)
//		}
//	})
//}

func TestAppsCreateNilAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			panic(err)
		}

		app := map[string]interface{}{"email": EMAIL}

		handshakejslogic.Setup(REDIS_URL)
		_, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestAppsCreateSpacedAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			panic(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": " "}

		handshakejslogic.Setup(REDIS_URL)
		_, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestAppsCreateAppNameWithSpaces(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			panic(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": "combine these"}

		handshakejslogic.Setup(REDIS_URL)
		result, _ := handshakejslogic.AppsCreate(app)
		if result["app_name"] != "combinethese" {
			t.Errorf("Incorrect combining of app_name " + result["app_name"].(string))
		}
	})
}
