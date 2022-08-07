SELECT county_id, county_name, state_id, %s
FROM county
WHERE county_name != '32767'
ORDER BY %s %s 
LIMIT ?;