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
), _pr_stats as (
-- get the cycle time of PRs deployed by the deployments finished each month
    SELECT
        distinct pr.id,
        YEARWEEK(cdc.finished_date) AS week,
        ppm.pr_cycle_time
    FROM
        pull_requests pr
        join project_pr_metrics ppm on ppm.id = pr.id
        join project_mapping pm on pr.base_repo_id = pm.row_id and pm.`table` = 'repos'
        join cicd_deployment_commits cdc on ppm.deployment_commit_id = cdc.id
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        and pr.merged_date is not null
        and ppm.pr_cycle_time is not null
        and cdc.finished_date BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
),


_find_median_clt_each_month_ranks as(
    SELECT *, percent_rank() over(PARTITION BY week order by pr_cycle_time) as ranks
    FROM _pr_stats
),

_clt as(
    SELECT week, max(pr_cycle_time) as median_change_lead_time
    FROM _find_median_clt_each_month_ranks
    WHERE ranks <= 0.5
    group by week
)
    SELECT
        YEARWEEK(cw.week_date) AS data_key,
        CASE
            WHEN _clt.median_change_lead_time IS NULL THEN 0
            ELSE _clt.median_change_lead_time
        END AS data_value
    FROM
        calendar_weeks cw
        LEFT JOIN _clt ON YEARWEEK(cw.week_date) = _clt.week
    ORDER BY
        cw.week_date DESC
