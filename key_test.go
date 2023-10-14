package schemedetector

import (
	"testing"

	"gotest.tools/assert"
)

func TestKeyString(t *testing.T) {
	assert.Equal(t, (newKey("test", "")).String(), "test")
	assert.Equal(t, (newKey("key1_test_key2", "")).String(), "key1_test_key2")
	assert.Equal(t, (newKey("KEY1_TEST_KEY2", "")).String(), "KEY1_TEST_KEY2")
	assert.Equal(t, (newKey("", "")).String(), "")
}

func TestKeyNew(t *testing.T) {
	assert.Check(t, (newKey("test", "")).uri == nil)
	// test correct uri parsing
	assert.Check(t, (newKey("test", "http://user:pass@localhost:8080/testpath")).uri != nil)
	assert.Equal(t, (newKey("test", "http://user:pass@localhost:8080/testpath")).uri.Hostname(), "localhost")
	assert.Equal(t, (newKey("test", "http://user:pass@localhost:8080/testpath")).uri.Port(), "8080")
	assert.Equal(t, (newKey("test", "http://user:pass@localhost:8080/testpath")).uri.Path, "/testpath")
	assert.Equal(t, (newKey("test", "http://user:pass@localhost:8080/testpath")).uri.User.Username(), "user")
	// test hints detection
	assert.Check(t, (newKey("TEST_HOST", "")).hints.host)
	assert.Check(t, (newKey("TEST_PATH", "")).hints.path)
	assert.Check(t, (newKey("TEST_PORT", "")).hints.port)
	assert.Check(t, (newKey("TEST_USER", "")).hints.user)
	assert.Check(t, (newKey("TEST_PASS", "")).hints.pass)
}

func TestKeyFindSimilars(t *testing.T) {

	checkFindSimilars(t, newKey("DATABASE_MAIN_HOST", ""),
		[]string{"DATABASE_MAIN_NAME", "DATABASE_MAIN_PASSWORD", "DATABASE_MAIN_PORT", "DATABASE_MAIN_SCHEMA", "DATABASE_MAIN_USERNAME", "DATABASE_REPLICA_URL", "DATABASE_SECOND_HOST", "DATABASE_SECOND_PASSWORD", "DATABASE_SECOND_PORT", "DATABASE_SECOND_SCHEMA", "DATABASE_SECOND_USERNAME", "RABBIT_MQ_HOST"},
		[]string{"DATABASE_MAIN_NAME", "DATABASE_MAIN_PASSWORD", "DATABASE_MAIN_PORT", "DATABASE_MAIN_USERNAME"},
	)

}

func checkFindSimilars(t *testing.T, k *key, c []string, answer []string) {
	var candidates []*key

	for _, candidate := range c {
		candidates = append(candidates, newKey(candidate, ""))
	}
	var result []string
	resultKeys := k.findSimilars(candidates)
	for _, r := range resultKeys {
		result = append(result, r.name)
	}
	assert.DeepEqual(t, answer, result)
}
