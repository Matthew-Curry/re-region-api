package controller

/* Holds functions to handle API requests */

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Matthew-Curry/re-region-api/apperrors"
	"github.com/Matthew-Curry/re-region-api/model"
)

// helper method which returns whether an HTTP request method is GET or HEAD (supported methods)
func getHTTPMethod(r *http.Request) (bool, bool, string) {
	if r.Method == http.MethodGet {
		return true, false, ""
	} else if r.Method == http.MethodOptions {
		return false, true, ""
	} else if r.Method == http.MethodHead {
		return false, false, ""
	}

	// invalid method
	return false, false, fmt.Sprintf("The provided HTTP method %s is unsupported", r.Method)
}

// handler for requests for county resource
func CountyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get county called")
	start := time.Now()
	// params
	id, name, fs, res, dep, income, errStr := getCountyParams(r)
	if errStr != "" {
		writeGotBadParams(w, errStr)
		return
	}
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	}
	if errStr != "" {
		writeStatusNotImpl(w, errStr)
		return
	}

	// call the appropriate service method based on the provided params
	var county *model.County
	var err *apperrors.AppError
	if name != "" {
		logger.Info("Getting county", name)
		county, err = countyService.GetCountyByName(name, fs, res, dep, income)
	} else {
		logger.Info("Getting county %v", id)
		county, err = countyService.GetCountyById(id, fs, res, dep, income)
	}

	// check errors, write the response based on county value
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			writeNoEntityAvailable(w, isGet, "county", nameOrId(name, id))
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			writeUnableToGetEntity(w, err, isGet, "county", nameOrId(name, id))
		}
	} else {
		b, err := county.MarshallCounty()
		if err != nil {
			writeGotMarshallError(w, err, isGet, "county", nameOrId(name, id))
		} else {
			write200Response(w, isGet, start, b)
		}

	}
}

// handle get requests for the state resource
func StateHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get state called")
	start := time.Now()
	// params
	id, name, fs, _, dep, income, errStr := getStateParams(r)
	if errStr != "" {
		writeGotBadParams(w, errStr)
		return
	}
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	}
	if errStr != "" {
		writeStatusNotImpl(w, errStr)
		return
	}

	// call the appropriate service method
	var state *model.State
	var err *apperrors.AppError
	if name != "" {
		logger.Info("Getting state", name)
		state, err = stateService.GetStateByName(name, fs, dep, income)
	} else {
		logger.Info("Getting state %v", id)
		state, err = stateService.GetStateById(id, fs, dep, income)
	}

	// check errors, write the response based on state value
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			writeNoEntityAvailable(w, isGet, "state", nameOrId(name, id))
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			writeUnableToGetEntity(w, err, isGet, "state", nameOrId(name, id))
		}
	} else {
		b, err := state.MarshallState()
		if err != nil {
			writeGotMarshallError(w, err, isGet, "state", nameOrId(name, id))
		} else {
			write200Response(w, isGet, start, b)
		}

	}
}

// handle get requests for ranked list of counties by a given metric
func CountyListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get county list called")
	start := time.Now()
	// params
	metricName, size, desc, errStr := getListParams(r)
	if errStr != "" {
		writeGotBadParams(w, errStr)
		return 
	}
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	}
	if errStr != "" {
		writeStatusNotImpl(w, errStr)
		return
	}

	// get the county list
	logger.Info("Getting county list for metric %s", metricName)
	countyList, err := countyService.GetCountyList(metricName, size, desc)

	// write the response based on county value
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			writeNoEntityAvailable(w, isGet, "metric", metricName)
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			writeUnableToGetEntity(w, err, isGet, "metric", metricName)
		}
	} else {
		b, err := countyList.MarshallCountyList()
		if err != nil {
			writeGotMarshallError(w, err, isGet, "metric", metricName)
		} else {
			write200Response(w, isGet, start, b)
		}

	}

}

// handle get requests for list of states ordered by metric
func StateListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get state list called")
	start := time.Now()
	// params
	metricName, size, desc, errStr := getListParams(r)
	if errStr != "" {
		writeGotBadParams(w, errStr)
		return
	}
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	}
	if errStr != "" {
		writeStatusNotImpl(w, errStr)
		return
	}

	// retrieve the state list
	logger.Info("Getting state list for metric %s", metricName)
	stateList, err := stateService.GetStateList(metricName, size, desc)

	// write the response based on state value
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			writeNoEntityAvailable(w, isGet, "metric", metricName)
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			writeUnableToGetEntity(w, err, isGet, "metric", metricName)
		}
	} else {
		b, err := stateList.MarshallStateList()
		if err != nil {
			writeGotMarshallError(w, err, isGet, "metric", metricName)
		} else {
			write200Response(w, isGet, start, b)
		}

	}
}

// handle get requests for county tax information
func CountyTaxesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get county tax info called")
	start := time.Now()
	// params
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	}
	if errStr != "" {
		writeGotBadParams(w, errStr)
		return
	} else if idStr == "" && name == "" {
		writeGotBadParams(w, "A county id or name must be provided to retrieve tax information.")
		return
	}

	var countyTaxList *model.CountyTaxList
	var err *apperrors.AppError
	var id int
	if idStr != "" {
		id, convErr := strconv.Atoi(idStr)
		if convErr == nil {
			logger.Info("Getting tax information for county %v", id)
			countyTaxList, err = countyService.GetCountyTaxListById(id)
		} else {
			writeNoEntityAvailable(w, isGet, "county", fmt.Sprint(id))
		}
	} else if name != "" {
		logger.Info("Getting tax information for county %s", name)
		countyTaxList, err = countyService.GetCountyTaxListByName(name)
	}

	// handle response based on response from county service
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			writeNoEntityAvailable(w, isGet, "county", nameOrId(name, id))
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			writeUnableToGetEntity(w, err, isGet, "county", nameOrId(name, id))
		}
	} else {
		b, err := countyTaxList.MarshallCountyTaxList()
		if err != nil {
			writeGotMarshallError(w, err, isGet, "county", nameOrId(name, id))
		} else {
			write200Response(w, isGet, start, b)
		}
	}
}

// handle get requests for state tax information
func StateTaxesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get state tax info called")
	start := time.Now()
	// params
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	} else if errStr != "" {
		writeGotBadParams(w, errStr)
		return
	} else if idStr == "" && name == "" {
		writeGotBadParams(w, "A state id or name must be provided to retrieve tax information.")
		return
	}

	var stateTaxInfo *model.StateTaxInfo
	var err *apperrors.AppError
	var id int
	if idStr != "" {
		id, convErr := strconv.Atoi(idStr)
		if convErr == nil {
			logger.Info("Getting tax information for state %v", id)
			stateTaxInfo, err = stateService.GetStateTaxInfoById(id)
		} else {
			writeNoEntityAvailable(w, isGet, "state", fmt.Sprint(id))
			return
		}
	} else if name != "" {
		logger.Info("Getting tax information for state %s", name)
		stateTaxInfo, err = stateService.GetStateTaxInfoByName(name)
	}

	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			writeNoEntityAvailable(w, isGet, "state", nameOrId(name, id))
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			writeUnableToGetEntity(w, err, isGet, "state", nameOrId(name, id))
		}
	} else {
		b, err := stateTaxInfo.MarshallStateTaxInfo()
		if err != nil {
			writeGotMarshallError(w, err, isGet, "state", nameOrId(name, id))
		} else {
			write200Response(w, isGet, start, b)
		}
	}
}

// handle get requests for federal tax information
func FederalTaxesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Get federal tax info called")
	start := time.Now()
	// http method validation
	isGet, isOption, errStr := getHTTPMethod(r)
	if isOption {
		writePreFlightRequest(w)
		return
	}
	if errStr != "" {
		writeGotBadParams(w, errStr)
		return
	}

	logger.Info("Getting federal tax info")
	federlTaxInfo, err := federalService.GetFederalTaxInfo()
	if err != nil {
		if err.IsKind(apperrors.InternalError) || err != nil {
			writeResponse(w, isGet, http.StatusInternalServerError, []byte("Unable to retrieve federal tax information due to an internal error."))
		}
	}
	b, err := federlTaxInfo.MarshallFederalTaxInfo()
	if err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte("Unable to retrieve federal tax information due to an internal error."))
	} else {
		write200Response(w, isGet, start, b)
	}

}

// health endpoint of the app
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, true, http.StatusOK, []byte("API is healthy"))
}

// 
