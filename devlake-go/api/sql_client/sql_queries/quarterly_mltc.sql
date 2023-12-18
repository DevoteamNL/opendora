
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
), _pr_stats as (
-- get the cycle time of PRs deployed by the deployments finished each quarter
    SELECT
        distinct pr.id,
        cdc.finished_date AS quarter,
        ppm.pr_cycle_time
    FROM
        pull_requests pr
        join project_pr_metrics ppm on ppm.id = pr.id
        join project_mapping pm on pr.base_repo_id = pm.row_id and pm.`table` = 'repos'
        join cicd_deployment_commits cdc on ppm.deployment_commit_id = cdc.id
        join repos ON cdc.repo_id = repos.id
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


_find_median_clt_each_quarter_ranks as(
    SELECT *, percent_rank() over(PARTITION BY quarter order by pr_cycle_time) as ranks
    FROM _pr_stats
),

_clt as(
    SELECT quarter, max(pr_cycle_time) as median_change_lead_time
    FROM _find_median_clt_each_quarter_ranks
    WHERE ranks <= 0.5
    group by quarter
)
    SELECT
        cq.quarter_date AS data_key,
        CASE
            WHEN _clt.median_change_lead_time IS NULL THEN 0
            ELSE _clt.median_change_lead_time
        END AS data_value
    FROM
        calendar_quarters cq
        LEFT JOIN _clt ON cq.quarter_date = _clt.quarter
    ORDER BY
        cq.quarter_date DESC
