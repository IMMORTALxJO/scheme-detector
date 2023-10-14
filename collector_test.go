package schemedetector

import (
	"os"
	"sort"
	"testing"

	"gotest.tools/assert"
)

func testDetect(t *testing.T, result []*Scheme, answers []string) {
	var resultStrings []string

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
		"DB": "pgsql://user:password@127.0.0.1:5432/test",
	}), []string{"pgsql://user:password@127.0.0.1:5432/test"})
	testDetect(t, FromMap(map[string]string{
		"DB": "pgsql://user@127.0.0.1:5432/test",
	}), []string{"pgsql://user@127.0.0.1:5432/test"})
	testDetect(t, FromMap(map[string]string{
		"DB": "pgsql://user@127.0.0.1:5432/test?option1=value1&option2=value2",
	}), []string{"pgsql://user@127.0.0.1:5432/test?option1=value1&option2=value2"})
	testDetect(t, FromMap(map[string]string{
		"DATABASE_HOST": "127.0.0.1",
		"DATABASE_USER": "user",
		"DATABASE_PASS": "password",
		"DATABASE_PORT": "5432",
		"DATABASE_NAME": "test",
	}), []string{"pgsql://user:password@127.0.0.1:5432/test"})
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
		"pgsql://user:password@127.0.0.1:5432/test",
		"pgsql://user2:password2@127.0.0.2:5432/test2",
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
		"DB_ADDRESS": "pgsql://user:password@127.0.0.1:5432/test",
		"DB_PORT":    "1234",
		"DB_HOST":    "localhost",
		"DB_NAME":    "db",
		"DB_USER":    "custom",
		"DB_PASS":    "newpassword",
	}), []string{"pgsql://user:password@127.0.0.1:5432/test"})

}

func TestCollectorFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("API_URL", "https://postman:password@postman-echo.com/basic-auth")
	os.Setenv("DB_ENDPOINT", "pgsql://user@127.0.0.1:5432/test?option1=value1&option2=value2")
	testDetect(t, FromEnv(), []string{
		"https://postman:password@postman-echo.com:443/basic-auth",
		"pgsql://user@127.0.0.1:5432/test?option1=value1&option2=value2",
	})

}

func TestCollectorHelpers(t *testing.T) {
	assert.Check(t, FromMap(map[string]string{
		"DB": "pgsql://user@127.0.0.1:5432/test",
	})[0].IsIP())
	assert.Check(t, !FromMap(map[string]string{
		"DB": "pgsql://user@127.0.0.1:5432/test",
	})[0].IsDNSName())
	assert.Check(t, FromMap(map[string]string{
		"DB": "pgsql://user@localhost:5432/test",
	})[0].IsDNSName())
	assert.Check(t, !FromMap(map[string]string{
		"DB": "pgsql://user@localhost:5432/test",
	})[0].IsIP())

}

func TestComplicated(t *testing.T) {
	os.Clearenv()
	testDetect(t, FromMap(map[string]string{
		"EXAMPLE_API_ENDPOINT":     "http://api.example.com",
		"EXAMPLE_API_PASSWORD":     "123pass",
		"EXAMPLE_API_USERNAME":     "api_username",
		"DATABASE_MAIN_HOST":       "base-master.default.svc.cluster.local",
		"DATABASE_MAIN_NAME":       "api",
		"DATABASE_MAIN_PASSWORD":   "MAIN_base_PASS",
		"DATABASE_MAIN_PORT":       "6432",
		"DATABASE_MAIN_SCHEMA":     "pgsql",
		"DATABASE_MAIN_USERNAME":   "api_user_main",
		"DATABASE_REPLICA_URL":     "pgsql://api_user_main:MAIN_base_PASS@base-repl.default.svc.cluster.local:5432/api",
		"DATABASE_SECOND_HOST":     "second-base.default.svc.cluster.local",
		"DATABASE_SECOND_PASSWORD": "SECOND_base_PASS",
		"DATABASE_SECOND_PORT":     "3306",
		"DATABASE_SECOND_SCHEMA":   "mysql",
		"DATABASE_SECOND_USERNAME": "api_user_second",
		"RABBIT_MQ_HOST":           "rabbitmq.default.svc.cluster.local",
		"RABBIT_MQ_PASSWORD":       "rabbitmq_pass",
		"RABBIT_MQ_USER":           "rabbitmq_user",
		"RABBIT_MQ_PORT":           "5672",
	}), []string{
		"amqp://rabbitmq_user:rabbitmq_pass@rabbitmq.default.svc.cluster.local:5672",
		"http://api_username:123pass@api.example.com:80",
		"mysql://api_user_second:SECOND_base_PASS@second-base.default.svc.cluster.local:3306",
		"pgsql://api_user_main:MAIN_base_PASS@base-master.default.svc.cluster.local:6432/api",
		"pgsql://api_user_main:MAIN_base_PASS@base-repl.default.svc.cluster.local:5432/api",
	})
}

func TestFilter(t *testing.T) {
	os.Clearenv()
	os.Setenv("SCHEME_DETECTOR_EXCLUDE", ".*PASS")
	testDetect(t, FromMap(map[string]string{
		"DATABASE_HOST": "127.0.0.1",
		"DATABASE_USER": "user",
		"DATABASE_PASS": "password",
		"DATABASE_PORT": "5432",
		"DATABASE_NAME": "test",
	}), []string{"pgsql://user@127.0.0.1:5432/test"})
}
