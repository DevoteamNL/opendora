with _pr_stats as (
    SELECT
        distinct pr.id,
        date_format(cdc.finished_date,'%y/%m') as month,
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

_find_median_clt_each_month_ranks as(
    SELECT *, percent_rank() over(PARTITION BY month order by pr_cycle_time) as ranks
    FROM _pr_stats
),

_clt as(
    SELECT month, max(pr_cycle_time) as median_change_lead_time
    FROM _find_median_clt_each_month_ranks
    WHERE ranks <= 0.5
    group by month
)

SELECT
    cm.month as data_key,
    case
        when _clt.median_change_lead_time is null then 0
        else _clt.median_change_lead_time/60 
    end as data_value
FROM
    calendar_months cm
    LEFT JOIN _clt on cm.month = _clt.month
  WHERE cm.month_timestamp BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)