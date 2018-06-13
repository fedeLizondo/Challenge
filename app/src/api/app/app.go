package app

import (
	"api/app/items"
	"database/sql"

	//	"os"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	// Needed to sql lite 3
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

var (
	r *gin.Engine
)

const (
	port string = ":8080"
)

// StartApp ...
func StartApp() {
	r = gin.Default()
	db := configDataBase()
	// data := []byte("state")
	// sessionName := "goquestsession"
	// google.Setup("http://localhost:8080/callback", "./client_secret.json", []string{"https://www.googleapis.com/auth/drive"}, data)
	//
	// r.Use(google.Session(sessionName))
	//
	// r.Use(google.Auth())
	items.Configure(r, db)
	items.ConfigureForFiles(r)
	r.Run(port)
}

func configDataBase() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", "user", "userpwd", "db", "db"))
	if err != nil {
		panic("Could not connect to the db")
	}
	for {
		err := db.Ping()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		// This is bad practice... You should create a schema.sql with all the definitions
		createTable(db)
		return db
	}

}

func createTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS items(
		id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		description VARCHAR(255)
	);`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}
