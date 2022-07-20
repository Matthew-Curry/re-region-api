SELECT 
    county.county_id,
    county.county_name,
    county.state_id,
    county.pop,
    county.male_pop,
    county.female_pop,
    county.median_income,
    county.average_rent,
    county.commute,
    tax_locale.tax_locale_id,
    tax_locale.tax_locale,
    tax_locale.resident_desc,
    tax_locale.resident_rate,
    tax_locale.resident_month_fee,
    tax_locale.resident_year_fee,
    tax_locale.resident_pay_period_fee,
    tax_locale.resident_state_rate,
    tax_locale.nonresident_desc,
    tax_locale.nonresident_rate,
    tax_locale.nonresident_month_fee,
    tax_locale.nonresident_year_fee,
    tax_locale.nonresident_pay_period_fee,
    tax_locale.nonresident_state_rate
FROM county INNER JOIN tax_locale ON county.county_id = tax_locale.county_id
WHERE LOWER(TRIM(county.county_name)) = ? AND county.county_id != 32767;