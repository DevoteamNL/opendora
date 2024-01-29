with _pr_stats as (
-- get the cycle time of PRs deployed by the deployments finished in the selected period
    SELECT
        distinct pr.id,
        ppm.pr_cycle_time
    FROM
        pull_requests pr
        join project_pr_metrics ppm on ppm.id = pr.id
        join project_mapping pm on pr.base_repo_id = pm.row_id
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
    WHEN median_change_lead_time < 60 then "lt-1hour"
    WHEN median_change_lead_time < 7 * 24 * 60 then "lt-1week"
    WHEN median_change_lead_time < 180 * 24 * 60 then "week-6month"
    WHEN median_change_lead_time >= 180 * 24 * 60 then "mt-6month"
    ELSE "N/A"
    END as median_change_lead_time
FROM _median_change_lead_time