package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

type Expense struct {
	Id     string   `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

var dbUrl string

func init() {
	dbUrl = os.Getenv("DATABASE_URL")
}

func CreateTable() {
	db := openDB()
	defer db.Close()

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
					id SERIAL PRIMARY KEY,
					title TEXT,
					amount FLOAT,
					note TEXT,
					tags TEXT[]
				);`

	_, err := db.Exec(createTb)
	if err != nil {
		log.Fatal("Can't create table", err)
	}

	log.Println("Create table success")
}

func CreateData(input Expense) (int, error) {
	db := openDB()
	defer db.Close()

	row := db.QueryRow("INSERT INTO expenses (title, amount, note , tags) values ($1, $2 , $3 , $4) RETURNING id",
		input.Title, input.Amount, input.Note, pq.Array(input.Tags))
	var id int

	err := row.Scan(&id)
	if err != nil {
		log.Println("Can't scan id", err)
		return -1, err
	}

	log.Println("Insert success id : ", id)

	return id, nil
}

func ReadData(id int) (Expense, error) {
	db := openDB()
	defer db.Close()

	result := Expense{}

	stmt, err := db.Prepare("SELECT id,	title, amount, note, tags  FROM expenses where id=$1")
	if err != nil {
		return result, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err != nil {
		return result, err
	}

	return result, nil
}

func openDB() *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Connect ro database error", err)
	}

	return db
}
