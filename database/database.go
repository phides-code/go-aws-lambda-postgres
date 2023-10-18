package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("host")
	password = os.Getenv("password")
	port     = 5432
	user     = "postgres"
	database = "postgres"
)

func GetConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, database)

	return sql.Open("postgres", psqlInfo)
}

const createEmployeesTableSQL = `
	CREATE TABLE employees (
		id serial,
		email varchar,
		first_name varchar,
		last_name varchar
	);
`

func CreateEmployeesTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createEmployeesTableSQL)
	return err
}
