
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
), _deployments as (
    SELECT
        cdc.cicd_deployment_id as deployment_id,
        max(cdc.finished_date) as deployment_finished_date
    FROM
        cicd_deployment_commits cdc
        JOIN project_mapping pm on cdc.cicd_scope_id = pm.row_id and pm.`table` = 'cicd_scopes'
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        and cdc.result = 'SUCCESS'
        and cdc.environment = 'PRODUCTION'
    GROUP BY 1
),

_failure_caused_by_deployments as (
    SELECT
        d.deployment_id,
        d.deployment_finished_date,
        count(distinct case when i.type = 'INCIDENT' then d.deployment_id else null end) as has_incident
    FROM
        _deployments d
        left join project_issue_metrics pim on d.deployment_id = pim.deployment_id
        left join issues i on pim.id = i.id
    GROUP BY 1,2
),

_change_failure_rate_for_each_quarter as (
    SELECT
        DATE_ADD(
            MAKEDATE(YEAR(deployment_finished_date), 1),
            INTERVAL QUARTER(deployment_finished_date) -1 QUARTER
        ) AS quarter_date,
        case
            when count(deployment_id) is null then null
            else sum(has_incident)/count(deployment_id) end as change_failure_rate
    FROM
        _failure_caused_by_deployments
    GROUP BY 1
)

SELECT
    cq.quarter_date AS data_key,
    cfr.change_failure_rate AS data_value
FROM
    calendar_quarters cq
    LEFT JOIN _change_failure_rate_for_each_quarter cfr on cq.quarter_date = cfr.quarter_date
ORDER BY
cq.quarter_date DESC