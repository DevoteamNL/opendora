package sql_client

import (
	"database/sql"
	"devlake-go/group-sync/api/models"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectToDatabase() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "lake",
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
	fmt.Println("Connected!")
}

func queryDeployments(query string, args ...any) ([]models.DataPoint, error) {
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

func QueryTotalDeploymentsWeekly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	return queryDeployments(WEEKLY_DEPLOYMENT_SQL, to, to, from, to, projectName, from, to)
}

func QueryTotalDeploymentsMonthly(projectName string, from int64, to int64) ([]models.DataPoint, error) {
	return queryDeployments(MONTHLY_DEPLOYMENT_SQL, projectName, from, to)
}
