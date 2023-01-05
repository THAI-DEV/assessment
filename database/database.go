package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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

	createTb := `CREATE TABLE IF NOT EXISTS public.expenses (
					id SERIAL PRIMARY KEY,
					title TEXT,
					amount FLOAT,
					note TEXT,
					tags TEXT[]
				);`

	_, err := db.Exec(createTb)
	if err != nil {
		log.Println("Can't create table", err)
	}

	log.Println("Create table success")
}

func CreateData(input Expense) (int, error) {
	db := openDB()
	defer db.Close()

	row := db.QueryRow("INSERT INTO public.expenses (title, amount, note , tags) values ($1, $2 , $3 , $4) RETURNING id",
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

func UpdateData(input Expense) (int, error) {
	db := openDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE public.expenses SET title=$2, amount=$3, note=$4 , tags=$5 WHERE id=$1")
	if err != nil {
		log.Println("Can't prepare statment update", err)
	}

	if _, err := stmt.Exec(input.Id, input.Title, input.Amount, input.Note, pq.Array(input.Tags)); err != nil {
		log.Panicln("Error execute update ", err)
		return -1, err
	}

	fmt.Println("Update success id : ", input.Id)

	id, _ := strconv.Atoi(input.Id)
	return id, nil
}

func ReadData(id int) (Expense, error) {
	db := openDB()
	defer db.Close()

	result := Expense{}

	stmt, err := db.Prepare("SELECT id,	title, amount, note, tags  FROM public.expenses where id=$1")
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

func ReadAllData() ([]Expense, error) {
	db := openDB()
	defer db.Close()

	result := []Expense{}

	stmt, err := db.Prepare("SELECT id,	title, amount, note, tags  FROM public.expenses")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Println("Can't query all expenses", err)
		return nil, err
	}

	for rows.Next() {
		output := Expense{}

		err := rows.Scan(&output.Id, &output.Title, &output.Amount, &output.Note, pq.Array(&output.Tags))
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}

		result = append(result, output)
	}

	return result, nil
}

func openDB() *sql.DB {
	fmt.Println("Connect db : ", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Connect ro database error", err)
	}

	return db
}
