with _pr_stats AS (
    SELECT
        DISTINCT pr.id,
        ppm.pr_cycle_time
    FROM
        pull_requests pr
        join project_pr_metrics ppm on ppm.id = pr.id
        join cicd_deployment_commits cdc on ppm.deployment_commit_id = cdc.id
        JOIN repos ON cdc.repo_id = repos.id
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
    SELECT max(pr_cycle_time) AS median_change_lead_time
    FROM _median_change_lead_time_ranks
    WHERE ranks <= 0.5
)

SELECT
  CASE
    WHEN median_change_lead_time < 60 then "lt-1hour"
    WHEN median_change_lead_time < 7 * 24 * 60 then "lt-1week"
    WHEN median_change_lead_time < 180 * 24 * 60 then "week-6month"
    WHEN median_change_lead_time >= 180 * 24 * 60 then "mt-6month"
    ELSE "mt-6month"
    END as data_key
FROM _median_change_lead_time