package handshakejslogic_test

import (
	"../handshakejslogic"
	"github.com/stvp/tempredis"
	"testing"
)

const (
	EMAIL     = "app0@mailinator.com"
	APP_NAME  = "app0"
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
		result, err := handshakejslogic.AppsCreate(app)
		if err != nil {
			t.Errorf("Error", err)
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
		result, err := handshakejslogic.AppsCreate(app)
		if err != nil {
			t.Errorf("Error", err)
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
		result, err := handshakejslogic.AppsCreate(app)
		if err != nil {
			t.Errorf("Error", err)
		}

		if result["salt"] == "" {
			t.Errorf("It should generate a salt if blank.")
		}
	})
}
