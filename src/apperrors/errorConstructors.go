package apperrors

import (
	"fmt"
)

/* Public constructors of app errors for each type of error the app can experience */

// Dao DB + query execution errors
func DBConnectionError(source error) *AppError{
	message := fmt.Sprintf("Cannot connect to the DB.")
	kind := InternalError
	return &AppError{message: message, kind: kind, source: source}
}

func NoSQLFileMappedToId(identifer string, source error) *AppError{
	message := fmt.Sprintf("There is no mapping of the identifer %s to a SQL file.", identifer)
	kind := InternalError
	return &AppError{message: message, kind: kind, source: source}
}

func SQLFileReadError(sqlFileName string, source error) *AppError{
	message := fmt.Sprintf("Cannot read given SQL file %s due to error %s", sqlFileName, source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: source}
}

func QueryExecutionError(source error) *AppError{
	message := fmt.Sprintf("Could not run query due to error:%s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: source}
}

func CannotGetMetrics(source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve list of metrics from the database: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func NoRows() *AppError{
	message := fmt.Sprintf("No rows returns by query.")
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

// DataNotFound Errors for domain
func InvalidCountyMetric() *AppError {
	message := fmt.Sprintf("The provided metric is not in the valid set")
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

func CountyNameNotFound(county string) *AppError {
	message := fmt.Sprintf("County %s is not in the system", county)
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

func CountyIDNotFound(county int) *AppError {
	message := fmt.Sprintf("County id %v is not in the system", county)
	kind := DataNotFound
	a := AppError{message: message, kind: kind, source: nil}
	return &a
}

func StateIDNotFound(state_id int) *AppError{
	message := fmt.Sprintf("State ID not in state cache %v", state_id)
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

func StateNameNotFound(state_name string) *AppError{
	message := fmt.Sprintf("State name not in state cache %s", state_name)
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

func StateIDNotInTaxCache(state_id int) *AppError{
	message := fmt.Sprintf("State ID not in state tax cache %v", state_id)
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

func StateNameNotInTaxCache(state_name string) *AppError{
	message := fmt.Sprintf("State name not in state tax cache %s", state_name)
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

func InvalidStateMetric(metric_name string) *AppError{
	message := fmt.Sprintf("The provided metric %s is not in the state cache", metric_name)
	kind := DataNotFound
	return &AppError{message: message, kind: kind, source: nil}
}

// Internal errors involving accessing domain
func UnableToGetStateCensus(source error) *AppError{
	message := fmt.Sprintf("Unable to retrieve state census data from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func UnableToGetStateTax(source error) *AppError{
	message := fmt.Sprintf("Unable to retrieve state tax data from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func UnableToGetCountyList(source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve county list from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func UnableToGetCountyName(county string, source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve data for county %s: %s", county, source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func UnableToGetCountyID(county int, source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve data for county id %v: %s", county, source.Error())
	kind := InternalError
	a := AppError{message: message, kind: kind, source: nil}
	return &a
}

func UnableToGetFederalTax(source error) *AppError{
	message := fmt.Sprintf("Unable to retrieve federal tax data from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func EmptyFederalCache() *AppError{
	return &AppError{message: "The federal tax cache is empty", kind: InternalError, source: nil}
}

// marhshalling errors
func UnableToMarshall(source error) *AppError{
	message := fmt.Sprintf("Unable to marshall response object: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: source}
}