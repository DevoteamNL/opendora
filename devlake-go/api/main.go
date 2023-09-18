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

type DataPoint struct {
	Key   string
	Value int
}

// TO TO FROM TO PROJECT FROM TO FROM TO
const WEEKLY_DEPLOYMENT_SQL = `
with calendar_weeks as(
-- construct the last few calendar months within the selected time period in the top-right corner
    SELECT CAST((FROM_UNIXTIME(?)-INTERVAL (T+U) WEEK) AS date) week
    FROM ( SELECT 0 T
            UNION ALL SELECT  10 UNION ALL SELECT  20 UNION ALL SELECT  30
            UNION ALL SELECT  40 UNION ALL SELECT  50
        ) T CROSS JOIN ( SELECT 0 U
            UNION ALL SELECT   1 UNION ALL SELECT   2 UNION ALL SELECT   3
            UNION ALL SELECT   4 UNION ALL SELECT   5 UNION ALL SELECT   6
            UNION ALL SELECT   7 UNION ALL SELECT   8 UNION ALL SELECT   9
        ) U
    WHERE
        (FROM_UNIXTIME(?)-INTERVAL (T+U) WEEK) BETWEEN FROM_UNIXTIME(?) AND FROM_UNIXTIME(?)
),
 _deployments as(
-- When deploying multiple commits in one pipeline, GitLab and BitBucket may generate more than one deployment. However, DevLake consider these deployments as ONE production deployment and use the last one's finished_date as the finished date.
    SELECT
        YEARWEEK(deployment_finished_date) as week,
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
    YEARWEEK(cw.week) as year_week,
    case when d.deployment_count is null then 0 else d.deployment_count end as deployment_count
FROM
    calendar_weeks cw
    LEFT JOIN _deployments d on YEARWEEK(cw.week) = d.week
	WHERE cw.week BETWEEN FROM_UNIXTIME(?) AND FROM_UNIXTIME(?)
	ORDER BY cw.week DESC
`

// projectname, from, to, from, to
const MONTHLY_DEPLOYMENT_SQL = `
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

func queryDeployments(query string, args ...any) ([]DataPoint, error) {
	var dataPoints []DataPoint

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dataPoint DataPoint
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

func dfTotalHandler(w http.ResponseWriter, queries url.Values) {
	aggregations, exists := queries["aggregation"]
	aggregation := "weekly"
	if exists && len(aggregations) > 0 {
		aggregation = aggregations[0]
	}

	to := time.Now().Unix()
	from := to - (60 * 60 * 24 * 30 * 6)

	switch aggregation {
	case "weekly":
		data, err := queryDeployments(WEEKLY_DEPLOYMENT_SQL, to, to, from, to, "my-project", from, to)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "data: %v", data)
	case "monthly":
		data, err := queryDeployments(MONTHLY_DEPLOYMENT_SQL, "my-project", from, to)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "data: %v", data)
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
