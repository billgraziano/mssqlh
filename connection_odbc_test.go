package mssqlh

import (
	"testing"

	"github.com/billgraziano/mssqlh/odbch"
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
		{"host", Connection{Server: "test"}, "Driver={SQL Server Native Client 11.0}; Server=test; Trusted_Connection=Yes;"},
		{"fqdn", Connection{Server: "test.example.com"}, "Driver={SQL Server Native Client 11.0}; Server=test.example.com; Trusted_Connection=Yes;"},
		{"host-instance", Connection{Server: "test", Instance: "junk"}, "Driver={SQL Server Native Client 11.0}; Server=test\\junk; Trusted_Connection=Yes;"},
		{"host-port", Connection{Server: "test", Port: 1433}, "Driver={SQL Server Native Client 11.0}; Server=test,1433; Trusted_Connection=Yes;"},
		{"user", Connection{User: "u1"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; UID=u1; PWD=;"},
		{"user-pass", Connection{User: "u1", Password: "pass"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; UID=u1; PWD=pass;"},
		{"appname", Connection{AppName: "appy"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; App=appy;"},
		{"database", Connection{Database: "db"}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Database=db;"},
		{"dial timeout", Connection{DialTimeout: 10}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Timeout=10;"},
		{"connect timeout", Connection{ConnectTimeout: 11}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Timeout=11;"},
		{"both timeout", Connection{ConnectTimeout: 13, DialTimeout: 10}, "Driver={SQL Server Native Client 11.0}; Server=localhost; Trusted_Connection=Yes; Timeout=23;"},
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
