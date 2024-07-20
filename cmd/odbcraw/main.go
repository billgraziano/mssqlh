package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/alexbrainman/odbc"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
)

func main() {
	var err error
	var file = flag.String("file", "servers.txt", "file with list of servers")
	flag.Parse()

	var servers []string

	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	kerberos := color.New(color.FgWhite).SprintFunc()
	ntlm := color.New(color.FgYellow).SprintFunc()

	if len(flag.Args()) == 0 {
		servers, err = readfile(*file)
		if err != nil {
			fmt.Fprintf(color.Output, "%s\n", red(err))
			os.Exit(1)
		}
	} else {
		servers = flag.Args()
	}

	if len(servers) == 0 {
		fmt.Fprintf(color.Output, "%s\n", red("No servers. Edit servers.txt or provide on the command line"))
		os.Exit(1)
	}

	for _, v := range servers {
		s := strings.TrimSpace(v)
		fmt.Fprintf(color.Output, "%s", cyan(s))

		db, err := sqlx.Open("odbc", s)
		if err != nil {
			msg := sqlerror(err)
			fmt.Fprintf(color.Output, " ==> %s\n", red(msg))
			continue
		}

		stmt := `
			SELECT	@@SERVERNAME AS ServerName,
				auth_scheme, 
				program_name,
				client_version,
				client_interface_name,
				login_name,
				net_transport
			--,* 
			FROM sys.dm_exec_connections c
			JOIN sys.dm_exec_sessions s ON s.session_id = c.session_id
			WHERE c.session_id = @@SPID
		`
		spid := SPID{}
		err = db.Get(&spid, stmt)
		if err != nil {
			msg := sqlerror(err)
			fmt.Fprintf(color.Output, " ==> %s\n", red(msg))
			continue
		}

		msg := fmt.Sprintf("%s (%s - %s)", spid.ServerName, spid.AuthScheme, spid.NetTransport)
		if spid.AuthScheme == "KERBEROS" {
			fmt.Fprintf(color.Output, " ==> %s\n", kerberos(msg))
		} else {
			fmt.Fprintf(color.Output, " ==> %s\n", ntlm(msg))
		}
	}
}

// sqlerror formats a SQL Server error
func sqlerror(err error) string {
	e, ok := err.(mssql.Error)
	if !ok {
		return err.Error()
	}
	str := fmt.Sprintf("Error: %s (Msg %d, Level %d, State %d, Line %d)", e.Message, e.Number, e.Class, e.State, e.LineNo)
	return str
}

// readfile returns a list of servers to test
func readfile(file string) ([]string, error) {
	var servers []string
	f, err := os.Open(file)
	if err != nil {
		return servers, errors.Wrap(err, "os.open")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		// no prefix and not empty
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			servers = append(servers, line)
		}
	}
	return servers, nil
}

// SPID holds the details for a connection
type SPID struct {
	ServerName          string `db:"ServerName"`
	AuthScheme          string `db:"auth_scheme"`
	ProgramName         string `db:"program_name"`
	ClientVersion       int32  `db:"client_version"`
	ClientInterfaceName string `db:"client_interface_name"`
	LoginName           string `db:"login_name"`
	NetTransport        string `db:"net_transport"`
}
