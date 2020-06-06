package mssqlh

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
	ExtraValues    map[string]string
}

/*
NewConnection returns a connection with sane defaults
*/
func NewConnection(server, user, password, app string) Connection {
	conn := Connection{User: user, Password: password, AppName: app}
	conn.SetInstance(server)
	conn.setDefaults()
	conn.ExtraValues = make(map[string]string)
	return conn
}

// String returns a connection string for the given connection.
// Setting Driver to an invalid type returns an unusable connection string
func (c Connection) String() string {
	c.setDefaults()
	switch c.Driver {
	case DriverMSSQL:
		return c.mssqlString()
	case DriverODBC:
		return c.odbcString()
	}
	return fmt.Sprintf("invalid driver: %s", c.Driver)
}

// Redacted returns a connection string with the password
// replaced with "redacted"
func (c Connection) Redacted() string {
	if c.Password != "" {
		c.Password = "redacted"
	}
	return c.String()
}

// Open connects to the SQL Server
func (c Connection) Open() (*sql.DB, error) {
	c.setDefaults()
	return sql.Open(c.Driver, c.String())
}

// setDefaults sets defaults for the driver and host name.
// This allows and empty connection to be used
func (c *Connection) setDefaults() {
	if c.Driver == "" {
		c.Driver = DriverMSSQL
	}
	if c.Server == "" {
		c.Server = "localhost"
	}
}

// Open connects to a SQL Server.  It accepts "host[\instance]",
// "host:port", or "host,port".
func Open(fqdn string) (*sql.DB, error) {
	var conn Connection
	conn.SetInstance(fqdn)
	return conn.Open()
}

// SetInstance takes a server in the format "host[\instance]",
// "host:port", or "host,port" and
// assigns the Server, Instance, and Port
func (c *Connection) SetInstance(s string) {
	if s == "" {
		return
	}
	c.Server, c.Instance, c.Port = parseFQDN(s)
}

// appName gets the default app name which is the fully
// qualified executable name
func appName() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err, "os.executable")
	}
	return filepath.Base(exe), nil
}

// parse FQDN splits a host\instance with an optional port
func parseFQDN(s string) (host, instance string, port int) {
	var err error
	parts := strings.FieldsFunc(s, hostSplitter)
	host = parts[0]
	if len(parts) == 1 {
		return host, "", 0
	}
	if len(parts) == 2 {
		port, err = strconv.Atoi(parts[1])
		if err == nil {
			return host, "", port
		}
		instance = parts[1]
		return host, instance, 0
	}
	if len(parts) == 3 {
		instance = parts[1]
		port, _ = strconv.Atoi(parts[2])
		return host, instance, port
	}

	return host, instance, port
}

// hostSplitter splits a string on :,\ and is used to split FQDN names
func hostSplitter(r rune) bool {
	return r == ':' || r == ',' || r == '\\'
}
