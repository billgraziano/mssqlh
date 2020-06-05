package mssqlh

import "strings"

// QuoteName wraps a string in [brackets].
// It tries to match the functionality of SQL Server's QUOTENAME()
func QuoteName(s string) string {
	return "[" + strings.Replace(s, "]", "]]", -1) + "]"
}

// QuoteString wraps a string in single-quotes and tries to handle
// embedded quotes
func QuoteString(s string) string {
	return "'" + strings.Replace(s, "'", "''", -1) + "'"
}

// SplitBatch (TODO) splits a batch on GO
// Handles GO [count]
// Handle comments?
// Must be on a line byitself
func SplitBatch(batch string) []string {
	return []string{batch}
}
