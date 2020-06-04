package mssqlh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedacted(t *testing.T) {
	assert := assert.New(t)
	conn := Connection{
		User:     "u1",
		Password: "pass",
	}
	assert.Equal("sqlserver://u1:redacted@localhost", conn.Redacted())
	assert.Equal("pass", conn.Password)
}
