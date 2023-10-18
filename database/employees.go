package database

import (
	"context"
	"database/sql"
)

const createEmployeeSQL = `
	INSERT INTO
		employees (email, first_name, last_name)
	VALUES ($1, $2, $3);
`

func CreateEmployee(ctx context.Context, db *sql.DB, email string, firstName string, lastName string) error {
	_, err := db.ExecContext(ctx, createEmployeeSQL, email, firstName, lastName)

	return err
}

const getEmployeesSQL = `
	SELECT * from employees;
`

type Employee struct {
	ID        int
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func GetEmployees(ctx context.Context, db *sql.DB) ([]*Employee, error) {

	rows, err := db.QueryContext(ctx, getEmployeesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		employee := &Employee{}

		if err := rows.Scan(&employee.ID, &employee.Email, &employee.FirstName, &employee.LastName); err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
