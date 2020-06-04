package mssqlh

// Server holds information about an instance of SQL Server
type Server struct {
	// Server holds the result of @@SERVERNAME
	Server   string
	Computer string
	Instance string
	Domain   string
	// Version, Instance, Azure, etc.
}
