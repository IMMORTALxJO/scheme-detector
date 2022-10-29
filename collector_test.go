package schemedetector

import (
	"os"
	"sort"
	"testing"

	"gotest.tools/assert"
)

func testDetect(t *testing.T, result []*Scheme, answers []string) {
	resultStrings := []string{}
	for _, r := range result {
		resultStrings = append(resultStrings, r.String())
	}
	sort.Strings(resultStrings)
	sort.Strings(answers)
	assert.DeepEqual(t, resultStrings, answers)
}

func TestCollectorMap(t *testing.T) {
	testDetect(t, FromMap(map[string]string{
		"API_ENDPOINT": "https://example.com/test",
	}), []string{"https://example.com:443/test"})
	testDetect(t, FromMap(map[string]string{
		"API_ENDPOINT": "https://user@example.com/test",
	}), []string{"https://user@example.com:443/test"})
	testDetect(t, FromMap(map[string]string{
		"DB": "postgres://user:password@127.0.0.1:5432/test",
	}), []string{"postgres://user:password@127.0.0.1:5432/test"})
	testDetect(t, FromMap(map[string]string{
		"DB": "postgres://user@127.0.0.1:5432/test",
	}), []string{"postgres://user@127.0.0.1:5432/test"})
	testDetect(t, FromMap(map[string]string{
		"DB": "postgres://user@127.0.0.1:5432/test?option1=value1&option2=value2",
	}), []string{"postgres://user@127.0.0.1:5432/test?option1=value1&option2=value2"})
	testDetect(t, FromMap(map[string]string{
		"DATABASE_HOST": "127.0.0.1",
		"DATABASE_USER": "user",
		"DATABASE_PASS": "password",
		"DATABASE_PORT": "5432",
		"DATABASE_NAME": "test",
	}), []string{"postgres://user:password@127.0.0.1:5432/test"})
	testDetect(t, FromMap(map[string]string{
		"DATABASE_HOST":  "127.0.0.1",
		"DATABASE_USER":  "user",
		"DATABASE_PASS":  "password",
		"DATABASE_PORT":  "5432",
		"DATABASE_NAME":  "test",
		"DATABASE2_HOST": "127.0.0.2",
		"DATABASE2_USER": "user2",
		"DATABASE2_PASS": "password2",
		"DATABASE2_PORT": "5432",
		"DATABASE2_NAME": "test2",
		"API_ENDPOINT":   "https://example.com/test",
		"API_USER":       "api_user",
		"API_PASS":       "api_password",
	}), []string{
		"postgres://user:password@127.0.0.1:5432/test",
		"postgres://user2:password2@127.0.0.2:5432/test2",
		"https://api_user:api_password@example.com:443/test",
	})
	// Unknown port => unknown schema => no result
	testDetect(t, FromMap(map[string]string{
		"DATABASE_HOST": "127.0.0.1",
		"DATABASE_USER": "user",
		"DATABASE_PASS": "password",
		"DATABASE_PORT": "0",
		"DATABASE_NAME": "test",
	}), []string{})
	// Unknown schema => unknown port => no result
	testDetect(t, FromMap(map[string]string{
		"ENDPOINT_URL": "unknown://127.0.0.1/test",
	}), []string{})

	// Overwrite is not allowed
	testDetect(t, FromMap(map[string]string{
		"DB_ADDRESS": "postgres://user:password@127.0.0.1:5432/test",
		"DB_PORT":    "1234",
		"DB_HOST":    "localhost",
		"DB_NAME":    "db",
		"DB_USER":    "custom",
		"DB_PASS":    "newpassword",
	}), []string{"postgres://user:password@127.0.0.1:5432/test"})

}

func TestCollectorFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("API_URL", "https://postman:password@postman-echo.com/basic-auth")
	os.Setenv("DB_ENDPOINT", "postgres://user@127.0.0.1:5432/test?option1=value1&option2=value2")
	testDetect(t, FromEnv(), []string{
		"https://postman:password@postman-echo.com:443/basic-auth",
		"postgres://user@127.0.0.1:5432/test?option1=value1&option2=value2",
	})

}

func TestCollectorHelpers(t *testing.T) {
	assert.Check(t, FromMap(map[string]string{
		"DB": "postgres://user@127.0.0.1:5432/test",
	})[0].IsIP())
	assert.Check(t, !FromMap(map[string]string{
		"DB": "postgres://user@127.0.0.1:5432/test",
	})[0].IsDNSName())
	assert.Check(t, FromMap(map[string]string{
		"DB": "postgres://user@localhost:5432/test",
	})[0].IsDNSName())
	assert.Check(t, !FromMap(map[string]string{
		"DB": "postgres://user@localhost:5432/test",
	})[0].IsIP())

}
