
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
    SELECT
        DISTINCT pr.id,
        cdc.finished_date AS quarter,
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


_find_median_clt_each_quarter_ranks AS(
    SELECT *, percent_rank() OVER(PARTITION BY quarter ORDER BY pr_cycle_time) as ranks
    FROM _pr_stats
),

_clt as(
    SELECT quarter, max(pr_cycle_time) AS median_change_lead_time
    FROM _find_median_clt_each_quarter_ranks
    WHERE ranks <= 0.5
    GROUP BY quarter
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
