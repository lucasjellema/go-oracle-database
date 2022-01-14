package main

// this example handles both a traditional database connection based on host, port, service name and username/password
// as well as a connection also based on an Oracle Wallet (that in addition to the traditional parameters requires the absolute wallet location)
// note that some of the data and functions used in this example are defined in the file common.go

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/sijms/go-ora/v2"
)

func doDBThings(dbParams map[string]string) {
	connectionString := "oracle://" + dbParams["username"] + ":" + dbParams["password"] + "@" + dbParams["server"] + ":" + dbParams["port"] + "/" + dbParams["service"]
	if val, ok := dbParams["walletLocation"]; ok && val != "" {
		connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(dbParams["walletLocation"])
	}
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}

	someAdditionalActions(db)
}
