package controller

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/Matthew-Curry/re-region-api/model"
)

/* Holds validator functions used by the controllers */

func getCountyParams(r *http.Request) (int, string, model.FilingStatus, bool, int, int, string) {
	
	return getGeoParams("county", r)
}

func getStateParams(r *http.Request) (int, string, model.FilingStatus, bool, int, int, string) {

	return getGeoParams("state", r)
}

func getListParams(r *http.Request) (string, int, bool, string) {
	metric := r.URL.Query().Get("metric_name")
	sizeStr := r.URL.Query().Get("size")
	descStr := r.URL.Query().Get("desc")

	if metric == "" {
		return "", 0, false, "A metric must be provided to generate the list."
	}

	desc, err := strconv.ParseBool(descStr)
	if err != nil {
		return "", 0, false, "A boolean like value must be given for whether to make the list descending."
	}

	size, err := strconv.Atoi(sizeStr) 
	if err != nil {
		return "", 0, false, "The size of the list must be an integer."
	}

	return metric, size, desc, ""

}

func getGeoParams(geo string, r *http.Request) (int, string, model.FilingStatus, bool, int, int, string) {
	// concat issues with parametes as encountered for the response
	errorStr := ""

	// read in the expected parameters as strings
	idStr := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	fsStr := r.URL.Query().Get("fs")
	resStr := r.URL.Query().Get("res")
	depStr := r.URL.Query().Get("dep")
	incomeStr := r.URL.Query().Get("income")

	// non string vars
	var id int
	var fs model.FilingStatus
	var res bool
	var dep int
	var income int

	// error
	var err error

	// process each string, check for validations and populate non string vars

	// name or id must be provided
	if name == "" && idStr == "" {
		errorStr = errorStr + fmt.Sprintf("A %s name or id must be provided.", geo)
	} else if name == "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			errorStr = errorStr + fmt.Sprintf("\nThe provided %s id must be an integer.", geo)
		}
	}

	// fs must be valid
	fs, err = model.ToFilingStatus(fsStr)
	if err != nil {
		errorStr = errorStr + "\nThe provided filing status must indicate 'S', 'H', or 'M'."
	}

	// res must be interpretable as bool
	res, err = strconv.ParseBool(resStr)
	if err != nil {
		errorStr = errorStr + "\nThe provided resident flag must be interpretable as a boolean"
	}

	// the dep must be an integer
	dep, err = strconv.Atoi(depStr)
	if err != nil {
		errorStr = errorStr + "\nThe provided number of dependents must be an integer."
	}

	// the income must be an integer
	income, err = strconv.Atoi(incomeStr)
	if err != nil {
		errorStr = errorStr + "\nThe provided income must be an integer."
	}

	return id, name, fs, res, dep, income, errorStr
}
