package services

/* Interface for the Re-Region API federal service */

import (
	"github.com/Matthew-Curry/re-region-api/apperrors"
	"github.com/Matthew-Curry/re-region-api/model"
)

type FederalServiceInterface interface {
	// public method for controller get overall federal tax information
	GetFederalTaxInfo() (*model.FederalTaxInfo, *apperrors.AppError)
	// return estimated federal liability
	GetFederalLiability(filingStatus model.FilingStatus, dependents int, income int) int
}