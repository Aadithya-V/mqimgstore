package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var (
	MySQLAddr = "127.0.0.1:3306"
)

func NewMySQLSession() *sql.DB {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "go",  //os.Getenv("DBUSER"),
		Passwd: "123", //os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   MySQLAddr,
		DBName: "test",
	}
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("MySQL Connected!")

	return db
}
