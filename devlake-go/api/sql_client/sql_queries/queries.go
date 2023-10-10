package sql_queries

import _ "embed"

//go:embed weekly_deployment_count.sql
var WeeklyDeploymentSql string

//go:embed monthly_deployment_count.sql
var MonthlyDeploymentSql string

//go:embed quarterly_deployment_count.sql
var QuarterlyDeploymentSql string
