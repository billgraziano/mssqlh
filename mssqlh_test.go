//go:build windows

package mssqlh_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/billgraziano/mssqlh/odbch"

	"github.com/billgraziano/mssqlh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/microsoft/go-mssqldb"
)

func TestMSSQL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	env := os.Getenv("MSSQLH_SERVERS")
	if env == "" {
		return
	}
	r := csv.NewReader(strings.NewReader(env))
	servers, err := r.Read()
	require.NoError(err)
	for _, server := range servers {
		fmt.Printf("server: %s\r\n", server)
		server = strings.TrimSpace(server)

		// Test with "mssql" driver
		m, err := mssqlh.Open(server, "")
		require.NoError(err, "mssql open failed: %s", server)

		_, err = mssqlh.GetServer(context.TODO(), m)
		assert.NoError(err, "mssql getserver failed: %s", server)

		_, err = mssqlh.GetSession(context.TODO(), m)
		assert.NoError(err, "mssql getsession failed: %s", server)

		cxn := mssqlh.Connection{}
		cxn.FQDN = server
		cxn.Driver = mssqlh.DriverODBC

		// Test with each installed "odbc" driver
		drivers, err := odbch.InstalledDrivers()
		require.NoError(err)
		for _, d := range drivers {
			fmt.Printf("server: %s  driver: %s\r\n", server, d)
			cxn.ODBCDriver = d
			o, err := cxn.Open()
			require.NoError(err, "odbc getsession failed: %s", server)

			_, err = mssqlh.GetServer(context.TODO(), o)
			assert.NoError(err, "odbc getserver failed: %s", server)

			_, err = mssqlh.GetSession(context.TODO(), o)
			assert.NoError(err, "odbc getsession failed: %s", server)
		}
	}
}
