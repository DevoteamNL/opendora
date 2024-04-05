WITH _deployments as (
    SELECT
        cdc.cicd_deployment_id AS deployment_id,
        max(cdc.finished_date) AS deployment_finished_date
    FROM 
        cicd_deployment_commits cdc
        JOIN project_mapping pm ON cdc.cicd_scope_id = pm.row_id AND pm.`table` = 'cicd_scopes'
    WHERE
        (
            :project = ""
            OR LOWER(repos.name) LIKE CONCAT('%/', LOWER(:project))
        )
        AND cdc.result = 'SUCCESS'
        AND cdc.environment = 'PRODUCTION'
    GROUP BY 1
),

_failure_caused_by_deployments AS (
    SELECT
        d.deployment_id,
        d.deployment_finished_date,
        count(DISTINCT CASE WHEN i.type = 'INCIDENT' THEN d.deployment_id ELSE NULL END) AS has_incident
    FROM
        _deployments d
        LEFT JOIN project_issue_metrics pim ON d.deployment_id = pim.deployment_id
        LEFT JOIN issues i ON pim.id = i.id
    GROUP BY 1,2
),

_change_failure_rate AS (
    SELECT 
        CASE 
            WHEN count(deployment_id) IS NULL THEN NULL
            ELSE sum(has_incident)/count(deployment_id) END AS change_failure_rate
    FROM
        _failure_caused_by_deployments
),

_is_collected_data AS(
    SELECT
        CASE 
        WHEN COUNT(i.id) = 0 AND COUNT(cdc.id) = 0 THEN 'No All'
        WHEN COUNT(i.id) = 0 THEN 'No Incidents' 
        WHEN COUNT(cdc.id) = 0 THEN 'No Deployments'
        END AS is_collected
FROM
    (SELECT 1) AS dummy
LEFT JOIN
    issues i ON i.type = 'INCIDENT'
LEFT JOIN
    cicd_deployment_commits cdc ON 1=1
)


SELECT
  CASE
        WHEN is_collected = "No All"  THEN "na-deployments-incidents"
        WHEN is_collected = "No Incidents"  THEN "na-incidents"
        WHEN is_collected = "No Deployments"  THEN "na-deployments"
        WHEN change_failure_rate <= .05 THEN  "elite"
        WHEN change_failure_rate <= .10 THEN  "high"
        WHEN change_failure_rate <= .15 THEN  "medium"
        WHEN change_failure_rate > .15 THEN  "low"
        ELSE "na-deployments-incidents"
        END AS data_key,
  CASE
        WHEN change_failure_rate <= .05 THEN round(change_failure_rate*100,1)
        WHEN change_failure_rate <= .10 THEN round(change_failure_rate*100,1)
        WHEN change_failure_rate <= .15 THEN round(change_failure_rate*100,1)
        WHEN change_failure_rate > .15 THEN round(change_failure_rate*100,1)
        ELSE ""
        END AS data_value
        
FROM 
    _change_failure_rate, _is_collected_data