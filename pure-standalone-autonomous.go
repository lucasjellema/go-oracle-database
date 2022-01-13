package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

func Header(columns []string) {

}

func Record(columns []string, values []driver.Value) {
	for i, c := range values {
		fmt.Printf("%-25s: %v\n", columns[i], c)
	}
	fmt.Println()
}
func handleQuery(query string, conn *go_ora.Connection) {
	stmt := go_ora.NewStmt(query, conn)

	defer stmt.Close()

	rows, err := stmt.Query(nil)
	handleError("Can't query", err)
	defer rows.Close()

	columns := rows.Columns()

	values := make([]driver.Value, len(columns))

	Header(columns)
	for {
		err = rows.Next(values)
		if err != nil {
			break
		}
		Record(columns, values)
	}
	if err != io.EOF {
		handleError("Can't Next", err)
	}
}

func queryData(conn *go_ora.Connection) error {
	t := time.Now()

	stmt := go_ora.NewStmt("SELECT SYSTIMESTAMP FROM DUAL", conn)
	rows, err := stmt.Query_(nil)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()
	var Date time.Time
	for rows.Next_() {
		err = rows.Scan(&Date)
		if err != nil {
			return err
		}
		fmt.Println("Current Timestamp in Database: ", Date)
	}
	fmt.Println("Finish query rows: ", time.Now().Sub(t))
	return nil
}

func executeStatement(sqlText string, conn *go_ora.Connection) {
	stmt := go_ora.NewStmt(sqlText, conn)
	var err error
	defer func() {
		err = stmt.Close()
		if err != nil {
			fmt.Println("Can't close stmt", err)
		}
	}()

	result, err := stmt.Exec(nil)
	if err != nil {
		fmt.Println("Can't execute sql", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Println("Rows affected: ", rowsAffected)
}

func main() {
	service := "k8j2fvxbaujdcfy_daprdb_low.adb.oraclecloud.com"
	user := "demo"
	server := "adb.us-ashburn-1.oraclecloud.com"
	port := 1522
	pw := "Modem123mode"
	walletLocation := "/home/lucas/dapr-work/components-contrib/state/oracledatabase/Wallet_daprDB/"

	urlOptions := map[string]string{
		"TRACE FILE": "trace.log",
		"SSL":        "enable",
		"SSL Verify": "false",
		"WALLET":     walletLocation,
	}
	databaseURL := go_ora.BuildUrl(server, port, service, user, pw, urlOptions)
	conn, err := go_ora.NewConnection(databaseURL)
	if err != nil {
		fmt.Println("Can't create connection: ", err)
		return
	}
	err = conn.Open()
	if err != nil {
		fmt.Println("Can't open connection: ", err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()
	var ctx = context.Background()
	err = conn.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	err = queryData(conn)
	if err != nil {
		panic(fmt.Errorf("error querying db: %w", err))
	}
	handleQuery("select to_char(systimestamp,'DD-MM-YYYY HH24:MI:SS') as the_time from dual", conn)
	handleQuery("select 42 as magic_number from dual", conn)

	executeStatement(createTableStatement, conn)
	executeStatement("INSERT INTO TEMP_TABLE ( NAME , VALUE) VALUES ('Johnny',391)", conn)
	executeStatement("INSERT INTO TEMP_TABLE ( NAME , VALUE) VALUES ('Marianne',312)", conn)
	handleQuery("select name, creation_time, value from TEMP_TABLE", conn)

	executeStatement(dropTableStatement, conn)
}
