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
	// Leaving this blank defaults to DriverMSSQL (the native GO driver).
	Driver         string
	Protocol       string
	FQDN           string
	User           string
	Password       string
	Database       string
	Application    string
	DialTimeout    int
	ConnectTimeout int
	ODBCDriver     string
	Encrypt        string
}

// NewConnection returns a connection with sane defaults.
// You can specify the server "host[\instance]",
// "host:port", or "host,port" format.
// You can prefix the host with the protocol.  For example, "tcp:host,1433".
// Supported protocols are "tcp", "np" (named pipes), "lpc" (shared memory).
// Optionally you can set the Protocol directly in the Connection object
func NewConnection(server, user, password, database, app string) Connection {
	conn := Connection{FQDN: server, User: user, Password: password, Database: database, Application: app}
	//conn.SetInstance(server)
	conn.setDefaults()
	//conn.ExtraValues = make(map[string]string)
	return conn
}

// Open connects to a SQL Server.  It accepts "host[\instance]",
// "host:port", or "host,port".
// You can prefix the host with the protocol.  For example, "tcp:host,1433".
// Supported protocols are "tcp", "np" (named pipes), "lpc" (shared memory).
// Optionally you can set the Protocol directly in the Connection object
func Open(fqdn, database string) (*sql.DB, error) {
	conn := Connection{FQDN: fqdn, Database: database}
	conn.setDefaults()
	return conn.Open()
}

// ServerName buids a string in the format server\instance or server:host.
// Most likely you won't have an instance and a port.  Plus I don't think
// that works.  This should be roughly what it tries to connect to.
func (c Connection) ServerName() string {
	c.setDefaults()
	return c.FQDN
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
// replaced with "redacted".  Optionally, you can specify
// how many characters of the password to include
func (c Connection) Redacted(n int) string {
	c.setDefaults()
	if c.Password == "" {
		return c.String()
	}
	if n >= len(c.Password) {
		return c.String()
	}
	if n == 0 {
		c.Password = "redacted"
		return c.String()
	}
	pwd := c.Password[:n]
	c.Password = fmt.Sprintf("%s_redacted", pwd)
	return c.String()
}

// Open connects to the SQL Server
func (c Connection) Open() (*sql.DB, error) {
	c.setDefaults()
	return sql.Open(c.Driver, c.String())
}

// Computer returns the computer (or host) name from FQDN
func (c Connection) Computer() string {
	c.setDefaults()
	_, computer, _, _ := parseFQDN(c.FQDN)
	return computer
}

// Instance returns the instance from FQDN
func (c Connection) Instance() string {
	c.setDefaults()
	_, _, instance, _ := parseFQDN(c.FQDN)
	return instance
}

// Port returns the port from FQDN.  It returns 0 if no port.
func (c Connection) Port() int {
	c.setDefaults()
	_, _, _, port := parseFQDN(c.FQDN)
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
	protocol, server := stripProtocol(c.FQDN)
	if protocol != "" {
		c.Protocol = protocol
		c.FQDN = server
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

// parse FQDN splits a host\instance with optional port and protocol
func parseFQDN(s string) (protocol, host, instance string, port int) {
	var err error
	protocol, s = stripProtocol(s)
	parts := strings.FieldsFunc(s, hostSplitter)
	host = parts[0]
	if len(parts) == 1 {
		return protocol, host, "", 0
	}
	if len(parts) == 2 {
		port, err = strconv.Atoi(parts[1])
		if err == nil {
			return protocol, host, "", port
		}
		instance = parts[1]
		return protocol, host, instance, 0
	}
	if len(parts) == 3 {
		instance = parts[1]
		port, _ = strconv.Atoi(parts[2])
		return protocol, host, instance, port
	}

	return protocol, host, instance, port
}

// hostSplitter splits a string on :,\ and is used to split FQDN names
func hostSplitter(r rune) bool {
	return r == ':' || r == ',' || r == '\\'
}

// stripProtocol accepts a server name or FQDN and the format [protocol:]FQDN and returns protocol and FQDN.
// It supports "np", "lpc" (shared memory), and "tcp".
// This function is case-sensitive.
func stripProtocol(server string) (string, string) {
	prefixes := []string{"tcp:", "np:", "lpc:"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(server, prefix) {
			return strings.TrimSuffix(prefix, ":"), strings.TrimPrefix(server, prefix)
		}
	}
	return "", server
}
