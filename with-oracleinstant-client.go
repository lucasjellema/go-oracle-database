package main

import (
	"database/sql"
	"fmt"

	_ "github.com/godror/godror"
)

func doDBThingsThroughInstantClient(dbParams map[string]string) {
	var db *sql.DB
	var err error
	if val, ok := dbParams["walletLocation"]; ok && val != "" {
		db, err = sql.Open("godror", fmt.Sprintf(`user="%s" password="%s"
		connectString="tcps://%s:%s/%s?wallet_location=%s"
		   `, dbParams["username"], dbParams["password"], dbParams["server"], dbParams["port"], dbParams["service"], dbParams["walletLocation"]))
	}
	if val, ok := dbParams["walletLocation"]; !ok || val == "" {
		connectionString := "oracle://" + dbParams["username"] + ":" + dbParams["password"] + "@" + dbParams["server"] + ":" + dbParams["port"] + "/" + dbParams["service"]
		db, err = sql.Open("oracle", connectionString)
	}

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
