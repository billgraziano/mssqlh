TestKerb
========

This is a simple utility to see if KERBEROS connectivity is working.  It's also pretty good at general connectivity testing.

* It reads from `servers.txt` or the command-line for servers to test.

* It uses the driver at [https://github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb).

* It expects `servers.txt` file to include raw connection strings in that format.  Samples are provided.

* Alternatively, you can provide connection strings on the command-line.  Provide them space separated:

      testkerb.exe  sqlserver://host/instance1  sqlserver://host/instance2

Parameters
-----

```
  -app string
        sets the application name (default "testkerb")
  -debug
        enables debug messages
  -file string
        file with list of servers (default "servers.txt")
  -log string
        see github.com/denisenkom/go-mssqldb
```

Setting `-app ""` will omit the application name.

