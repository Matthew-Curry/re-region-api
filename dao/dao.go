package dao

/* Interface for the Re-Region API dao */

import (
	"github.com/Matthew-Curry/re-region-api/apperrors"
)

type DaoInterface interface {
	// state data access methods
	GetStateCensusData() ([][]interface{}, *apperrors.AppError)
	GetStateTax() ([][]interface{}, *apperrors.AppError)
	// county data access method (pull both tax and census information at the same time)
	GetCountyDataById(county_id int) ([][]interface{}, *apperrors.AppError)
	GetCountyDataByName(county_name string) ([][]interface{}, *apperrors.AppError)
	// to pull top listing for a metric for counties
	GetCountyList(metric string, n int, desc bool) ([][]interface{}, *apperrors.AppError)
	// federal tax data access
	GetFederalTaxData() ([][]interface{}, *apperrors.AppError)
}