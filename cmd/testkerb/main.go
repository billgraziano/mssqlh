package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func main() {
	var err error
	var file = flag.String("file", ".\\servers.txt", "list of servers")
	var driverLog = flag.String("log", "", "log=3 for driver logging (messy)")
	var appName = flag.String("app", "testkerb", "sets the application name")
	var debug = flag.Bool("debug", false, "enable debug messages")
	flag.Parse()

	var servers []string

	if len(flag.Args()) == 0 {
		servers, err = readfile(*file)
		if err != nil {
			log.Fatal(errors.Wrap(err, "readfile"))
		}
	} else {
		servers = flag.Args()
	}

	if len(servers) == 0 {
		log.Fatal("no servers. edit servers.txt or provide on command line")
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	kerberos := color.New(color.FgWhite).SprintFunc()
	ntlm := color.New(color.FgYellow).SprintFunc()

	for _, v := range servers {
		dirty := false
		s := strings.TrimSpace(v)
		fmt.Fprintf(color.Output, "%s", cyan(s))

		// Parse and adjust the connections
		sqlurl, err := url.Parse(s)
		if err != nil {
			fmt.Fprintf(color.Output, " ==> %s\n", red(err.Error()))
			continue
		}

		query := sqlurl.Query()

		if *driverLog != "" {
			query.Set("log", *driverLog)
			dirty = true
		}

		if *appName != "" {
			query.Set("app name", *appName)
			dirty = true
		}

		if dirty || *debug {
			sqlurl.RawQuery = query.Encode()
			s = sqlurl.String()
			if *debug {
				fmt.Printf(" ==> %s", s)
			}
		}

		db, err := sqlx.Open("mssql", s)
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
		if *driverLog != "" {
			fmt.Println("\n--- Log ------------------------------")
		}
		err = db.Get(&spid, stmt)
		if err != nil {
			msg := sqlerror(err)
			fmt.Fprintf(color.Output, " ==> %s\n", red(msg))
			continue
		}
		if *driverLog != "" {
			fmt.Println("--------------------------------------")
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
	csvFile, err := os.Open(file)
	if err != nil {
		return servers, errors.Wrap(err, "os.open")
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comment = '#'
	lines, err := reader.ReadAll()
	if err != nil {
		return servers, errors.Wrap(err, "reader.readall")
	}
	for _, l := range lines {
		servers = append(servers, l[0])
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
