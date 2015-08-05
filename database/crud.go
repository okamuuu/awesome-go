package main

/*
http://blog.yohei.org/go-mysql-crud/
*/

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: os.Getenv("") を const に代入できない理由
var DB_USER = os.Getenv("DB_USER")
var DB_PASSWORD = os.Getenv("DB_PASSWORD")
var DB_NAME = "test_crud"

func main() {
	fmt.Println("** mysql start ** ")
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", DB_USER, DB_PASSWORD, DB_NAME),
	)
	if err != nil {
		log.Fatal("db error: %v", err)
	}
	defer db.Close()

	truncate(db)
	read(db)
	insert(db)
	read(db)
	update(db, 1)
	read(db)
	delete(db, 1)
	read(db)

	fmt.Println("** mysql end ** ")
}

func truncate(db *sql.DB) {
	query := "TRUNCATE TABLE user"
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("truncate error: ", err)
	}
}

func read(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM `user`")
	defer rows.Close()
	if err != nil {
		log.Fatal("query error: ", err)
	}

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("scan error: ", err)
		}
		fmt.Println(id, name)
	}
}

func insert(db *sql.DB) {
	query := "INSERT INTO user (name) values (?)"
	result, err := db.Exec(query, "Coffee")
	if err != nil {
		log.Fatal("insert error: ", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		fmt.Println("insert last id: %d", lastId)
	}
}

func update(db *sql.DB, id int) {
	query := "UPDATE user SET name=? WHERE id=?"

	_, err := db.Exec(query, "poke", id)
	if err != nil {
		log.Fatal("update error: ", err)
	} else {
		fmt.Println("update complete! user id", id)
	}
}

func delete(db *sql.DB, id int) {
	query := "DELETE FROM user WHERE id=?"
	_, err := db.Exec(query, id)
	if err != nil {
		log.Fatal("delete error: ", err)
	} else {
		fmt.Println("delete complete! id = ", id)
	}
}
