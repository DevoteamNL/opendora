WITH _deployments AS(
    SELECT
        date_format(deployment_finished_date, '%y/%m') AS MONTH,
        count(cicd_deployment_id) AS deployment_count
    FROM
        (
            SELECT
                cdc.cicd_deployment_id,
                max(cdc.finished_date) AS deployment_finished_date
            FROM
                cicd_deployment_commits cdc
                JOIN project_mapping pm ON cdc.cicd_scope_id = pm.row_id
                AND pm.`table` = 'cicd_scopes'
            WHERE
                (
                    :project = ""
                    OR pm.project_name = :project
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
    cm.month AS data_key,
    CASE
        WHEN d.deployment_count IS NULL THEN 0
        ELSE d.deployment_count
    END AS data_value
FROM
    calendar_months cm
    LEFT JOIN _deployments d ON cm.month = d.month
WHERE
    cm.month_timestamp BETWEEN FROM_UNIXTIME(:from)
    AND FROM_UNIXTIME(:to)
