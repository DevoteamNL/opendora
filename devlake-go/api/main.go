package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectToDatabase() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "lake",
	}
	// Get a database handle.
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

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

const DEPLOYMENT_SQL = `
-- Metric 1: Number of deployments per month
with _deployments as(
-- When deploying multiple commits in one pipeline, GitLab and BitBucket may generate more than one deployment. However, DevLake consider these deployments as ONE production deployment and use the last one's finished_date as the finished date.
    SELECT
        date_format(deployment_finished_date,'%y/%m') as month,
        count(cicd_deployment_id) as deployment_count
    FROM (
        SELECT
            cdc.cicd_deployment_id,
            max(cdc.finished_date) as deployment_finished_date
        FROM cicd_deployment_commits cdc
        JOIN project_mapping pm on cdc.cicd_scope_id = pm.row_id and pm.` + "`table`" + ` = 'cicd_scopes'
        WHERE
            pm.project_name = ?
            and cdc.result = 'SUCCESS'
            and cdc.environment = 'PRODUCTION'
        GROUP BY 1
		-- WHERE max(cdc.finished_date) BETWEEN FROM_UNIXTIME(?) AND FROM_UNIXTIME(?)
    ) _production_deployments
    GROUP BY 1
)

SELECT
    cm.month,
    case when d.deployment_count is null then 0 else d.deployment_count end as deployment_count
FROM
    calendar_months cm
    LEFT JOIN _deployments d on cm.month = d.month
	WHERE cm.month_timestamp BETWEEN FROM_UNIXTIME(?) AND FROM_UNIXTIME(?)
	-- LIMIT 10,20
`

func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	// projectname, from, to, from, to
	to := time.Now().Unix()
	from := to - (60 * 60 * 24 * 30 * 6)
	fmt.Printf("to: %v from: %v\n", to, from)
	rows, err := db.Query(DEPLOYMENT_SQL, "my-project", from, to)
	fmt.Printf("query sql")
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	fmt.Printf("loop sql")
	for rows.Next() {
		var month string
		var count int

		// var alb Album
		if err := rows.Scan(&month, &count); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		fmt.Printf("month: %v count: %v\n", month, count)
		// albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func dfTotalHandler(w http.ResponseWriter, queries url.Values) {
	aggregations, exists := queries["aggregation"]
	aggregation := "weekly"
	if exists && len(aggregations) > 0 {
		aggregation = aggregations[0]
	}
	switch aggregation {
	case "weekly":
		fmt.Fprintf(w, "Hello, %q", html.EscapeString("weekly"))
	case "monthly":
		_, err := albumsByArtist("")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Hello, %q", html.EscapeString("monthly"))
	case "quarterly":
		fmt.Fprintf(w, "Hello, %q", html.EscapeString("quarterly"))
	default:
		http.Error(w, "aggregation should be provided as either weekly, monthly or quarterly", http.StatusBadRequest)
	}
}

func main() {
	connectToDatabase()
	http.HandleFunc("/dora/api/metric", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		metricTypes, exists := queries["type"]
		if !exists || len(metricTypes) != 1 {
			http.Error(w, "type should be provided as either df_average or df_total", http.StatusBadRequest)
			return
		}
		switch metricTypes[0] {
		case "df_total":
			dfTotalHandler(w, queries)
		case "df_average":
			fmt.Fprintf(w, "Hello, %q", html.EscapeString("average"))
		default:
			http.Error(w, "type should be provided as either df_average or df_total", http.StatusBadRequest)
		}
	})

	log.Fatal(http.ListenAndServe(":10666", nil))
}
