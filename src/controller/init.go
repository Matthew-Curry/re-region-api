package controller

import (
	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	"github.com/Matthew-Curry/re-region-api/src/dao"
	"github.com/Matthew-Curry/re-region-api/src/logging"
	"github.com/Matthew-Curry/re-region-api/src/services"
)

/* Initiliaze logger, core services for the controller */

var logger, _ = logging.GetLogger("file.log")

// services
var daoImpl dao.DaoInterface = nil
var federalService services.FederalServiceInterface = nil
var stateService services.StateServiceInterface = nil
var countyService services.CountyServiceInterface = nil

// public method called to initialize services if they have not been initilized
func InitServices(user, password, dbName, dbHost, dbPort string) error {
	// initialize any nil services in order of dependency
	var err *apperrors.AppError = nil
	if daoImpl == nil {
		daoImpl, err = dao.GetPostgresDao(user, password, dbName, dbHost, dbPort)
		if err != nil {
			logger.Error("Could not initialize dao service")
			return err
		}

		logger.Info("Successfully initialized dao service")
	}

	if federalService == nil {
		federalService, err = services.GetFederalServiceImpl(daoImpl)
		if err != nil {
			logger.Error("Could not initialize federal service")
			return err
		}

		logger.Info("Successfully initialized federal service")
	}

	if stateService == nil {
		stateService, err = services.GetStateServiceImpl(daoImpl, federalService)
		if err != nil {
			logger.Error("Could not initialize state service")
			return err
		}

		logger.Info("Successfully initialized state service")
	}

	if countyService == nil {
		countyService, err = services.GetCountyServiceImpl(daoImpl, stateService)
		if err != nil {
			logger.Error("Could not initialize county service")
			return err
		}

		logger.Info("Successfully initialized county service")
	}

	return nil

}
