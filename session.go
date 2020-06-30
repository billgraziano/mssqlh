package mssqlh

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Session stores information about the connection to SQL Server
type Session struct {
	Server          string    `db:"atat_server_name"`
	ID              int       `db:"session_id"`
	ConnectTime     time.Time `db:"connect_time"`
	LoginTime       time.Time `db:"login_time"`
	ClientInterface string    `db:"client_interface_name"`
	ClientVersion   int       `db:"client_version"`
	AuthScheme      string    `db:"auth_scheme"`
	Application     string
	Login           string
	Database        string
}

// GetSession gets details on the current connection to SQL Server
// TODO Create a queryer interface and accept that
// This requires VIEW DATABASE STATE in Azure
func GetSession(ctx context.Context, db *sql.DB) (Session, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	query := `
		SELECT	@@SERVERNAME AS atat_server_name
				,c.session_id
				,c.connect_time
				,COALESCE(c.auth_scheme, '') AS auth_scheme
				,s.login_time
				,COALESCE(s.client_interface_name, '') AS client_interface_name
				,COALESCE(s.client_version, 0) AS client_version
				,COALESCE(s.program_name, '') AS program_name 
				,COALESCE(s.login_name, '') AS login_name
				,COALESCE(DB_NAME(), '') AS database_name
		FROM	sys.dm_exec_connections c
		JOIN	sys.dm_exec_sessions s ON s.session_id = c.session_id
		WHERE	c.session_id = @@SPID;
		`
	row := db.QueryRowContext(ctx, query)
	var s Session
	err := row.Scan(&s.Server, &s.ID, &s.ConnectTime, &s.AuthScheme, &s.LoginTime, &s.ClientInterface, &s.ClientVersion, &s.Application, &s.Login, &s.Database)
	if err != nil {
		return s, errors.Wrap(err, "sql.scan")
	}

	return s, nil
}
