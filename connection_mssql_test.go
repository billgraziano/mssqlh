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
		{"host", Connection{FQDN: "test"}, "sqlserver://test"},
		{"fqdn", Connection{FQDN: "test.example.com"}, "sqlserver://test.example.com"},
		{"host-instance", Connection{FQDN: "test\\junk"}, "sqlserver://test/junk"},
		{"host-comma-port", Connection{FQDN: "test,1433"}, "sqlserver://test:1433"},
		{"host-colon-port", Connection{FQDN: "test:1433"}, "sqlserver://test:1433"},
		{"user", Connection{User: "u1"}, "sqlserver://u1:@localhost"},
		{"user-pass", Connection{User: "u1", Password: "pass"}, "sqlserver://u1:pass@localhost"},
		{"appname", Connection{Application: "appy"}, "sqlserver://localhost?app+name=appy"},
		{"database", Connection{Database: "db"}, "sqlserver://localhost?database=db"},
		{"dial timeout", Connection{DialTimeout: 10}, "sqlserver://localhost?dial+timeout=10"},
		{"connect timeout", Connection{ConnectTimeout: 11}, "sqlserver://localhost?connect+timeout=11"},
		{"encrypt-yes", Connection{Encrypt: EncryptYes}, "sqlserver://localhost?encrypt=true"},
		{"encrypt-no", Connection{Encrypt: EncryptNo}, "sqlserver://localhost?encrypt=false"},
		{"encrypt-optional", Connection{Encrypt: EncryptOptional}, "sqlserver://localhost"},
		{"encrypt-strict", Connection{Encrypt: EncryptStrict}, "sqlserver://localhost?encrypt=true"},
		{"encrypt-mandatory", Connection{Encrypt: EncryptMandatory}, "sqlserver://localhost?encrypt=true"},
		{"encrypt-other", Connection{Encrypt: "other"}, "sqlserver://localhost?encrypt=other"},
		{
			"big",
			Connection{
				FQDN:        "new\\in",
				User:        "u2",
				Password:    "nope",
				Application: "bonk",
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
