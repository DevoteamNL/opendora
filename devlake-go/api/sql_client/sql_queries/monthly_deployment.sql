WITH _deployments AS(
    SELECT
        DATE_ADD(
            MAKEDATE(YEAR(deployment_finished_date), 1),
            INTERVAL MONTH(deployment_finished_date) -1 MONTH
        ) AS deployment_month,
        COUNT(cicd_deployment_id) AS deployment_count
    FROM
        (
            SELECT
                cdc.cicd_deployment_id,
                MAX(cdc.finished_date) AS deployment_finished_date
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
        cm.month AS data_key,
        CASE
            WHEN d.deployment_count IS NULL THEN 0
            ELSE d.deployment_count
        END AS data_value
    FROM
        calendar_months cm
        LEFT JOIN _deployments d ON cm.month = DATE_FORMAT(d.deployment_month, '%y/%m')
    WHERE
        cm.month_timestamp BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
    ORDER BY
        cm.month DESC
)