package model

import (
	"encoding/json"
)

type County struct {
	County_id int
	County_name string
	State_id int
	State_name string
	// county metrics
	Pop int
    Male_pop int
    Female_pop int
    Median_income int
    Average_rent int
    Commute int
	// list of tax jurisdictions with tax information
	Tax_locale []TaxLocale
}

type TaxLocale struct {
	Locale_id int
	Locale_name string
	Total_tax int
	Federal_tax int
	State_tax int
	Locale_tax int
}


// marshaller for controller
func (c *County) MarshallCounty() ([]byte, error) {
	return json.Marshal(c)
}