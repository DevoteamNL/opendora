WITH _incidents AS (
    SELECT
      DISTINCT i.id,
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
        AND i.created_date BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
),

_median_mttr_ranks as(
    SELECT *, percent_rank() OVER(ORDER BY lead_time_minutes) AS ranks
    FROM _incidents
),

_median_mttr as(
    SELECT max(lead_time_minutes) as median_time_to_resolve
    FROM _median_mttr_ranks
    WHERE ranks <= 0.5
)

SELECT 
    CASE
        WHEN median_time_to_resolve < 60 THEN "elite"
        WHEN median_time_to_resolve < 24 * 60 THEN "high"
        WHEN median_time_to_resolve < 7 * 24 * 60 THEN "medium"
        WHEN median_time_to_resolve >= 7 * 24 * 60 THEN "low"
        ELSE "na-incidents"
        END AS data_key,
    CASE
        WHEN median_time_to_resolve > 0 THEN round(median_time_to_resolve/60,1)
        ELSE ""
        END AS data_value
    FROM 
        _median_mttr
