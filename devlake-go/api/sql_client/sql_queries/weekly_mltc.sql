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
-- get the cycle time of PRs deployed by the deployments finished each week
    SELECT
        DISTINCT pr.id,
        YEARWEEK(cdc.finished_date) AS week,
        ppm.pr_cycle_time
    FROM
        pull_requests pr
        JOIN project_pr_metrics ppm ON ppm.id = pr.id
        JOIN project_mapping pm ON pr.base_repo_id = pm.row_id AND pm.`table` = 'repos'
        JOIN cicd_deployment_commits cdc ON ppm.deployment_commit_id = cdc.id
        JOIN repos ON cdc.repo_id = repos.id
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        AND pr.merged_date IS NOT NULL
        AND ppm.pr_cycle_time IS NOT NULL
        AND cdc.finished_date BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
),


_find_median_clt_each_week_ranks AS(
    SELECT *, percent_rank() over(PARTITION BY week ORDER BY pr_cycle_time) AS ranks
    FROM _pr_stats
),

_clt as(
    SELECT week, max(pr_cycle_time) AS median_change_lead_time
    FROM _find_median_clt_each_week_ranks
    WHERE ranks <= 0.5
    GROUP BY week
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
