package mssqlh

import (
	"database/sql"
	"errors"
)

// Connection is the basis for building connection strings
type Connection struct {
	// Driver sets the GO driver that will be used.
	// Leaving this blank defaults to DriverMSSQL.
	Driver         string
	Server         string
	Instance       string
	Port           int
	User           string
	Password       string
	Database       string
	AppName        string
	DialTimeout    int
	ConnectTimeout int
	ODBCDriver     string
	Parameters     map[string]string // Attributes? Properties? Values?
}

// String returns a connection string for the given connection.
// If no server is specified, localhost is used
// The executable name is used for app name if possible
func (c Connection) String() string {
	if c.Driver == DriverMSSQL || c.Driver == "" {
		return c.mssqlString()
	}
	return "tbd"
}

// Redacted returns a connection string with the password
// replaced with _redacted_
func (c Connection) Redacted() string {
	if c.Password != "" {
		c.Password = "redacted"
	}
	return c.String()
}

// Open returns a sql.DB pool
// Open("localhost"), Open("D10\Server"),
// Open("D10,1433"), Open("D10:1433")
func Open(string) (sql.DB, error) {
	return sql.DB{}, errors.New("not implemented")
}
