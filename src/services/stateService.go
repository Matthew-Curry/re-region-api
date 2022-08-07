package services

/* Interface for the Re-Region API state service */

import (
	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	"github.com/Matthew-Curry/re-region-api/src/model"
)

type StateServiceInterface interface {
	// public methods to request a state. Also takes filing status, dependents, and income to estimate taxes
	GetStateById(id int, fs model.FilingStatus, dependents int, income int) (*model.State, *apperrors.AppError)
	GetStateByName(name string, fs model.FilingStatus, dependents int, income int) (*model.State, *apperrors.AppError)
	// public methods to request state list by metric name, list size, and whether the list is ascending or descending
	GetStateList(metricName string, n int, desc bool) (*model.StateList, *apperrors.AppError)
	// public methods to request the tax info for a state
	GetStateTaxInfoById(id int) (*model.StateTaxInfo, *apperrors.AppError)
	GetStateTaxInfoByName(name string) (*model.StateTaxInfo, *apperrors.AppError)
	// internal methods to the package
	// lookup of state id to name
	getStateNameById(id int) (string, *apperrors.AppError)
	// process state tax liability given the id
	processTaxLiabilityById(id int, filingStatus model.FilingStatus, dependents int, income int) (int, int, int)
}
