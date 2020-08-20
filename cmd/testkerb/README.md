TESTKERB
========

This is a simple utility to see if KERBEROS connectivity is working.
It reads from `servers.txt` or the command-line for servers to test.  It also works well as a simple connectivity test tool.

It uses the driver at [https://github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb).

It expects the `servers.txt` file to include raw connection strings in that format.  Samples are provided.

Alternatively, you can provide connection strings on the command-line.  Provide them space separated:

`testkerb.exe sqlserver://d40/sql2017 sqlserver://d40/sql2019`

Usage
-----

```
Usage of testkerb:
  -app string
        sets the application name (default "testkerb")
  -debug
        enable debug messages
  -file string
        list of servers (default ".\\servers.txt")
  -log string
        log=3 for driver logging (can be messy)
```