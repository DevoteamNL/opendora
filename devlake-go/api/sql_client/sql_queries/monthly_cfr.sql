with _deployments AS (
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

_change_failure_rate_for_each_month as (
    SELECT
        date_format(deployment_finished_date,'%y/%m') as month,
        case
            when count(deployment_id) is null then null
            else sum(has_incident)/count(deployment_id) end as change_failure_rate
    FROM
        _failure_caused_by_deployments
    GROUP BY 1
)

SELECT
    cm.month,
    cfr.change_failure_rate
FROM
    calendar_months cm
    LEFT JOIN _change_failure_rate_for_each_month cfr on cm.month = cfr.month
    WHERE cm.month_timestamp BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)