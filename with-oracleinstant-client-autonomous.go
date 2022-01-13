package main

import (
	"database/sql"
	"fmt"

	_ "github.com/godror/godror"
)

func main4() {
	service := "k8j2fvxbaujdcfy_daprdb_low.adb.oraclecloud.com"
	user := "demo"
	server := "adb.us-ashburn-1.oraclecloud.com"
	port := "1522"
	pw := "Modem123mode"
	walletLocation := "/home/lucas/dapr-work/components-contrib/state/oracledatabase/Wallet_daprDB/"
	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s"
	connectString="tcps://%s:%s/%s?wallet_location=%s"
	   `, user, pw, server, port, service, walletLocation))

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

	var queryResultColumnOne string
	row := db.QueryRow("SELECT systimestamp FROM dual")
	err = row.Scan(&queryResultColumnOne)
	if err != nil {
		panic(fmt.Errorf("error scanning db: %w", err))
	}
	fmt.Println("The time in the database ", queryResultColumnOne)

	someAdditionalActions(db)
}
