//go:build linux

package odbch

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDriver(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	memfs := afero.NewMemMapFs()
	fs = memfs
	err := afero.WriteFile(fs, "/etc/odbcinst.ini", []byte("[Test]\nv2=2]\n"), 0700)
	require.NoError(err)
	dd, err := getDrivers()
	assert.NoError(err, "getdrivers")
	assert.Equal(2, len(dd)) // The driver adds a DEFAULT section
}

func TestBestDriver(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	memfs := afero.NewMemMapFs()
	fs = memfs
	err := afero.WriteFile(fs, "/etc/odbcinst.ini", []byte(
		`[ODBC Driver 17 for SQL Server]
		v2=2
		
		[ODBC Driver 13 for SQL Server]
		another = 1
		`),
		0700)
	require.NoError(err)
	d, err := BestDriver()
	if err != nil {
		t.Error("Best Driver error: ", err)
	}
	assert.Equal(ODBC17, d)
}

func TestInstalledDrivers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	memfs := afero.NewMemMapFs()
	fs = memfs
	err := afero.WriteFile(fs, "/etc/odbcinst.ini", []byte(
		`[ODBC Driver 17 for SQL Server]
		v2=2
		
		[ODBC Driver 13 for SQL Server]
		another = 1
		`),
		0700)
	require.NoError(err)
	t.Log("Available Drivers")
	t.Log("=====================================")
	d, err := InstalledDrivers()
	assert.Equal(2, len(d))
	if err != nil {
		t.Error("available drivers: ", err)
	}
	for _, s := range d {
		t.Log(s)
	}
}

func TestValidDrivers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	memfs := afero.NewMemMapFs()
	fs = memfs
	err := afero.WriteFile(fs, "/etc/odbcinst.ini", []byte(
		`[ODBC Driver 17 for SQL Server]
		v2=2
		
		[ODBC Driver 13 for SQL Server]
		another = 1
		`),
		0700)
	require.NoError(err)

	err = ValidDriver("test")
	assert.EqualError(err, ErrInvalidDriver.Error())

	err = ValidDriver("ODBC Driver 17 for SQL Server")
	assert.NoError(err)
}
