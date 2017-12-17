package snapshot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateNameTag(t *testing.T) {
	expected := testNameTag()
	actual := GenerateNameTag(testInstanceName(), testDeviceName())

	assert.Equal(t, expected, actual, "Generated nametag meets expectations.")
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice(testStringTrue(), exampleSlice()), "Expected string found in slice")
	assert.False(t, StringInSlice(testStringFalse(), exampleSlice()), "Unexpected string not found in slice")
}
