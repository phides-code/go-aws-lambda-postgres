package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"inchworm/database"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var logger *zap.Logger
var db *sql.DB

func init() {
	l, _ := zap.NewProduction()
	logger = l

	dbConnection, err := database.GetConnection()

	if err != nil {
		logger.Error("Error connecting to database", zap.Error(err))
		panic(err)
	}

	dbConnection.Ping()

	if err != nil {
		logger.Error("Error pinging database", zap.Error(err))
		panic(err)
	}

	db = dbConnection

	defer logger.Sync() // flush buffer, if any
}

type DefaultResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetEmployeesResponse struct {
	Employees []*database.Employee `json:"employees"`
}

func MyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	var res *events.APIGatewayProxyResponse

	logger.Info("received event", zap.Any("method", event.HTTPMethod),
		zap.Any("path", event.Path),
		zap.Any("body", event.Body),
	)

	if event.Path == "/migrate" {
		err := database.CreateEmployeesTable(ctx, db)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  fmt.Sprint(http.StatusInternalServerError),
				Message: "could not create employees table",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&DefaultResponse{
			Status:  fmt.Sprint(http.StatusOK),
			Message: "ran migrations!",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else if event.Path == "/employees" && event.HTTPMethod == http.MethodPost {
		// create a new employeE
		employee := &database.Employee{}
		err := json.Unmarshal([]byte(event.Body), &employee)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  fmt.Sprint(http.StatusBadRequest),
				Message: err.Error(),
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       string(body),
			}, nil
		}

		log.Println()

		err = database.CreateEmployee(ctx, db, employee.Email, employee.FirstName, employee.LastName)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  fmt.Sprint(http.StatusInternalServerError),
				Message: err.Error(),
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&DefaultResponse{
			Status:  fmt.Sprint(http.StatusOK),
			Message: "success!",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else if event.Path == "/employees" && event.HTTPMethod == http.MethodGet {
		// get all employees
		employees, err := database.GetEmployees(ctx, db)
		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  fmt.Sprint(http.StatusInternalServerError),
				Message: "error fetching employees",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&GetEmployeesResponse{
			Employees: employees,
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else {
		body, _ := json.Marshal(&DefaultResponse{
			Status:  fmt.Sprint(http.StatusOK),
			Message: "default path",
		})

		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	}

	return res, nil
}

func main() {
	lambda.Start(MyHandler)
}
