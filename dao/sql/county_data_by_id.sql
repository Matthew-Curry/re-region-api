SELECT 
    county.county_id
    county.county_name
    county.state_id
    county.pop
    county.male_pop
    county.female_pop
    county.median_income
    county.average_rent
    county.commute
    tax_jurisdiction.tax_jurisdiction_id
    tax_jurisdiction.tax_jurisdiction_name
    tax_jurisdiction.resident_desc
    tax_jurisdiction.resident_rate
    tax_jurisdiction.resident_month_fee
    tax_jurisdiction.resident_year_fee
    tax_jurisdiction.resident_pay_period_fee
    tax_jurisdiction.resident_state_rate
    tax_jurisdiction.nonresident_desc
    tax_jurisdiction.nonresident_rate
    tax_jurisdiction.nonresident_month_fee
    tax_jurisdiction.nonresident_year_fee
    tax_jurisdiction.nonresident_pay_period_fee
    tax_jurisdiction.nonresident_state_rate
FROM county INNER JOIN tax_jurisdiction ON county.county_id = tax_jurisdiction.county_id
WHERE county.county_id = ?