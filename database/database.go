package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func CreateTable(dbUrl string) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Connect ro database error", err)
	}

	defer db.Close()

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
					id SERIAL PRIMARY KEY,
					title TEXT,
					amount FLOAT,
					note TEXT,
					tags TEXT[]
				);`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	log.Println("create table success")
}
