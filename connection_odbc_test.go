package mssqlh

import (
	"testing"

	"github.com/billgraziano/mssqlh/v2/odbch"
	"github.com/stretchr/testify/assert"
)

func TestODBCString(t *testing.T) {
	assert := assert.New(t)
	mock = true
	var tests = []struct {
		name     string
		in       Connection
		expected string
	}{
		{"empty", Connection{}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes;"},
		{"host", Connection{FQDN: "test"}, "Driver={SQL Server Native Client 11.0}; Server=test; Trusted_Connection=Yes;"},
		{"fqdn", Connection{FQDN: "test.example.com"}, "Driver={SQL Server Native Client 11.0}; Server=test.example.com; Trusted_Connection=Yes;"},
		{"host-instance", Connection{FQDN: "test\\junk"}, "Driver={SQL Server Native Client 11.0}; Server=test\\junk; Trusted_Connection=Yes;"},
		{"host-port", Connection{FQDN: "test,1433"}, "Driver={SQL Server Native Client 11.0}; Server=test,1433; Trusted_Connection=Yes;"},
		{"host-colon-port", Connection{FQDN: "test:1433"}, "Driver={SQL Server Native Client 11.0}; Server=test,1433; Trusted_Connection=Yes;"},
		{"user", Connection{User: "u1"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; UID=u1; PWD=;"},
		{"user-pass", Connection{User: "u1", Password: "pass"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; UID=u1; PWD=pass;"},
		{"appname", Connection{Application: "appy"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; App=appy;"},
		{"database", Connection{Database: "db"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Database=db;"},
		{"dial timeout", Connection{DialTimeout: 10}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Timeout=10;"},
		{"connect timeout", Connection{ConnectTimeout: 11}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Timeout=11;"},
		{"both timeout", Connection{ConnectTimeout: 13, DialTimeout: 10}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Timeout=23;"},
		{"encrypt-yes", Connection{Encrypt: EncryptYes}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Encrypt=Yes;"},
		{"encrypt-no", Connection{Encrypt: EncryptNo}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Encrypt=No;"},
		{"encrypt-strict", Connection{Encrypt: EncryptStrict}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Encrypt=Strict;"},
		{"encrypt-mandatory", Connection{Encrypt: EncryptMandatory}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Encrypt=Mandatory;"},
		{"encrypt-optional", Connection{Encrypt: EncryptOptional}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Encrypt=Optional;"},
		{"encrypt-other", Connection{Encrypt: "other"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Encrypt=other;"},
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
			"Driver={SQL Server Native Client 11.0}; Server=new\\in; UID=u2; PWD=nope; Database=txn; App=bonk; Timeout=5;",
		},
	}

	for _, v := range tests {
		v.in.Driver = DriverODBC
		v.in.ODBCDriver = odbch.NativeClient11

		actual := v.in.String()
		assert.Equal(v.expected, actual, v.name)
	}
}

func TestODBC18RequiresEncrypt(t *testing.T) {
	assert := assert.New(t)
	mock = true
	var tests = []struct {
		driver   string
		name     string
		in       Connection
		expected string
	}{
		{odbch.NativeClient11, "base", Connection{}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes;"},
		{odbch.ODBC17, "odbc17", Connection{}, "Driver={ODBC Driver 17 for SQL Server}; Server=localhost; Trusted_Connection=Yes;"},
		{odbch.ODBC18, "odbc18", Connection{}, "Driver={ODBC Driver 18 for SQL Server}; Server=localhost; Trusted_Connection=Yes; Encrypt=Optional;"},
		{odbch.ODBC18, "odbc18-EncryptNo", Connection{Encrypt: EncryptNo}, "Driver={ODBC Driver 18 for SQL Server}; Server=localhost; Trusted_Connection=Yes; Encrypt=No;"},
	}

	for _, v := range tests {
		v.in.Driver = DriverODBC
		v.in.ODBCDriver = v.driver

		actual := v.in.String()
		assert.Equal(v.expected, actual, v.name)
	}
}
