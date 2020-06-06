package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/alexbrainman/odbc"
	"github.com/billgraziano/mssqlh"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: mssqlh.exe fqdn")
	}
	fqdn := os.Args[1]
	log.Printf("connecting to: %s...\r\n", fqdn)
	db, err := mssqlh.Open(fqdn)
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

	session, err := mssqlh.GetSession(context.Background(), db)
	if err != nil {
		log.Fatal(errors.Wrap(err, "mssqlh.getsession"))
	}
	log.Printf("Connected to %s (%s) on session %d via %s\r\n", server.Name, server.Domain, session.ID, session.AuthScheme)
}

func test() {
	s := mssqlh.NewConnection("ab.c.com", "", "", "myapp").String()
	fmt.Println(s)
}
