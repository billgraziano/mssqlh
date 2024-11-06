package mssqlh

// DriverMSSQL uses the https://github.com/microsoft/go-mssqldb library.
// This is the driver for this package.
var DriverMSSQL = "sqlserver"

// DriverODBC uses the https://github.com/alexbrainman/odbc library
var DriverODBC = "odbc"

// used to enable mocks for testing
var mock bool
