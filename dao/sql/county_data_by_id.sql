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
    COALESCE(tax_locale.tax_locale_id, 0),
    COALESCE(tax_locale.tax_locale, ''),
    COALESCE(tax_locale.resident_desc, ''),
    COALESCE(tax_locale.resident_rate, 0),
    COALESCE(tax_locale.resident_month_fee, 0),
    COALESCE(tax_locale.resident_year_fee, 0),
    COALESCE(tax_locale.resident_pay_period_fee, 0),
    COALESCE(tax_locale.resident_state_rate, 0),
    COALESCE(tax_locale.nonresident_desc, ''),
    COALESCE(tax_locale.nonresident_rate, 0),
    COALESCE(tax_locale.nonresident_month_fee, 0),
    COALESCE(tax_locale.nonresident_year_fee, 0),
    COALESCE(tax_locale.nonresident_pay_period_fee, 0),
    COALESCE(tax_locale.nonresident_state_rate, 0)
FROM county LEFT JOIN tax_locale ON county.county_id = tax_locale.county_id
WHERE county.county_id = ? AND county.county_id != 32767;