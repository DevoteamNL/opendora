
WITH last_few_calendar_months AS(
    SELECT CAST((SYSDATE()-INTERVAL (H+T+U) DAY) AS date) day
    FROM ( SELECT 0 H
            UNION ALL SELECT 100 UNION ALL SELECT 200 UNION ALL SELECT 300
        ) H CROSS JOIN ( SELECT 0 T
            UNION ALL SELECT  10 UNION ALL SELECT  20 UNION ALL SELECT  30
            UNION ALL SELECT  40 UNION ALL SELECT  50 UNION ALL SELECT  60
            UNION ALL SELECT  70 UNION ALL SELECT  80 UNION ALL SELECT  90
        ) T CROSS JOIN ( SELECT 0 U
            UNION ALL SELECT   1 UNION ALL SELECT   2 UNION ALL SELECT   3
            UNION ALL SELECT   4 UNION ALL SELECT   5 UNION ALL SELECT   6
            UNION ALL SELECT   7 UNION ALL SELECT   8 UNION ALL SELECT   9
        ) U
    WHERE
        (SYSDATE() - INTERVAL (H + T + U) DAY) > FROM_UNIXTIME(:from)
        AND (SYSDATE() - INTERVAL (H + T + U) DAY) < FROM_UNIXTIME(:to)
),

_production_deployment_days AS(
    SELECT
        cdc.cicd_deployment_id AS deployment_id,
        max(DATE(cdc.finished_date)) AS day
    FROM cicd_deployment_commits cdc
    JOIN project_mapping pm ON cdc.cicd_scope_id = pm.row_id AND pm.`table` = 'cicd_scopes'
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        and cdc.result = 'SUCCESS'
        and cdc.environment = 'PRODUCTION'
    GROUP BY 1
),

_days_weekly_deploy AS(
    SELECT
            date(DATE_ADD(last_few_calendar_months.day, INTERVAL -WEEKDAY(last_few_calendar_months.day) DAY)) as week,
            MAX(if(_production_deployment_days.day is not null, 1, null)) AS weeks_deployed,
            COUNT(distinct _production_deployment_days.day) AS days_deployed
    FROM 
        last_few_calendar_months
        LEFT JOIN _production_deployment_days ON _production_deployment_days.day = last_few_calendar_months.day
    GROUP BY week
    ),

_days_monthly_deploy AS(
    SELECT
            date(DATE_ADD(last_few_calendar_months.day, INTERVAL -DAY(last_few_calendar_months.day)+1 DAY)) as month,
            MAX(if(_production_deployment_days.day is not null, 1, null)) AS months_deployed,
          COUNT(distinct _production_deployment_days.day) AS days_deployed
    FROM 
        last_few_calendar_months
        LEFT JOIN _production_deployment_days ON _production_deployment_days.day = last_few_calendar_months.day
    GROUP BY month
    ),

_days_six_months_deploy AS (
  SELECT
    month,
    SUM(days_deployed) OVER (
      ORDER BY month
      ROWS BETWEEN 5 PRECEDING AND CURRENT ROW
    ) AS days_deployed_per_six_months,
    COUNT(months_deployed) OVER (
      ORDER BY month
      ROWS BETWEEN 5 PRECEDING AND CURRENT ROW
    ) AS months_deployed_count,
    ROW_NUMBER() OVER (
      PARTITION BY DATE_FORMAT(month, '%Y-%m') DIV 6
      ORDER BY month DESC
    ) AS rn
  FROM _days_monthly_deploy
),

_median_number_of_deployment_days_per_week_ranks AS(
    SELECT *, percent_rank() over(order by days_deployed) AS ranks
    FROM _days_weekly_deploy
),

_median_number_of_deployment_days_per_week AS(
    SELECT max(days_deployed) AS median_number_of_deployment_days_per_week
    FROM _median_number_of_deployment_days_per_week_ranks
    WHERE ranks <= 0.5
),

_median_number_of_deployment_days_per_month_ranks AS(
    SELECT *, percent_rank() over(order by days_deployed) AS ranks
    FROM _days_monthly_deploy
),

_median_number_of_deployment_days_per_month AS(
    SELECT max(days_deployed) AS median_number_of_deployment_days_per_month
    FROM _median_number_of_deployment_days_per_month_ranks
    WHERE ranks <= 0.5
),

_days_per_six_months_deploy_by_filter AS (
SELECT
  month,
  days_deployed_per_six_months,
  months_deployed_count
FROM _days_six_months_deploy
WHERE rn%6 = 1
),


_median_number_of_deployment_days_per_six_months_ranks AS(
    SELECT *, percent_rank() over(order by days_deployed_per_six_months) AS ranks
    FROM _days_per_six_months_deploy_by_filter
),

_median_number_of_deployment_days_per_six_months AS(
    SELECT min(days_deployed_per_six_months) AS median_number_of_deployment_days_per_six_months, min(months_deployed_count) as is_collected
    FROM _median_number_of_deployment_days_per_six_months_ranks
    WHERE ranks >= 0.5
)

SELECT 

  CASE
        WHEN median_number_of_deployment_days_per_week >= 7 THEN 'week-elite'
        WHEN median_number_of_deployment_days_per_week >= 1 THEN 'week-high'
        WHEN median_number_of_deployment_days_per_month >= 1 THEN 'month-medium'
        WHEN median_number_of_deployment_days_per_month < 1 and is_collected is not null THEN 'month-low'
        ELSE "na-deployments" END AS data_key,
  CASE
        WHEN median_number_of_deployment_days_per_week >= 7 THEN median_number_of_deployment_days_per_week
        WHEN median_number_of_deployment_days_per_week >= 1 THEN median_number_of_deployment_days_per_week
        WHEN median_number_of_deployment_days_per_month >= 1 THEN median_number_of_deployment_days_per_month
        WHEN median_number_of_deployment_days_per_month < 1 and is_collected is not null THEN median_number_of_deployment_days_per_month
        ELSE "" END AS data_value
FROM _median_number_of_deployment_days_per_week, _median_number_of_deployment_days_per_month, _median_number_of_deployment_days_per_six_months