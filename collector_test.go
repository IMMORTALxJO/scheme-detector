package schemedetector

import (
	"testing"

	"gotest.tools/assert"
)

func testDetect(t *testing.T, keys map[string]string, answers []string) {
	result := []string{}
	for _, s := range FromMap(keys) {
		result = append(result, s.String())
	}
	for _, a := range answers {
		assert.Check(t, stringInArray(a, result))
	}
}

func TestCollector(t *testing.T) {
	testDetect(t, map[string]string{
		"API_ENDPOINT": "https://example.com/test",
	}, []string{"https://example.com:443/test"})
	testDetect(t, map[string]string{
		"API_ENDPOINT": "https://user@example.com/test",
	}, []string{"https://user@example.com:443/test"})
	testDetect(t, map[string]string{
		"DB": "postgres://user:password@127.0.0.1:5432/test",
	}, []string{"postgres://user:password@127.0.0.1:5432/test"})
	testDetect(t, map[string]string{
		"DB": "postgres://user@127.0.0.1:5432/test",
	}, []string{"postgres://user@127.0.0.1:5432/test"})
	testDetect(t, map[string]string{
		"DATABASE_HOST": "127.0.0.1",
		"DATABASE_USER": "user",
		"DATABASE_PASS": "password",
		"DATABASE_PORT": "5432",
		"DATABASE_NAME": "test",
	}, []string{"postgres://user:password@127.0.0.1:5432/test"})
	testDetect(t, map[string]string{
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
	}, []string{
		"postgres://user:password@127.0.0.1:5432/test",
		"postgres://user2:password2@127.0.0.2:5432/test2",
		"https://api_user:api_password@example.com:443/test",
	})

}
