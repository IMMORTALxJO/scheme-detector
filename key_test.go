package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestKeyHasHints(t *testing.T) {
	assert.Check(t, (&key{"test", ""}).hasHints([]string{"test"}))
	assert.Check(t, (&key{"key1_test", ""}).hasHints([]string{"test"}))
	assert.Check(t, (&key{"test_test", ""}).hasHints([]string{"test"}))
	assert.Check(t, (&key{"test_key1", ""}).hasHints([]string{"test"}))
	assert.Check(t, (&key{"test__key1", ""}).hasHints([]string{"test"}))
	assert.Check(t, (&key{"key1_test_key2", ""}).hasHints([]string{"test"}))
	assert.Check(t, (&key{"key1_test_key2", ""}).hasHints([]string{"notfound", "test"}))
	assert.Check(t, (&key{"key1_test_key2", ""}).hasHints([]string{"test", "notfound"}))
	assert.Check(t, (&key{"key1_test_key2", ""}).hasHints([]string{"key1", "key2"}))
	assert.Check(t, (&key{"key1_test_key2", ""}).hasHints([]string{"key1", "test", "key2"}))
	assert.Check(t, (&key{"KEY1_TEST_KEY2", ""}).hasHints([]string{"key1", "test", "key2"}))
	assert.Check(t, !(&key{"key1_test_key2", ""}).hasHints([]string{""}))
	assert.Check(t, !(&key{"key1_test_key2", ""}).hasHints([]string{"notfound"}))
}

func TestKeyString(t *testing.T) {
	assert.Equal(t, (&key{"test", ""}).String(), "test")
	assert.Equal(t, (&key{"key1_test_key2", ""}).String(), "key1_test_key2")
	assert.Equal(t, (&key{"KEY1_TEST_KEY2", ""}).String(), "KEY1_TEST_KEY2")
	assert.Equal(t, (&key{"", ""}).String(), "")
}

func checkReplaceHint(t *testing.T, h *hintedKey, matched bool, result string, matches []string) {
	assert.Equal(t, h.matched, matched)
	assert.Equal(t, h.name, result)
	assert.DeepEqual(t, h.matches, matches)
}

func TestKeyReplaceHints(t *testing.T) {
	checkReplaceHint(t,
		(&key{"string", ""}).replaceHints([]string{"test"}),
		false, "string", []string{},
	)
	checkReplaceHint(t,
		(&key{"string", ""}).replaceHints([]string{"string"}),
		true, "XXX", []string{"string"},
	)
	checkReplaceHint(t,
		(&key{"string_key", ""}).replaceHints([]string{"key"}),
		true, "string_XXX", []string{"key"},
	)
	checkReplaceHint(t,
		(&key{"string_key_key", ""}).replaceHints([]string{"key"}),
		true, "string_XXX_XXX", []string{"key"},
	)
	checkReplaceHint(t,
		(&key{"KEY1_TEST_KEY2", ""}).replaceHints([]string{"key1", "key2"}),
		true, "XXX_test_XXX", []string{"key1", "key2"},
	)
	checkReplaceHint(t,
		(&key{"KEY_KEY1_KEY2", ""}).replaceHints([]string{"key"}),
		true, "XXX_key1_key2", []string{"key"},
	)
}

func checkFindSimilars(t *testing.T, k *key, c []string, hints []string, answer []string) {
	candidates := []*key{}
	for _, candidate := range c {
		candidates = append(candidates, &key{candidate, ""})
	}
	var result []string
	resultKeys := k.findSimilars(candidates, hints)
	for _, r := range resultKeys {
		result = append(result, r.name)
	}
	assert.DeepEqual(t, answer, result)
}
func TestKeyFindSimilars(t *testing.T) {
	checkFindSimilars(t, &key{"DATABASE_PORT", ""},
		[]string{"DATABASE_PORT", "DATABASE_HOST", "DATABASE_USER"}, []string{"user", "port", "host"},
		[]string{"DATABASE_PORT", "DATABASE_HOST", "DATABASE_USER"},
	)
	checkFindSimilars(t, &key{"DATABASE", ""},
		[]string{"DATABASE_PORT", "DATABASE_HOST", "DATABASE_USER"}, []string{"user", "port", "host"},
		nil,
	)
	checkFindSimilars(t, &key{"HOST_PORT", ""},
		[]string{"HOST_USER", "HOST_HOST", "HOST_PASS", "CLIENT_ID"}, []string{"host", "user", "port", "pass"},
		[]string{"HOST_USER", "HOST_HOST", "HOST_PASS"},
	)
	checkFindSimilars(t, &key{"API_ENDPOINT", ""},
		[]string{"API_USER", "API_PASS"}, []string{"endpoint", "user", "pass"},
		[]string{"API_USER", "API_PASS"},
	)

}
