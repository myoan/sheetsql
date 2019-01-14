package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/myoan/sheetsql"
)

func main() {
	var (
		k = flag.String("key", "", "google service account key path")
		i = flag.String("id", "", "google sheet ID")
	)
	flag.Parse()
	if len(*k) == 0 || len(*i) == 0 {
		fmt.Println("Service Account Key or SheetID are undefined")
		os.Exit(1)
	}
	dsn := *k + "|" + *i
	conn, _ := sql.Open("sheetsql", dsn)
	stmt, _ := conn.Query("SELECT * FROM Sheet1")
	for stmt.Next() {
		var id, name, age string
		stmt.Scan(&id, &name, &age)
		fmt.Println("-- data")
		fmt.Printf("id: %v\n", id)
		fmt.Printf("name: %v\n", name)
		fmt.Printf("age: %+v\n", age)
	}
}
