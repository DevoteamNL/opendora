package sql_client

import (
	"fmt"
	"log"
	"os"

	"github.com/devoteamnl/opendora/api/models"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ClientInterface interface {
	QueryDeployments(query string, params QueryParams) ([]models.DataPoint, error)
	QueryBenchmark(query string, params QueryParams) (string, error)
}

type Client struct {
}

var db *sqlx.DB

func New() Client {
	client := Client{}
	client.connectToDatabase()
	return client
}

func (client Client) connectToDatabase() {
	cfg := mysql.Config{
		User:                 os.Getenv("DEVLAKE_DBUSER"),
		Passwd:               os.Getenv("DEVLAKE_DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DEVLAKE_DBADDRESS"),
		DBName:               os.Getenv("DEVLAKE_DBNAME"),
		AllowNativePasswords: true,
	}

	var err error
	db, err = sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("connected to DevLake database")
}

type QueryParams struct {
	To      int64
	From    int64
	Project string
}

func queryRows[R models.DataPoint | models.BenchmarkResponse](query string, params QueryParams, makeStruct func() R) ([]R, error) {
	if db == nil {
		return nil, fmt.Errorf("first connect to the database before querying")
	}
	var responseRows []R
	rows, err := db.NamedQuery(query, &params)
	if err != nil {
		return nil, err
	}

	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		dataPoint := makeStruct()
		if err := rows.StructScan(&dataPoint); err != nil {
			return nil, err
		}
		responseRows = append(responseRows, dataPoint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return responseRows, nil
}

func (client Client) QueryDeployments(query string, params QueryParams) ([]models.DataPoint, error) {
	return queryRows(query, params, func() models.DataPoint { return models.DataPoint{} })
}

func (client Client) QueryBenchmark(query string, params QueryParams) (string, error) {
	response, err := queryRows(query, params, func() models.BenchmarkResponse { return models.BenchmarkResponse{} })

	if err != nil {
		return "", err
	}

	return response[0].Key, nil
}
