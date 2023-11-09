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
		User:   os.Getenv("DEVLAKE_DBUSER"),
		Passwd: os.Getenv("DEVLAKE_DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DEVLAKE_DBADDRESS"),
		DBName: os.Getenv("DEVLAKE_DBNAME"),
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

func (client Client) QueryDeployments(query string, params QueryParams) ([]models.DataPoint, error) {
	if db == nil {
		return nil, fmt.Errorf("first connect to the database before querying deployments")
	}
	var dataPoints []models.DataPoint
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
		dataPoint := models.DataPoint{}
		if err := rows.StructScan(&dataPoint); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, dataPoint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dataPoints, nil
}

func (client Client) QueryBenchmark(query string, params QueryParams) (string, error) {
	if db == nil {
		return "", fmt.Errorf("first connect to the database before querying benchmarks")
	}
	rows, err := db.NamedQuery(query, &params)
	if err != nil {
		return "", err
	}

	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	rows.Next()
	response := models.BenchmarkResponse{}
	if err := rows.StructScan(&response); err != nil {
		return "", err
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return response.Key, nil
}
