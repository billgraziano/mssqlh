package main

import (
	"context"
	"flag"
	"log"

	"github.com/billgraziano/mssqlh"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

func main() {
	var srv = flag.String("s", "", "server (server.net.com or server,port or server:port)")
	var user = flag.String("u", "", "user name")
	var pwd = flag.String("p", "", "password")
	var dbname = flag.String("d", "", "database name")
	flag.Parse()
	var cxn mssqlh.Connection
	if srv != nil {
		cxn.FQDN = *srv
	}
	if user != nil {
		cxn.User = *user
	}
	if pwd != nil {
		cxn.Password = *pwd
	}
	if dbname != nil {
		cxn.Database = *dbname
	}

	log.Printf("connecting to: '%s'...\r\n", cxn.ServerName())
	db, err := cxn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	server, err := mssqlh.GetServer(nil, db)
	if err != nil {
		log.Fatal(errors.Wrap(err, "mssqlh.getserver"))
	}

	//session := mssqlh.Session{}
	session, err := mssqlh.GetSession(context.Background(), db)
	if err != nil {
		log.Fatal(errors.Wrap(err, "mssqlh.getsession"))
	}
	log.Printf("Connected to %s (%s) as '%s' on session %d via '%s' in [%s] using '%s'\r\n",
		server.Name, server.Domain, session.Login, session.ID, session.AuthScheme, session.Database, session.Application)
}
