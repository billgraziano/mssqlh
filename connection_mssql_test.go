package mssqlh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMSSQLString(t *testing.T) {
	assert := assert.New(t)
	mock = true
	var tests = []struct {
		name     string
		in       Connection
		expected string
	}{
		{"empty", Connection{}, "sqlserver://localhost"},
		{"host", Connection{Server: "test"}, "sqlserver://test"},
		{"fqdn", Connection{Server: "test.example.com"}, "sqlserver://test.example.com"},
		{"host-instance", Connection{Server: "test", Instance: "junk"}, "sqlserver://test/junk"},
		{"host-port", Connection{Server: "test", Port: 1433}, "sqlserver://test:1433"},
		{"user", Connection{User: "u1"}, "sqlserver://u1:@localhost"},
		{"user-pass", Connection{User: "u1", Password: "pass"}, "sqlserver://u1:pass@localhost"},
		{"appname", Connection{AppName: "appy"}, "sqlserver://localhost?app+name=appy"},
		{"database", Connection{Database: "db"}, "sqlserver://localhost?database=db"},
		{"dial timeout", Connection{DialTimeout: 10}, "sqlserver://localhost?dial+timeout=10"},
		{"connect timeout", Connection{ConnectTimeout: 11}, "sqlserver://localhost?connect+timeout=11"},
		{
			"big",
			Connection{
				Server:      "new",
				Instance:    "in",
				User:        "u2",
				Password:    "nope",
				AppName:     "bonk",
				Database:    "txn",
				DialTimeout: 5,
			},
			"sqlserver://u2:nope@new/in?app+name=bonk&database=txn&dial+timeout=5",
		},
	}

	for _, v := range tests {
		actual := v.in.String()
		assert.Equal(v.expected, actual, v.name)
	}
}
