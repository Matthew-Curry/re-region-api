package model

import (
	"encoding/json"

	"github.com/Matthew-Curry/re-region-api/src/apperrors"
)

type CountyTaxList struct {
	County_name string
	County_id   int
	State_name  string
	State_id    int
	Tax_locales []TaxLocaleInfo
}

type TaxLocaleInfo struct {
	Locale_id  int
	Local_name string
	// resident fields
	Resident_desc           string
	Resident_rate           float64
	Resident_month_fee      float64
	Resident_year_fee       float64
	Resident_pay_period_fee float64
	Resident_state_rate     float64
	// non-resident fields
	Nonresident_desc           string
	Nonresident_rate           float64
	Nonresident_month_fee      float64
	Nonresident_year_fee       float64
	Nonresident_pay_period_fee float64
	Nonresident_state_rate     float64
}

// marshallers for controller

func (c *CountyTaxList) MarshallCountyTaxList() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(c)

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}

func (c *TaxLocaleInfo) MarshallTaxLocaleInfo() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(c)

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}
