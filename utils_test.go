package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestUtilsStringInArray(t *testing.T) {
	assert.Check(t, stringInArray("test", []string{"test"}))
	assert.Check(t, stringInArray("test", []string{"test", "test"}))
	assert.Check(t, stringInArray("test", []string{"notfound", "test"}))
	assert.Check(t, stringInArray("test", []string{"test", "notfound"}))
	assert.Check(t, stringInArray("", []string{""}))
	assert.Check(t, !stringInArray("", []string{"test", "notfound"}))
}
