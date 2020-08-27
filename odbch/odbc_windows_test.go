// +build windows
package odbch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestDriver(t *testing.T) {
	s, err := BestDriver()
	if err != nil {
		t.Error("Best Driver error: ", err)
	}
	t.Log("Best Driver: ", s)
}

func TestODBCDriver(t *testing.T) {
	d := NativeClient11
	t.Log("Native 11: ", d)
}

func TestInstalledDrivers(t *testing.T) {
	t.Log("Available Drivers")
	t.Log("=====================================")
	d, err := InstalledDrivers()
	if err != nil {
		t.Error("available drivers: ", err)
	}
	for _, s := range d {
		t.Log(s)
	}
}

func TestValidDrivers(t *testing.T) {
	assert := assert.New(t)
	var err error
	err = ValidDriver("test")
	if err != ErrInvalidDriver {
		t.Error("'test' should be invalid driver")
	}

	err = ValidDriver("ODBC Driver 17 for SQL Server")
	assert.NoError(err)
}
