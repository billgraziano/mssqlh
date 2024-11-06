
* Happy path is the GO driver - Can override with SetDriver()
* Server object GetServer(sql.DB - but my interface) returns Server
* Session object GetSession(sql.DB - but my interface) returns Session
* Interface with just the stuff I need (like sqlx does)
  - Open, Query, etc.
* Connection object return a String() or a sql pool (or sqlx pool?)
Connection.Open() returns sql.DB
Connection.String() returns connection string - default to localhost if nothing given?
MustString, MustConnect

mssqlh.Open(fqdn) - defaults, exe name for app, returns sql.DB
mssqlh.ConnectionString(fqdn) - defaults but only a connection string
* read connection string attributes from KV map?  This helps with TOML and such.
* some private field to do sql or sqlx

QuoteName
QuoteString
FixSlashes
BatchSeparator

ParseConnectionString() - ?

Do we want to set the pool information here?
