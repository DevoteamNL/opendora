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
),_incidents AS (
    SELECT
      DISTINCT i.id,
        YEARWEEK(i.created_date) AS week,
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

_find_median_mttr_each_week_ranks as(
    SELECT *, percent_rank() OVER(PARTITION BY week ORDER BY lead_time_minutes) AS ranks
    FROM _incidents
),

_mttr as(
    SELECT week, max(lead_time_minutes) AS median_time_to_resolve
    FROM _find_median_mttr_each_week_ranks
    WHERE ranks <= 0.5
    GROUP BY week
)

SELECT 
        YEARWEEK(cw.week_date) AS data_key,
    CASE 
        WHEN m.median_time_to_resolve IS NULL THEN 0 
        ELSE m.median_time_to_resolve/60 END AS data_value
FROM 
    calendar_weeks cw
    LEFT JOIN _mttr m ON YEARWEEK(cw.week_date) = m.week
    ORDER BY
        cw.week_date DESC