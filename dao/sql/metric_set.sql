-- get list of available county metrics. Select all fields from county other than name and id
SELECT column_name
FROM information_schema.columns
WHERE table_name = 'county'
    AND column_name NOT IN ('county_id', 'county_name', 'state_id');
