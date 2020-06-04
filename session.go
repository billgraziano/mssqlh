package mssqlh

import (
	"database/sql"
	"errors"
)

// Session stores information about the connection to SQL Server
type Session struct {
	SessionID int
}

func GetSession(cxn *sql.DB) (Session, error) {
	return Session{}, errors.New("not implemented")
}
