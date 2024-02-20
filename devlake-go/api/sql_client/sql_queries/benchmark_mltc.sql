with _pr_stats AS (
-- get the cycle time of PRs deployed by the deployments finished in the selected period
    SELECT
        DISTINCT pr.id,
        ppm.pr_cycle_time
    FROM
        pull_requests pr
        JOIN project_pr_metrics ppm ON ppm.id = pr.id
        JOIN project_mapping pm ON pr.base_repo_id = pm.row_id
        JOIN cicd_deployment_commits cdc ON ppm.deployment_commit_id = cdc.id
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        and pr.merged_date IS NOT NULL
        and ppm.pr_cycle_time IS NOT NULL
        and cdc.finished_date BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
),

_median_change_lead_time_ranks AS(
    SELECT *, percent_rank() OVER(ORDER BY pr_cycle_time) AS ranks
    FROM _pr_stats
),

_median_change_lead_time AS(
-- use median PR cycle time as the median change lead time
    SELECT max(pr_cycle_time) AS median_change_lead_time
    FROM _median_change_lead_time_ranks
    WHERE ranks <= 0.5
)

SELECT
  CASE
    WHEN median_change_lead_time < 60 THEN "lt-1hour"
    WHEN median_change_lead_time < 7 * 24 * 60 THEN "lt-1week"
    WHEN median_change_lead_time < 180 * 24 * 60 THEN "week-6month"
    WHEN median_change_lead_time >= 180 * 24 * 60 THEN "mt-6month"
    ELSE "N/A"
    END AS median_change_lead_time
FROM _median_change_lead_time