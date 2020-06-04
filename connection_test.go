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

func TestSetInstance(t *testing.T) {
	assert := assert.New(t)
	var c0 Connection
	c0.SetInstance("host")
	assert.Equal("host", c0.Server, "host only")
	assert.Equal("", c0.Instance, "host only")

	var c1 Connection
	c1.SetInstance("host\\instance")
	assert.Equal("host", c1.Server, "host instance")
	assert.Equal("instance", c1.Instance, "host instance")
}
