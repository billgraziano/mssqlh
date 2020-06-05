package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/alexbrainman/odbc"
	"github.com/billgraziano/mssqlh"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

func main() {
	log.Print("starting mssqlh.exe")
	db, err := mssqlh.Open("D40\\SQL2016")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	session, err := mssqlh.GetSession(context.Background(), db)
	if err != nil {
		log.Fatal(errors.Wrap(err, "mssqlh.getsession"))
	}
	log.Println("server: ", session.ServerName)
}

func test() {
	s := mssqlh.NewConnection("ab.c.com", "", "", "myapp").String()
	fmt.Println(s)
}
