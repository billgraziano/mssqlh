package mssqlh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteName(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		in       string
		expected string
	}{
		{"test", "[test]"},
		{"te]st", "[te]]st]"},
		{"[test]", "[[test]]]"},
		{"te[]st", "[te[]]st]"},
		{"te[st", "[te[st]"},
	}

	for _, v := range tests {
		actual := QuoteName(v.in)
		assert.Equal(v.expected, actual, "bad quotename")
	}
}
