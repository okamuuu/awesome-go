package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/shogo82148/go-sql-proxy"
)

func main() {
	sql.Register("sqlite3-proxy", proxy.NewProxy(&sqlite3.SQLiteDriver{}, &proxy.Hooks{
		Open: func(_ interface{}, conn driver.Conn) error {
			log.Println("Open")
			return nil
		},
		Exec: func(_ interface{}, stmt *proxy.Stmt, args []driver.Value, result driver.Result) error {
			log.Printf("Exec: %s; args = %v\n", stmt.QueryString, args)
			return nil
		},
		Query: func(_ interface{}, stmt *proxy.Stmt, args []driver.Value, rows driver.Rows) error {
			log.Printf("Query: %s; args = %v\n", stmt.QueryString, args)
			return nil
		},
		Begin: func(_ interface{}, conn *proxy.Conn) error {
			log.Println("Begin")
			return nil
		},
		Commit: func(_ interface{}, tx *proxy.Tx) error {
			log.Println("Commit")
			return nil
		},
		Rollback: func(_ interface{}, tx *proxy.Tx) error {
			log.Println("Rollback")
			return nil
		},
	}))

	db, err := sql.Open("sqlite3-proxy", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 4; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
}
