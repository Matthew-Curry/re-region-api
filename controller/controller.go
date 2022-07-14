package controller

/* Holds functions to handle API requests */

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/Matthew-Curry/re-region-api/model"
	"github.com/Matthew-Curry/re-region-api/apperrors"
)

// helper method to write the response
func writeResponse(w http.ResponseWriter, isGet bool, statusCode int, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// only write the response for get requests
	if isGet {
		w.Write(b)
	}
}

// helper method which returns whether an HTTP request method is GET or HEAD (supported methods)
func isGet(r *http.Request) (bool, string) {
	if r.Method == http.MethodGet{
		return true, ""
	} else if r.Method == http.MethodHead {
		return false, ""
	} 

	// invalid method
	return false, fmt.Sprintf("The provided HTTP method %s is unsupported", r.Method)
}

// helper method that will return the name if not empty, else will return the given int as a string
func nameOrId(name string, id int) string {
	if name != ""{
		return name
	}

	return fmt.Sprint(id)
}

// returns no entity available message for entity and identifier
func noEntityAvailable(entity, identifier string) string {
	return fmt.Sprintf("There is no %s %s available.", entity, identifier)
}

// returns unable to retrieve message for entity and its identifier
func unableToRetrieveEntity(entity, identifier string) string {
	return fmt.Sprintf("Unable to retrieve %s %s", entity, identifier)
}

// handler for requests for county resource
func CountyHandler(w http.ResponseWriter, r *http.Request) {
	// params
	id, name, fs, res, dep, income, errStr := getCountyParams(r)
	if errStr != "" {
		writeResponse(w, true, http.StatusBadRequest, []byte(errStr))
	}
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusNotImplemented, []byte(errStr))
	}

	// call the appropriate service method based on the provided params
	var county *model.County
	var err *apperrors.AppError
	if name != "" {
		county, err = countyService.GetCountyByName(name, fs, res, dep, income)
	} else {
		county, err = countyService.GetCountyById(id, fs, res, dep, income)
	}

	// check errors, write the response based on county value
	if err.IsKind(apperrors.DataNotFound) {
		writeResponse(w, isGet, http.StatusNotFound, []byte(noEntityAvailable("county", nameOrId(name, id))))
	} else if err.IsKind(apperrors.InternalError) || err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("county", nameOrId(name, id))))
	} else {
		b, err := county.MarshallCounty()
		if err != nil {
			writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("county", nameOrId(name, id))))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}

	}
}

// handle get requests for the state resource
func StateHandler(w http.ResponseWriter, r *http.Request) {
	// params
	id, name, fs, _, dep, income, errStr := getStateParams(r)
	if errStr != "" {
		writeResponse(w, true, http.StatusBadRequest, []byte(errStr))
	}
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusNotImplemented, []byte(errStr))
	}

	// call the appropriate service method
	var state *model.State
	var err *apperrors.AppError
	if name != "" {
		state, err = stateService.GetStateByName(name, fs, dep, income)
	} else {
		state, err = stateService.GetStateById(id, fs, dep, income)
	}

	// check errors, write the response based on state value
	if err.IsKind(apperrors.DataNotFound) {
		writeResponse(w, isGet, http.StatusNotFound, []byte(noEntityAvailable("state", nameOrId(name, id))))
	} else if err.IsKind(apperrors.InternalError) || err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("state", nameOrId(name, id))))
	} else {
		b, err := state.MarshallState()
		if err != nil {
			writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("state", nameOrId(name, id))))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}

	}
}

// handle get requests for ranked list of counties by a given metric
func CountyListHandler(w http.ResponseWriter, r *http.Request) {
	// params
	metricName, size, errStr := getListParams(r)
	if errStr != "" {
		writeResponse(w, true, http.StatusBadRequest, []byte(errStr))
	}
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusNotImplemented, []byte(errStr))
	}

	// get the county list
	countyList, err := countyService.GetCountyList(metricName, size)

	// write the response based on county value
	if err.IsKind(apperrors.DataNotFound) {
		writeResponse(w, isGet, http.StatusNotFound, []byte(noEntityAvailable("metric", metricName)))
	} else if err.IsKind(apperrors.InternalError) || err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("metric", metricName)))
	} else {
		b, err := countyList.MarshallCountyList()
		if err != nil {
			writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("metric", metricName)))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}

	}

}

// handle get requests for list of states ordered by metric
func StateListHandler(w http.ResponseWriter, r *http.Request) {
	// params
	metricName, size, errStr := getListParams(r)
	if errStr != "" {
		writeResponse(w, true, http.StatusBadRequest, []byte(errStr))
	}
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusNotImplemented, []byte(errStr))
	}

	// if there is an error, return bad response, else call the appropriate service method
	stateList, err := stateService.GetStateList(metricName, size)

	// write the response based on county value
	if err.IsKind(apperrors.DataNotFound) {
		writeResponse(w, isGet, http.StatusNotFound, []byte(noEntityAvailable("metric", metricName)))
	} else if err.IsKind(apperrors.InternalError) || err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("metric", metricName)))
	} else {
		b, err := stateList.MarshallStateList()
		if err != nil {
			writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("metric", metricName)))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}

	}
}

// handle get requests for county tax information
func CountyTaxesHandler(w http.ResponseWriter, r *http.Request) {
	// params
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusNotImplemented, []byte(errStr))
	}

	var countyTaxList *model.CountyTaxList
	var err *apperrors.AppError
	var id int
	if idStr != "" {
		id, convErr := strconv.Atoi(idStr)
		if convErr == nil {
			countyTaxList, err = countyService.GetCountyTaxListById(id)
		} else {
			writeResponse(w, isGet, http.StatusBadRequest, []byte(noEntityAvailable("county", fmt.Sprint(id))))
		}
	} else if name != "" {
		countyTaxList, err = countyService.GetCountyTaxListByName(name)
	}


	if err.IsKind(apperrors.DataNotFound) {
		writeResponse(w, isGet, http.StatusNotFound, []byte(noEntityAvailable("county", nameOrId(name, id))))
	} else if err.IsKind(apperrors.InternalError) || err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("county", nameOrId(name, id))))
	} else {
		b, err := countyTaxList.MarshallCountyTaxList()
		if err != nil {
			writeResponse(w, isGet, http.StatusBadRequest, []byte(unableToRetrieveEntity("county", nameOrId(name, id))))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}
	}
}

// handle get requests for state tax information
func StateTaxesHandler(w http.ResponseWriter, r *http.Request) {
	// params
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusBadRequest, []byte(errStr))
	}

	var stateTaxInfo *model.StateTaxInfo
	var err *apperrors.AppError
	var id int
	if idStr != "" {
		id, convErr := strconv.Atoi(idStr)
		if convErr == nil {
			stateTaxInfo, err = stateService.GetStateTaxInfoById(id)
		} else {
			writeResponse(w, isGet, http.StatusBadRequest, []byte(noEntityAvailable("state", fmt.Sprint(id))))
		}
	} else if name != "" {
		stateTaxInfo, err = stateService.GetStateTaxInfoByName(name)
	}

	if err.IsKind(apperrors.DataNotFound) {
		writeResponse(w, isGet, http.StatusNotFound, []byte(noEntityAvailable("state", nameOrId(name, id))))
	} else if err.IsKind(apperrors.InternalError) || err != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte(unableToRetrieveEntity("state", nameOrId(name, id))))
	} else {
		b, err := stateTaxInfo.MarshallStateTaxInfo()
		if err != nil {
			writeResponse(w, isGet, http.StatusBadRequest, []byte(unableToRetrieveEntity("state", nameOrId(name, id))))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}
	}
}

// handle get requests for federal tax information
func FederalTaxesHandler(w http.ResponseWriter, r *http.Request) {
	// is the request GET or HEAD?
	isGet, errStr := isGet(r)
	if errStr != "" {
		writeResponse(w, false, http.StatusBadRequest, []byte(errStr))
	}

	federlTaxInfo, rErr := federalService.GetFederalTaxInfo()
	if rErr.IsKind(apperrors.InternalError) || rErr != nil {
		writeResponse(w, isGet, http.StatusInternalServerError, []byte("Unable to retrieve federal tax information due to an internal error."))
	}

	b, err := federlTaxInfo.MarshallFederalTaxInfo()
		if err != nil {
			writeResponse(w, isGet, http.StatusInternalServerError, []byte("Unable to retrieve federal tax information due to an internal error."))
		} else {
			writeResponse(w, isGet, http.StatusOK, b)
		}
}

// health endpoint of the app
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, true, http.StatusOK, []byte("API is healthy"))
}
