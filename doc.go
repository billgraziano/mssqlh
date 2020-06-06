/*
Package mssqlh provides connection string building and helper routines for working with Microsoft SQL Server.

This package provides support for connecting to SQL Server using either
https://github.com/denisenkom/go-mssqldb (mssql driver) or https://github.com/alexbrainman/odbc (odbc driver).

Using the Connection type, you should be able to switch seamlessly between the two.
The package defaults to the "mssql" driver usless you specify the "odbc" driver.

Example using Open
	db, err := mssqlh.Open(fqdn)

This uses a trusted connection to the designated server using the "mssql" driver.  It accepts
server.domain.com, server\instance, server,port, or server:port.

Example code using NewConnection:
	cxn := mssqlh.NewConnection("localhost", "", "", "myapp")
	db, err := sql.Open("mssql", cxn.String())

If you don't pass a user and password it defaults to a trusted connection.

Example using the Connection type:
	cxn := mssqlh.Connect{
		Server:      "db-txn.corp.loc",
		Application: "myapp",
		DialTimeout: 15,
	}
	cxn.Database = "TXNDB"
	db, err := cxn.Open()

Defaults

The package provides the following defaults
(1) if no server is specified, it will use localhost,
(2) if no user is specified, it will default to a trusted connection
(3) if no application name is specied, it will default to the name of the executable

Using the ODBC driver

The subpackage odbch provides additional support for
using ODBC driver (https://github.com/alexbrainman/odbc)

Example code using the Connection object:
	cxn := mssqlh.Connect{
		Driver: mssqlh.DriverODBC,
		ODBCDriver: odbch.NativeClient11,
		Server: "localhost",
	}
	db, err := cxn.Open()

This connects using the ODBC driver.

Version Support

GetServer and GetSession should support SQL Server 2005 and beyond.  They
have been tested on SQL Server 2014 through SQL Server 2019.

There is limited testing with Azure SQL Databases.  The GetSession method
requires VIEW DATABASE STATE permission.

*/
package mssqlh
