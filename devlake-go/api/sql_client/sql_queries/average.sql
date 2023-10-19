SELECT
    t1.data_key,
    AVG(t2.data_value) AS data_value
FROM
    count t1
    JOIN count t2 ON t2.data_key <= t1.data_key
GROUP BY
    t1.data_key;