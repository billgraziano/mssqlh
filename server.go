package mssqlh

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Server holds information about an instance of SQL Server
type Server struct {
	// Name holds the result of @@SERVERNAME
	Name      string
	Computer  string
	Instance  string
	Domain    string
	DNSSuffix string
	FQDN      string

	EngineEdition       int
	ProductVersion      string
	ProductMajorVersion int
	// Version, Instance, Azure, etc.
}

// GetServer gets details on the SQL Server
// TODO Create a queryer interface and accept that (suppport sqlx)
func GetServer(ctx context.Context, db *sql.DB) (Server, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	query := `
		SET XACT_ABORT ON;
		SET NOCOUNT ON;
		
		SELECT	COALESCE(@@SERVERNAME, '') AS atat_server_name
				,COALESCE(DEFAULT_DOMAIN(), '') as domain
				,COALESCE(CAST(SERVERPROPERTY('ServerName') AS NVARCHAR(128)), '') AS computer
				,COALESCE(CAST(SERVERPROPERTY('InstanceName') AS NVARCHAR(128)), '') AS instance
				,COALESCE(CAST(SERVERPROPERTY('EngineEdition') AS INT), 0) AS engine_edition
				,COALESCE(CAST(SERVERPROPERTY('ProductVersion') AS NVARCHAR(128)), '') AS product_version
		`
	row := db.QueryRowContext(ctx, query)
	var s Server
	err := row.Scan(&s.Name, &s.Domain, &s.Computer, &s.Instance, &s.EngineEdition, &s.ProductVersion)
	if err != nil {
		return s, errors.Wrap(err, "sql.scan")
	}
	s.FQDN = s.Computer

	if s.EngineEdition == 2 || s.EngineEdition == 3 || s.EngineEdition == 4 {
		fqdn_query := `
		
			DECLARE	@Suffix		VARCHAR(1024),
				@has_perms	INT

			SELECT @has_perms = COALESCE(HAS_PERMS_BY_NAME('xp_regread', 'OBJECT', 'EXECUTE'), 0)  
			IF @has_perms = 1 
			BEGIN
				EXEC master..xp_regread
					@rootkey = 'HKEY_LOCAL_MACHINE',
					@key = 'system\currentcontrolset\services\tcpip\parameters\',
					@value_name = 'Domain',
					@value = @Suffix OUTPUT
			END
			SELECT COALESCE(@Suffix, '') as dns_suffix 

		`
		row = db.QueryRowContext(ctx, fqdn_query)
		err = row.Scan(&s.DNSSuffix)
		if err != nil {
			return s, errors.Wrap(err, "fqdn.sql.scan")
		}
		if s.DNSSuffix != "" {
			s.FQDN = fmt.Sprintf("%s.%s", s.Computer, s.DNSSuffix)
		}
	}

	version := strings.Split(s.ProductVersion, ".")
	if len(version) != 4 {
		s.ProductMajorVersion = 0
	}
	s.ProductMajorVersion, err = strconv.Atoi(version[0])
	if err != nil {
		return s, fmt.Errorf("invalid product version: %s", s.ProductVersion)
	}

	return s, nil
}
