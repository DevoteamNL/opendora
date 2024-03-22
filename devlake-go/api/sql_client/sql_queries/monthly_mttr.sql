WITH _incidents AS (
    SELECT
      DISTINCT i.id,
        DATE_FORMAT(i.created_date,'%y/%m') AS month,
        cast(lead_time_minutes AS signed) AS lead_time_minutes
    FROM
        issues i
      JOIN board_issues bi ON i.id = bi.issue_id
      JOIN boards b ON bi.board_id = b.id
      JOIN project_mapping pm ON b.id = pm.row_id AND pm.`table` = 'boards'
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        AND i.type = 'INCIDENT'
        AND i.lead_time_minutes IS NOT NULL
),

_find_median_mttr_each_month_ranks as(
    SELECT *, percent_rank() OVER(PARTITION BY month ORDER BY lead_time_minutes) AS ranks
    FROM _incidents
),

_mttr as(
    SELECT month, max(lead_time_minutes) AS median_time_to_resolve
    FROM _find_median_mttr_each_month_ranks
    WHERE ranks <= 0.5
    GROUP BY month
)

SELECT 
    cm.month AS data_key,
    CASE 
        WHEN m.median_time_to_resolve IS NULL THEN 0 
        ELSE m.median_time_to_resolve/60 
    END AS data_value
FROM 
    calendar_months cm
    LEFT JOIN _mttr m ON cm.month = m.month
  WHERE cm.month_timestamp BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)