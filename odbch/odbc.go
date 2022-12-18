package odbch

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ErrNoDrivers is returned if no valid ODBC SQL Server drivers are found
var ErrNoDrivers = errors.New("no drivers found")

// ErrInvalidDriver indiates that an ODBC SQL Server driver is invalid
var ErrInvalidDriver = errors.New("invalid driver")

// ODBCDriver is the name of an ODBC SQL Server Drive
//type ODBCDriver string

const (
	// NativeClient11 is an Native SQL Server Driver version 11
	NativeClient11 string = "SQL Server Native Client 11.0"

	// NativeClient10 is an Native SQL Server Driver version 10
	NativeClient10 string = "SQL Server Native Client 10.0"

	// ODBC18 is an ODBC SQL Server Driver version 18
	ODBC18 string = "ODBC Driver 18 for SQL Server"

	// ODBC17 is an ODBC SQL Server Driver version 17
	ODBC17 string = "ODBC Driver 17 for SQL Server"

	// ODBC13 is an ODBC SQL Server Driver version 13
	ODBC13 string = "ODBC Driver 13 for SQL Server"

	// ODBC11 is an ODBC SQL Server Driver version 11
	ODBC11 string = "ODBC Driver 11 for SQL Server"

	// GenericODBC is the Generic ODBC SQL Server driver
	GenericODBC string = "SQL Server"

	// NoDriver is an empty string. Usually used for error checking
	// NoDriver string = ""
)

var orderedDrivers = []string{
	ODBC18,
	ODBC17,
	ODBC13,
	NativeClient11,
	ODBC11,
	NativeClient10,
	GenericODBC,
}

var fs = afero.NewOsFs()
