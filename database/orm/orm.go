package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/k0kubun/pp"
)

var DB_USER = os.Getenv("DB_USER")
var DB_PASSWORD = os.Getenv("DB_PASSWORD")
var DB_NAME = "test_orm"

type Post struct {
	Id        int64
	Title     string
	Body      string
	CreatedAt int64
	UpdatedAt int64
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", DB_USER, DB_PASSWORD, DB_NAME),
	)
	checkErr(err, "sql.Open failed")
	defer db.Close()

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	t := dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	t.ColMap("Id").Rename("id")
	t.ColMap("Title").Rename("title")
	t.ColMap("Body").Rename("body")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("UpdatedAt").Rename("updated_at")

	// create the table
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	post := &Post{
		Title:     "title1",
		Body:      "body1",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = dbmap.Insert(post)
	checkErr(err, "Insert failed")

	primaryKey := 1
	p1, err := dbmap.Get(Post{}, primaryKey)

	pp.Println(p1)
	p := p1.(*Post)
	fmt.Println(p.Title)
	fmt.Println(p.Body)

}
