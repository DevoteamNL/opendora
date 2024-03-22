
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
),  _incidents AS (
    SELECT
      DISTINCT i.id,
        i.created_date AS quarter,
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

_find_median_mttr_each_quarter_ranks as(
    SELECT *, percent_rank() OVER(PARTITION BY quarter ORDER BY lead_time_minutes) AS ranks
    FROM _incidents
),

_mttr as(
    SELECT quarter, max(lead_time_minutes) AS median_time_to_resolve
    FROM _find_median_mttr_each_quarter_ranks
    WHERE ranks <= 0.5
    GROUP BY quarter
)

SELECT 
    cq.quarter_date AS data_key,
    CASE 
        WHEN m.median_time_to_resolve IS NULL THEN 0 
        ELSE m.median_time_to_resolve/60 
    END AS data_value
FROM 
        calendar_quarters cq
    LEFT JOIN _mttr m ON cq.quarter_date = m.quarter
    ORDER BY
        cq.quarter_date DESC