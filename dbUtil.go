package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "fitComp"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func initDb() {

	db := dbConn()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (username varchar(20) PRIMARY KEY NOT NULL, password TEXT NOT NULL, authtype TEXT NOT NULL)")
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec()
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS measurements (id INTEGER PRIMARY KEY AUTO_INCREMENT, username TEXT, date DATE, waist INTEGER, weight INTEGER)")
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS prize (id INTEGER PRIMARY KEY AUTO_INCREMENT, prize INTEGER, increase INTEGER, username TEXT, date DATE)")
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS invites (id INTEGER PRIMARY KEY AUTO_INCREMENT, code TEXT, username TEXT)")
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec()

	// insert admin user

	if !userExists(user{
		Username: "admin",
	}) {
		if !addUser(user{
			Username: "admin",
			Password: "admin",
			AuthType: "admin",
		}) {
			fmt.Println("ERROR: Unable to initialize admin user")
			os.Exit(1)
		}
	}

	defer db.Close()
}
