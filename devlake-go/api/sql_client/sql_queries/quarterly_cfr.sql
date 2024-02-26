
WITH RECURSIVE calendar_quarters AS (
    SELECT
        DATE_ADD(
            MAKEDATE(YEAR(FROM_UNIXTIME(:from)), 1),
            INTERVAL QUARTER(FROM_UNIXTIME(:from)) -1 QUARTER
        ) AS quarter_date
    UNION
    ALL
    SELECT
        DATE_ADD(quarter_date, INTERVAL 1 QUARTER)
    FROM
        calendar_quarters
    WHERE
        quarter_date < FROM_UNIXTIME(:to)
), _deployments AS (
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

_change_failure_rate_for_each_quarter AS (
    SELECT
        DATE_ADD(
            MAKEDATE(YEAR(deployment_finished_date), 1),
            INTERVAL QUARTER(deployment_finished_date) -1 QUARTER
        ) AS quarter_date,
        CASE
            WHEN count(deployment_id) IS NULL THEN NULL
            ELSE sum(has_incident)/count(deployment_id) END AS change_failure_rate
    FROM
        _failure_caused_by_deployments
    GROUP BY 1
)

SELECT
    cq.quarter_date AS data_key,
    cfr.change_failure_rate AS data_value
FROM
    calendar_quarters cq
    LEFT JOIN _change_failure_rate_for_each_quarter cfr ON cq.quarter_date = cfr.quarter_date
ORDER BY
cq.quarter_date DESC