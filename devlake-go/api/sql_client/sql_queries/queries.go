package sql_queries

import _ "embed"

//go:embed weekly_deployment_count.sql
var WEEKLY_DEPLOYMENT_SQL string

//go:embed monthly_deployment_count.sql
var MONTHLY_DEPLOYMENT_SQL string

//go:embed quarterly_deployment_count.sql
var QUARTERLY_DEPLOYMENT_SQL string
