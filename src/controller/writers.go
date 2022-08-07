package controller

import (
	"net/http"
	"fmt"
	"time"
)

/* Methods used to write different types of responses by handler functions */

// helper method to write the response
func writeResponse(w http.ResponseWriter, isGet bool, statusCode int, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.WriteHeader(statusCode)
	// only write the response for get requests
	if isGet {
		w.Write(b)
	}
}

// helper method that will return the name if not empty, else will return the given int as a string
func nameOrId(name string, id int) string {
	if name != ""{
		return name
	}

	return fmt.Sprint(id)
}

// write the response to the preflight options request
func writePreFlightRequest(w http.ResponseWriter) {
	writeResponse(w, false, http.StatusOK, []byte{})
}

// log and response message for no entity available
func writeNoEntityAvailable(w http.ResponseWriter, isGet bool, entity, identifier string) {
	m := fmt.Sprintf("There is no %s %s available.", entity, identifier)
	logger.Warn(m)
	writeResponse(w, isGet, http.StatusNotFound, []byte(m))
}

// log and response method for unable to retrieve resource
func writeUnableToGetEntity(w http.ResponseWriter, e error, isGet bool, entity, identifier string) {
	m := fmt.Sprintf("Unable to retrieve %s %s", entity, identifier)
	logger.Error(m, e.Error())
	writeResponse(w, isGet, http.StatusNotFound, []byte(m))
}

// log and response method for bad params
func writeGotBadParams(w http.ResponseWriter, errStr string) {
	logger.Warn("Bad params recieved, writing bad request response: %s", errStr)
	writeResponse(w, true, http.StatusBadRequest, []byte(errStr))
}

// log and response method for status not implemented
func writeStatusNotImpl(w http.ResponseWriter, errStr string) {
	logger.Warn("%s, writing status not implemented response", errStr)
	writeResponse(w, false, http.StatusNotImplemented, []byte(errStr))
}

// helper method to log a marshall error and write the response
func writeGotMarshallError(w http.ResponseWriter, e error, isGet bool, entity, identifier string) {
	logger.Error("Unable to marhsall", entity, identifier, e.Error())
	writeResponse(w, isGet, http.StatusInternalServerError, []byte(fmt.Sprintf("Unable to retrieve %s %s", entity, identifier)))
}

// helper method to log and write a 200 response
func write200Response(w http.ResponseWriter, isGet bool, start time.Time, b []byte) {
	writeResponse(w, isGet, http.StatusOK, b)
	elapsed := time.Since(start)
	logger.Info("Returned 200 response in %s", elapsed)
}