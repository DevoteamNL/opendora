package sql_client

const WEEKLY_DEPLOYMENT_SQL = `
with calendar_weeks as(
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

const MONTHLY_DEPLOYMENT_SQL = `
with _deployments as(
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
`
