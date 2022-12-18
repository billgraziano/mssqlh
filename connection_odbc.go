package mssqlh

import (
	"fmt"
	"strings"

	"github.com/billgraziano/mssqlh/odbch"
)

// odbcString returns a connection string for a GO ODBC connection
func (c Connection) odbcString() string {

	var str string

	// We copied c so we can make changes to it
	c.setDefaults()

	if c.ODBCDriver == "" {
		// We are swallowing this error.  Whatever uses this connection
		// string will fail and it can deal with the error then.  This
		// only happens if there are no drivers installed.
		c.ODBCDriver, _ = odbch.BestDriver()
	}
	str += fmt.Sprintf("Driver={%s}; ", c.ODBCDriver)

	str += fmt.Sprintf("Server=%s; ", c.getODBCServerName())

	if c.User == "" {
		str += "Trusted_Connection=Yes; "
	} else {
		str += fmt.Sprintf("UID=%s; PWD=%s; ", c.User, c.Password)
	}

	if c.Database != "" {
		str += fmt.Sprintf("Database=%s; ", c.Database)
	}

	if c.Application != "" {
		str += fmt.Sprintf("App=%s; ", c.Application)
	} else {
		if !mock {
			app, err := exeName()
			if err == nil { // swallow any error
				str += fmt.Sprintf("App=%s; ", app)
			}
		}
	}

	to := c.DialTimeout + c.ConnectTimeout
	if to > 0 {
		str += fmt.Sprintf("Timeout=%d; ", to)
	}

	c.Encrypt = strings.TrimSpace(c.Encrypt)
	// ODBC18 requires a value for Encrypt.
	// Use Optional if no value specified.
	if c.Encrypt == "" && c.ODBCDriver == odbch.ODBC18 {
		c.Encrypt = EncryptOptional
	}
	if c.Encrypt != "" {
		str += fmt.Sprintf("Encrypt=%s; ", c.Encrypt)
	}

	return strings.TrimSpace(str)
}

// combineHostInstance returns the Server and Instance as Server\Instance
func (c *Connection) getODBCServerName() string {
	str := c.Computer()
	if c.Instance() != "" {
		str += "\\" + c.Instance()
	}
	if c.Port() != 0 {
		str += fmt.Sprintf(",%d", c.Port())
	}
	return str
}
