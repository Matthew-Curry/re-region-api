package model

import (
	"encoding/json"

	"github.com/Matthew-Curry/re-region-api/apperrors"
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
func (c *County) MarshallCounty() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(c)

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}