package sql_client

import (
	"database/sql"
	"devlake-go/group-sync/api/models"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type ClientInterface interface {
	QueryTotalDeploymentsWeekly(projectName string, from int64, to int64) ([]models.DataPoint, error)
	QueryTotalDeploymentsMonthly(projectName string, from int64, to int64) ([]models.DataPoint, error)
}

type Client struct {
}

var db *sql.DB

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
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("connected to DevLake database")
}

func (client Client) queryDeployments(query string, args ...any) ([]models.DataPoint, error) {
	if db == nil {
		return nil, fmt.Errorf("first connect to the database before querying deployments")
	}
	var dataPoints []models.DataPoint
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dataPoint models.DataPoint
		if err := rows.Scan(&dataPoint.Key, &dataPoint.Value); err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, dataPoint)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dataPoints, nil
}

func (client Client) QueryTotalDeploymentsWeekly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	return client.queryDeployments(WEEKLY_DEPLOYMENT_SQL, to, to, from, to, projectName, from, to)
}

func (client Client) QueryTotalDeploymentsMonthly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	return client.queryDeployments(MONTHLY_DEPLOYMENT_SQL, projectName, from, to)
}
