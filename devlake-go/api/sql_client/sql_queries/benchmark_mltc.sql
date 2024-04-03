WITH _pr_stats AS (
    SELECT
        DISTINCT pr.id,
        ppm.pr_cycle_time
    FROM
        pull_requests pr 
        JOIN project_pr_metrics ppm ON ppm.id = pr.id
        JOIN project_mapping pm ON pr.base_repo_id = pm.row_id AND pm.`table` = 'repos'
        JOIN cicd_deployment_commits cdc ON ppm.deployment_commit_id = cdc.id
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        AND pr.merged_date IS NOT NULL
        AND ppm.pr_cycle_time IS NOT NULL
        and cdc.finished_date BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
),

_median_change_lead_time_ranks as(
    SELECT *, percent_rank() over(order by pr_cycle_time) as ranks
    FROM _pr_stats
),

_median_change_lead_time as(
-- use median PR cycle time as the median change lead time
    SELECT max(pr_cycle_time) as median_change_lead_time
    FROM _median_change_lead_time_ranks
    WHERE ranks <= 0.5
)

SELECT 
  CASE
        WHEN median_change_lead_time < 24 * 60 THEN "elite"
        WHEN median_change_lead_time < 7 * 24 * 60 THEN "high"
        WHEN median_change_lead_time < 30 * 24 * 60 THEN "medium"
        WHEN median_change_lead_time >= 30 * 24 * 60 THEN "low"
        ELSE "na-deployments-pullrequests"
        END AS data_key,
  CASE
        WHEN median_change_lead_time > 0 THEN round(median_change_lead_time/60,1)
        ELSE ""
        END AS data_value
FROM _median_change_lead_time