package schemedetector

import (
	"testing"

	"gotest.tools/assert"
)

func TestKeyString(t *testing.T) {
	assert.Equal(t, (&key{"test", ""}).String(), "test")
	assert.Equal(t, (&key{"key1_test_key2", ""}).String(), "key1_test_key2")
	assert.Equal(t, (&key{"KEY1_TEST_KEY2", ""}).String(), "KEY1_TEST_KEY2")
	assert.Equal(t, (&key{"", ""}).String(), "")
}

func TestKeyFindSimilars(t *testing.T) {
	checkFindSimilars(t, &key{"A_B", ""},
		[]string{"A_C", "A_D", "A_E"},
		[]string{"A_C", "A_D", "A_E"},
	)
	checkFindSimilars(t, &key{"A_B", ""},
		[]string{"A_C", "A_D", "A_E", "C_D", "A", "B"},
		[]string{"A_C", "A_D", "A_E"},
	)
	checkFindSimilars(t, &key{"A", ""},
		[]string{"A_B", "A_C", "C_A_B", "C_B"},
		[]string{"A_B", "A_C"},
	)
	checkFindSimilars(t, &key{"A_B_C", ""},
		[]string{"A_B_D", "A_B_E", "A_D_E"},
		[]string{"A_B_D", "A_B_E"},
	)
	checkFindSimilars(t, &key{"A", ""},
		[]string{"B", "C"},
		nil,
	)
	checkFindSimilars(t, &key{"A_B_C", ""},
		[]string{"D_E_F", "G_H_I"},
		nil,
	)
}

func checkFindSimilars(t *testing.T, k *key, c []string, answer []string) {
	candidates := []*key{}
	for _, candidate := range c {
		candidates = append(candidates, &key{candidate, ""})
	}
	var result []string
	resultKeys := k.findSimilars(candidates)
	for _, r := range resultKeys {
		result = append(result, r.name)
	}
	assert.DeepEqual(t, answer, result)
}
