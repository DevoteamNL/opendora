WITH last_few_calendar_months AS(
    SELECT
        CAST((SYSDATE() - INTERVAL (H + T + U) DAY) AS date) DAY
    FROM
        (
            SELECT
                0 H
            UNION
            ALL
            SELECT
                100
            UNION
            ALL
            SELECT
                200
            UNION
            ALL
            SELECT
                300
        ) H
        CROSS JOIN (
            SELECT
                0 T
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
            UNION
            ALL
            SELECT
                60
            UNION
            ALL
            SELECT
                70
            UNION
            ALL
            SELECT
                80
            UNION
            ALL
            SELECT
                90
        ) T
        CROSS JOIN (
            SELECT
                0 U
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
        ) U
    WHERE
        (SYSDATE() - INTERVAL (H + T + U) DAY) > FROM_UNIXTIME(:from)
        AND (SYSDATE() - INTERVAL (H + T + U) DAY) < FROM_UNIXTIME(:to)
),
_production_deployment_days AS(
    SELECT
        cdc.cicd_deployment_id AS deployment_id,
        max(DATE(cdc.finished_date)) AS DAY
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
),
_days_weeks_deploy AS(
    SELECT
        date(
            DATE_ADD(
                last_few_calendar_months.day,
                INTERVAL - WEEKDAY(last_few_calendar_months.day) DAY
            )
        ) AS week,
        MAX(
            IF(
                _production_deployment_days.day IS NOT NULL,
                1,
                0
            )
        ) AS weeks_deployed,
        COUNT(DISTINCT _production_deployment_days.day) AS days_deployed
    FROM
        last_few_calendar_months
        LEFT JOIN _production_deployment_days ON _production_deployment_days.day = last_few_calendar_months.day
    GROUP BY
        week
),
_monthly_deploy AS(
    SELECT
        DATE(
            DATE_ADD(
                last_few_calendar_months.day,
                INTERVAL - DAY(last_few_calendar_months.day) + 1 DAY
            )
        ) AS MONTH,
        MAX(
            IF(
                _production_deployment_days.day IS NOT NULL,
                1,
                0
            )
        ) AS months_deployed
    FROM
        last_few_calendar_months
        LEFT JOIN _production_deployment_days ON _production_deployment_days.day = last_few_calendar_months.day
    GROUP BY
        MONTH
),
_median_number_of_deployment_days_per_week_ranks AS(
    SELECT
        *,
        percent_rank() OVER(
            ORDER BY
                days_deployed
        ) AS ranks
    FROM
        _days_weeks_deploy
),
_median_number_of_deployment_days_per_week AS(
    SELECT
        max(days_deployed) AS median_number_of_deployment_days_per_week
    FROM
        _median_number_of_deployment_days_per_week_ranks
    WHERE
        ranks <= 0.5
),
_median_number_of_deployment_days_per_month_ranks AS(
    SELECT
        *,
        percent_rank() OVER(
            ORDER BY
                months_deployed
        ) AS ranks
    FROM
        _monthly_deploy
),
_median_number_of_deployment_days_per_month AS(
    SELECT
        max(months_deployed) AS median_number_of_deployment_days_per_month
    FROM
        _median_number_of_deployment_days_per_month_ranks
    WHERE
        ranks <= 0.5
)
SELECT
    CASE
        WHEN median_number_of_deployment_days_per_week >= 3 THEN 'on-demand'
        WHEN median_number_of_deployment_days_per_week >= 1 THEN 'week-month'
        WHEN median_number_of_deployment_days_per_month >= 1 THEN 'month-6month'
        ELSE 'lt-6month'
    END AS data_key
FROM
    _median_number_of_deployment_days_per_week,
    _median_number_of_deployment_days_per_month