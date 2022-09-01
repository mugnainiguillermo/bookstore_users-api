package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	Client *sql.DB

	username = os.Getenv("MYSQL_DB_USERNAME")
	password = os.Getenv("MYSQL_DB_PASSWORD")
	host     = os.Getenv("MYSQL_DB_HOST")
	schema   = os.Getenv("MYSQL_DB_SCHEMA")
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, password, host, schema)
	var err error
	Client, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		// TODO: Proper error handling
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		// TODO: Proper error handling
		panic(err)
	}
	log.Println("database successfully configured")
}
