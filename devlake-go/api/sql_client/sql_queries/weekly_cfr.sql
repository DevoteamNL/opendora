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

_deployments as (
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

_change_failure_rate_for_each_week as (
    SELECT
        YEARWEEK(deployment_finished_date) AS week,
        case
            when count(deployment_id) is null then null
            else sum(has_incident)/count(deployment_id) end as change_failure_rate
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

