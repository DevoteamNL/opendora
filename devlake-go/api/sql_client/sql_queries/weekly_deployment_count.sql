WITH calendar_weeks AS(
    SELECT
        CAST(
            (
                FROM_UNIXTIME(:to) - INTERVAL (ten_unit + unit) WEEK
            ) AS date
        ) week
    FROM
        (
            SELECT
                0 ten_unit
            UNION
            ALL
            SELECT
                10
            UNION
            ALL
            SELECT
                20
            UNION
            ALL
            SELECT
                30
            UNION
            ALL
            SELECT
                40
            UNION
            ALL
            SELECT
                50
        ) ten_unit
        CROSS JOIN (
            SELECT
                0 unit
            UNION
            ALL
            SELECT
                1
            UNION
            ALL
            SELECT
                2
            UNION
            ALL
            SELECT
                3
            UNION
            ALL
            SELECT
                4
            UNION
            ALL
            SELECT
                5
            UNION
            ALL
            SELECT
                6
            UNION
            ALL
            SELECT
                7
            UNION
            ALL
            SELECT
                8
            UNION
            ALL
            SELECT
                9
        ) unit
    WHERE
        (
            FROM_UNIXTIME(:to) - INTERVAL (ten_unit + unit) WEEK
        ) BETWEEN FROM_UNIXTIME(:from)
        AND FROM_UNIXTIME(:to)
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
                JOIN project_mapping pm ON cdc.cicd_scope_id = pm.row_id
                AND pm.` + "` TABLE `" + ` = 'cicd_scopes'
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
    YEARWEEK(cw.week) AS data_key,
    CASE
        WHEN d.deployment_count IS NULL THEN 0
        ELSE d.deployment_count
    END AS data_value
FROM
    calendar_weeks cw
    LEFT JOIN _deployments d ON YEARWEEK(cw.week) = d.week
WHERE
    cw.week BETWEEN FROM_UNIXTIME(:from)
    AND FROM_UNIXTIME(:to)
ORDER BY
    cw.week DESC