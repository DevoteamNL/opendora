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
),
_deployments AS(
    SELECT
        DATE_ADD(
            MAKEDATE(YEAR(deployment_finished_date), 1),
            INTERVAL QUARTER(deployment_finished_date) -1 QUARTER
        ) AS quarter_date,
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
                    OR repos.name LIKE CONCAT('%/',:project)
                )
                AND cdc.result = 'SUCCESS'
                AND cdc.environment = 'PRODUCTION'
            GROUP BY
                1
        ) _production_deployments
    GROUP BY
        1
)
SELECT
    cq.quarter_date AS data_key,
    CASE
        WHEN d.deployment_count IS NULL THEN 0
        ELSE d.deployment_count
    END AS data_value
FROM
    calendar_quarters cq
    LEFT JOIN _deployments d ON cq.quarter_date = d.quarter_date
ORDER BY
    cq.quarter_date DESC