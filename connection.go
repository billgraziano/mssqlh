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
	FQDN           string
	User           string
	Password       string
	Database       string
	Application    string
	DialTimeout    int
	ConnectTimeout int
	ODBCDriver     string
	//ExtraValues    map[string]string
	//Server         string
	//Instance       string
	//Port           int
}

// NewConnection returns a connection with sane defaults.
// You can specify the server "host[\instance]",
// "host:port", or "host,port" format.
func NewConnection(server, user, password, database, app string) Connection {
	conn := Connection{FQDN: server, User: user, Password: password, Database: database, Application: app}
	//conn.SetInstance(server)
	conn.setDefaults()
	//conn.ExtraValues = make(map[string]string)
	return conn
}

// Open connects to a SQL Server.  It accepts "host[\instance]",
// "host:port", or "host,port".
func Open(fqdn, database string) (*sql.DB, error) {
	conn := Connection{FQDN: fqdn, Database: database}
	return conn.Open()
}

// ServerName buids a string in the format server\instance or server:host.
// Most likely you won't have an instance and a port.  Plus I don't think
// that works.  This should be roughly what it tries to connect to.
func (c Connection) ServerName() string {
	c.setDefaults()
	return c.FQDN
	// s := c.Server
	// if c.Instance != "" {
	// 	s += "\\" + c.Instance
	// }
	// if c.Port != 0 {
	// 	s += fmt.Sprintf(":%d", c.Port)
	// }
	// return s
}

// String returns a connection string for the given connection.
// Setting Driver to an invalid type returns an unusable connection string
// but not an error.  It should be caught on Open
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

// SetInstance takes a server in the format "host[\instance]",
// "host:port", or "host,port" and
// assigns the Server, Instance, and Port
// func (c *Connection) SetInstance(s string) {
// 	if s == "" {
// 		return
// 	}
// 	c.Server, c.Instance, c.Port = parseFQDN(s)
// }

func (c Connection) Computer() string {
	c.setDefaults()
	computer, _, _ := parseFQDN(c.FQDN)
	return computer
}

func (c Connection) Instance() string {
	c.setDefaults()
	_, instance, _ := parseFQDN(c.FQDN)
	return instance
}

func (c Connection) Port() int {
	c.setDefaults()
	_, _, port := parseFQDN(c.FQDN)
	return port
}

// setDefaults sets defaults for the driver and host name.
// This allows and empty connection to be used
func (c *Connection) setDefaults() {
	if c.Driver == "" {
		c.Driver = DriverMSSQL
	}
	if c.FQDN == "" {
		c.FQDN = "localhost"
	}
}

// exeName gets the default app name which is the executable name
func exeName() (string, error) {
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
