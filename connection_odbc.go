package mssqlh

import (
	"fmt"
	"strings"

	"github.com/billgraziano/mssqlh/odbc"
)

// mssqlString returns a connection string for a GO MSSQL connection
func (c Connection) odbcString() string {

	var str string

	// We copied c so we can make changes to it

	if c.ODBCDriver == "" {
		// We are swallowing this error.  Whatever uses this connection
		// string will fail and it can deal with the error then.  This
		// only happens if there are no drivers installed.
		c.ODBCDriver, _ = odbc.BestDriver()
	}
	str += fmt.Sprintf("Driver={%s}; ", c.ODBCDriver)

	if c.Server == "" {
		c.Server = "localhost"
	}

	str += fmt.Sprintf("Server=%s; ", c.getODBCServerName())

	if c.User == "" {
		str += fmt.Sprintf("Trusted_Connection=Yes; ")
	} else {
		str += fmt.Sprintf("UID=%s; PWD=%s; ", c.User, c.Password)
	}

	if c.Database != "" {
		str += fmt.Sprintf("Database=%s; ", c.Database)
	}

	if c.AppName != "" {
		str += fmt.Sprintf("App=%s; ", c.AppName)
	} else {
		if !mock {
			app, err := appName()
			if err == nil {
				fmt.Println(app)
				str += fmt.Sprintf("App=%s; ", app)
			}
		}
	}

	to := c.DialTimeout + c.ConnectTimeout
	if to > 0 {
		str += fmt.Sprintf("Timeout=%d; ", to)
	}

	return strings.TrimSpace(str)
}

// combineHostInstance returns the Server and Instance as Server\Instance
func (c *Connection) getODBCServerName() string {
	str := c.Server
	if c.Instance != "" {
		str += "\\" + c.Instance
	}
	if c.Port != 0 {
		str += fmt.Sprintf(",%d", c.Port)
	}
	return str
}
