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
),
_deployments AS(
    SELECT
        YEARWEEK(deployment_finished_date) AS week,
        count(cicd_deployment_id) AS deployment_count
    FROM
        (
            SELECT
                cdc.cicd_deployment_id,
                max(cdc.finished_date) AS deployment_finished_date
            FROM
                cicd_deployment_commits cdc
                JOIN repos ON cdc.repo_id = repos.id
            WHERE
                (
                    :project = ""
                    OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
                )
                AND cdc.result = 'SUCCESS'
                AND cdc.environment = 'PRODUCTION'
            GROUP BY
                1
        ) _production_deployments
    GROUP BY
        1
),
count AS (
    SELECT
        YEARWEEK(cw.week_date) AS data_key,
        CASE
            WHEN d.deployment_count IS NULL THEN 0
            ELSE d.deployment_count
        END AS data_value
    FROM
        calendar_weeks cw
        LEFT JOIN _deployments d ON YEARWEEK(cw.week_date) = d.week
    ORDER BY
        cw.week_date DESC
)