# README for SQL Server and Linux (ODBC)

The best solution at this point is to specify the ODBC driver name when using Linux.

## Notes

* The `odbcinst.ini` file holds the installed drivers.  The possible locations for the file include:
    * `/etc/odbcinst.ini`
    * `/usr/local/etc/odbcinst.ini`
* Darwin (macOS) support is very limited

## Roadmap

* Maybe just start with a default of ODBC 17, try to find a better one, and ignore any errors along the way
* Accept array of paths to search for odbcinst.ini 
* Have default hard coded
* Search the file for a suitable driver
* Consider using `odbcinst -j` to find the file locatiion

## Utilities

* `odbcinst -j` lists the installed drivers
* `odbc_config --odbcinstini` lists the location of the `odbcinst.ini` file

## Resources

* https://docs.microsoft.com/en-us/sql/connect/odbc/linux-mac/installing-the-microsoft-odbc-driver-for-sql-server?view=sql-server-ver15#driver-files

* https://docs.microsoft.com/en-us/sql/connect/odbc/download-odbc-driver-for-sql-server?view=sql-server-ver15


