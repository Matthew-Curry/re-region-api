package apperrors

import (
	"fmt"
)

/* Public constructors of app errors for each type of error the app can experience */

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

func NoRows() *AppError{
	message := fmt.Sprintf("No rows returns by query.")
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func StateCensusNotFound(source error) *AppError{
	message := fmt.Sprintf("Unable to retrieve state census data from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func StateTaxNotFound(source error) *AppError{
	message := fmt.Sprintf("Unable to retrieve state tax data from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func CountyListNotFound(source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve county list from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func CountyNameNotFound(county string, source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve data for county %s: %s", county, source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func CountyIDNotFound(county int, source error) *AppError {
	message := fmt.Sprintf("Unable to retrieve data for county id %v: %s", county, source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func FederalTaxNotFound(source error) *AppError{
	message := fmt.Sprintf("Unable to retrieve federal tax data from DB: %s", source.Error())
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func EmptyFederalCache() *AppError{
	return &AppError{message: "The federal tax cache is empty", kind: InternalError, source: nil}
}

func StateIDNotFound(state_id int) *AppError{
	message := fmt.Sprintf("State ID not in state cache %v", state_id)
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func StateNameNotFound(state_name string) *AppError{
	message := fmt.Sprintf("State name not in state cache %s", state_name)
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func StateListNotFound(metric_name string) *AppError{
	message := fmt.Sprintf("Unable to generate state list for metric %s", metric_name)
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func StateNameNotInTaxCache(state_name string) *AppError{
	message := fmt.Sprintf("State name not in state tax cache %s", state_name)
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}

func StateIDNotInTaxCache(state_id int) *AppError{
	message := fmt.Sprintf("State ID not in state tax cache %v", state_id)
	kind := InternalError
	return &AppError{message: message, kind: kind, source: nil}
}
