package services

/* Interface for the Re-Region API County service */

import (
	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	"github.com/Matthew-Curry/re-region-api/src/model"
)

type CountyServiceInterface interface {
	// public methods to request a County
	GetCountyById(id int, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *apperrors.AppError)
	GetCountyByName(name string, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *apperrors.AppError)
	// public method to request County list by metric name and size
	GetCountyList(metricName string, n int, desc bool) (*model.CountyList, *apperrors.AppError)
	// public methods to request the tax info for a County
	GetCountyTaxListById(id int) (*model.CountyTaxList, *apperrors.AppError)
	GetCountyTaxListByName(name string) (*model.CountyTaxList, *apperrors.AppError)
}
