package mssqlh

import "strings"

// QuoteName wraps a string in brackets ([]).
// It tries to match the functionality of SQL Server's QUOTENAME()
func QuoteName(s string) string {
	return "[" + strings.Replace(s, "]", "]]", -1) + "]"
}
