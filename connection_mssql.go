package mssqlh

import (
	"fmt"
	"net/url"
	"strconv"
)

// mssqlString returns a connection string for a GO MSSQL connection
func (c Connection) mssqlString() string {
	// We copied c so we can make changes to it
	if c.Server == "" {
		c.Server = "localhost"
	}

	query := url.Values{}
	// TODO get the app name from the executable if blank (testing?)
	if c.AppName != "" {
		query.Add("app name", c.AppName)
	}
	if c.Database != "" {
		query.Add("database", c.Database)
	}
	if c.DialTimeout > 0 {
		query.Add("dial timeout", strconv.Itoa(c.DialTimeout))
	}
	if c.ConnectTimeout > 0 {
		query.Add("connect timeout", strconv.Itoa(c.ConnectTimeout))
	}

	// TODO apply any parameters
	u := &url.URL{
		Scheme:   "sqlserver",
		Host:     hostPortString(c.Server, c.Port),
		Path:     c.Instance,
		RawQuery: query.Encode(),
	}
	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	return u.String()
}

// hostPortString converts a host and port into host:port format
// If port is zero, it is ommitted
func hostPortString(host string, port int) string {
	if port == 0 {
		return host
	}
	return fmt.Sprintf("%s:%d", host, port)
}
