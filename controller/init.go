package controller

import (
	"github.com/Matthew-Curry/re-region-api/dao"
	"github.com/Matthew-Curry/re-region-api/logging"
	"github.com/Matthew-Curry/re-region-api/services"
)

/* Initiliaze logger, core services for the controller */

var logger, _ = logging.GetLogger("file.log")

// services
var daoImpl dao.DaoInterface = nil
var federalService services.FederalServiceInterface = nil
var stateService services.StateServiceInterface = nil
var countyService services.CountyServiceInterface = nil

// public method called to initialize services if they have not been initilized
func InitServices() error {
	// initialize any nil services in order of dependency
	var err error = nil
	if daoImpl == nil {
		daoImpl, err = dao.GetPostgresDao()
		if err != nil {
			return err
		}
	}

	if federalService == nil {
		federalService, err = services.GetFederalServiceImpl(daoImpl)
		if err != nil {
			return err
		}
	}

	if stateService == nil {
		stateService, err = services.GetStateServiceImpl(daoImpl, federalService)
		if err != nil {
			return err
		}
	}

	if countyService == nil {
		countyService, err = services.GetCountyServiceImpl(daoImpl, stateService)
		if err != nil {
			return err
		}
	}

	return nil

}