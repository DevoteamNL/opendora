WITH RECURSIVE calendar_weeks AS (
    SELECT
        STR_TO_DATE(
            CONCAT(YEARWEEK(FROM_UNIXTIME(:from)), ' Sunday'),
            '%X%V %W'
        ) AS week_date
    UNION
    ALL
    SELECT
        DATE_ADD(week_date, INTERVAL 1 WEEK)
    FROM
        calendar_weeks
    WHERE
        week_date < FROM_UNIXTIME(:to)
),

_deployments AS (
   SELECT
        cdc.cicd_deployment_id AS deployment_id,
        max(cdc.finished_date) AS deployment_finished_date
    FROM
        cicd_deployment_commits cdc
        JOIN project_mapping pm ON cdc.cicd_scope_id = pm.row_id AND pm.`table` = 'cicd_scopes'
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        AND cdc.result = 'SUCCESS'
        AND cdc.environment = 'PRODUCTION'
    GROUP BY 1
),

_failure_caused_by_deployments AS (
   SELECT
        d.deployment_id,
        d.deployment_finished_date,
        count(DISTINCT CASE WHEN i.type = 'INCIDENT' THEN d.deployment_id ELSE NULL END) AS has_incident
    FROM
        _deployments d
        LEFT JOIN project_issue_metrics pim ON d.deployment_id = pim.deployment_id
        LEFT JOIN issues i ON pim.id = i.id
    GROUP BY 1,2
),

_change_failure_rate_for_each_week as (
    SELECT
        YEARWEEK(deployment_finished_date) AS week,
        case
            WHERE count(deployment_id) IS NULL THEN NULL
            ELSE sum(has_incident)/count(deployment_id) END AS change_failure_rate
    FROM
        _failure_caused_by_deployments
    GROUP BY 1
)

    SELECT
        YEARWEEK(cw.week_date) AS data_key,
		cfr.change_failure_rate AS data_value
    FROM
        calendar_weeks cw
        LEFT JOIN _change_failure_rate_for_each_week cfr ON YEARWEEK(cw.week_date) = cfr.week
    ORDER BY
        cw.week_date DESC

