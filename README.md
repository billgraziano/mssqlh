# SQL Server Helper Library
Package `mssqlh` provides connection string building and helper routines for working with Microsoft SQL Server.

This package provides support for connecting to SQL Server using either:
* https://github.com/denisenkom/go-mssqldb (native GO driver) 
* https://github.com/alexbrainman/odbc (ODBC driver)

Using the Connection type, you should be able to switch seamlessly between the two.
The package defaults to the "mssql" driver (`mssqlh.DriverMSSQL`) usless you specify the "odbc" driver (`mssqlh.DriverODBC`).

Example using Open:
```go
db, err := mssqlh.Open(fqdn)
```


This uses a trusted connection to the designated server using the `mssql` driver.  It accepts
server.domain.com, server\instance, server,port, or server:port.

Example code using NewConnection:
```go
cxn := mssqlh.NewConnection("localhost", "", "", "myapp")
db, err := sql.Open("mssql", cxn.String())
```

If you don't pass user and password, it defaults to a trusted connection.

Example using the Connection type:
```go
cxn := mssqlh.Connect{
	FQDN:        "db-txn.corp.loc",
	Application: "myapp",
	DialTimeout: 15,
}
cxn.Database = "TXNDB"
db, err := cxn.Open()
```

## Defaults

The package provides the following defaults
1. If no server is specified, use localhost
2. If no user is specified, default to a trusted connection
3. If no application name is specified, default to the name of the executable

## Using the ODBC driver

The subpackage `odbch` provides additional support for
using ODBC driver (https://github.com/alexbrainman/odbc)

Example code using the Connection object:

	cxn := mssqlh.Connect{
		Driver:     mssqlh.DriverODBC,
		ODBCDriver: odbch.ODBC18,
		FQDN:       "localhost",
	}
	db, err := cxn.Open()

This connects using the specified ODBC driver.

## SQL Server Version Support

`GetServer` and `GetSession` should support SQL Server 2005 and beyond.  They
have been tested on SQL Server 2014 through SQL Server 2019.

There is limited testing with Azure SQL Databases.  The `GetSession` method
requires VIEW DATABASE STATE permission.

## Linux
It should support Linux but this has recived very little testing.

Linux looks for the following files to locate installed ODBC drivers:
* /usr/local/etc/odbcinst.ini
* /etc/odbcinst.ini

## Applications
The system comes with three sample applications 
* `mssqlh.exe` is a sample application
* `odbcraw.exe` can test ODBC connections from a `settings.txt` file
* `testkerbexe` can test connections for Kerberos from a `settings.txt` file
