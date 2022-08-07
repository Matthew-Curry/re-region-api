package dao

/* Postgres implementation of a DaoInterface */

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	"github.com/Matthew-Curry/re-region-api/src/logging"
)

const (
	// identifiers to sql query files
	GET_METRIC_SET      string = "GET_METRIC_SET"
	COUNTY_DATA_BY_ID   string = "COUNTY_DATA_BY_ID"
	COUNTY_DATA_BY_NAME string = "COUNTY_DATA_BY_NAME"
	FEDERAL_TAX_DATA    string = "FEDERAL_TAX_DATA"
	STATE_CENSUS_DATA   string = "STATE_CENSUS_DATA"
	STATE_TAX_DATA      string = "STATE_TAX_DATA"
	COUNTY_LIST_DATA    string = "COUNTY_LIST_DATA"

	// sql queries
	GET_METRIC_SET_QUERY      string = "metric_set.sql"
	COUNTY_DATA_BY_ID_QUERY   string = "county_data_by_id.sql"
	COUNTY_DATA_BY_NAME_QUERY string = "county_data_by_name.sql"
	FEDERAL_TAX_DATA_QUERY    string = "federal_tax_data.sql"
	STATE_CENSUS_DATA_QUERY   string = "state_census_data.sql"
	STATE_TAX_DATA_QUERY      string = "state_tax_data.sql"
	COUNTY_LIST_DATA_QUERY    string = "county_list.sql"
)

var logger, _ = logging.GetLogger("file.log")

type DaoImpl struct {
	// holds map of valid metrics to request, populated on startup to valid requested metrics
	metricSet map[string]int
	// the database connection
	con *sql.DB
	// map of identifiers to SQL queries to load in to pull data
	sqlMap map[string]string
}

// public constructor to return the postgres impl of the dao
func GetPostgresDao(user, password, dbname, host, port string) (DaoInterface, *apperrors.AppError) {
	// map of identifiers to sql files
	sqlMap := map[string]string{
		"GET_METRIC_SET":      GET_METRIC_SET_QUERY,
		"COUNTY_DATA_BY_ID":   COUNTY_DATA_BY_ID_QUERY,
		"COUNTY_DATA_BY_NAME": COUNTY_DATA_BY_NAME_QUERY,
		"FEDERAL_TAX_DATA":    FEDERAL_TAX_DATA_QUERY,
		"STATE_CENSUS_DATA":   STATE_CENSUS_DATA_QUERY,
		"STATE_TAX_DATA":      STATE_TAX_DATA_QUERY,
		"COUNTY_LIST_DATA":    COUNTY_LIST_DATA_QUERY,
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	logger.Info("Opening connection to postgres DB")
	c, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, apperrors.DBConnectionError(err)
	}
	logger.Info("Successfully instantiated DB connection")

	d := &DaoImpl{con: c, sqlMap: sqlMap, metricSet: map[string]int{}}
	logger.Info("Instantiated DAO")

	logger.Info("Reading in the valid metrics")
	ae := d.loadMetricSet()

	if err != nil {
		return nil, ae
	}

	return d, nil

}

// method called by constructor to load in valid list of metrics
// on the instantiation of the dao
func (d *DaoImpl) loadMetricSet() *apperrors.AppError {
	query, err := d.readSQLFileAsString(GET_METRIC_SET)

	if err != nil {
		return err
	}

	logger.Info("Getting valid metrics from database")
	res, err := d.getRowsFromQuery(query)
	if err != nil {
		return apperrors.CannotGetMetrics(err)
	}

	logger.Info("Loading the response into the metric set")
	for _, row := range res {
		d.metricSet[string(row[0].([]uint8))] = 0
	}

	return nil

}

func (d *DaoImpl) GetStateCensusData() ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(STATE_CENSUS_DATA)

	if err != nil {
		return nil, err
	}
	logger.Info("Executing State Census query")
	res, err := d.getRowsFromQuery(query)
	if err != nil {
		return nil, apperrors.UnableToGetStateCensus(err)
	}

	return res, nil
}

func (d *DaoImpl) GetCountyList(metric string, n int, desc bool) ([][]interface{}, *apperrors.AppError) {
	// verify the metric is valid
	_, ok := d.metricSet[metric]
	if !ok {
		return nil, apperrors.InvalidCountyMetric()
	}

	query, err := d.readSQLFileAsString(COUNTY_LIST_DATA)

	var order string
	if desc {
		order = "DESC"
	} else {
		order = ""
	}

	// substitute the metric into the query to be selected before passing to DB.
	query = fmt.Sprintf(query, metric, metric, order)

	if err != nil {
		return nil, err
	}
	logger.Info("Executing County list query")
	res, err := d.getRowsFromQuery(query, n)
	if err != nil {
		return nil, apperrors.UnableToGetCountyList(err)
	}
	fmt.Println("EREGERTHRTHRTHTRHRT")
	fmt.Println(res)

	return res, nil
}

func (d *DaoImpl) GetStateTax() ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(STATE_TAX_DATA)

	if err != nil {
		return nil, err
	}
	logger.Info("Executing State tax query")
	res, err := d.getRowsFromQuery(query)
	fmt.Println("STATE CENSUS DATA")
	fmt.Println(res)
	if err != nil {
		return nil, apperrors.UnableToGetStateTax(err)
	}

	return res, nil
}

func (d *DaoImpl) GetCountyDataByName(county_name string) ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(COUNTY_DATA_BY_NAME)

	if err != nil {
		return nil, err
	}

	// ensure name is trimmed and lowercased
	county_name = strings.TrimSpace(strings.ToLower(county_name))

	logger.Info("Executing County by name query")
	res, err := d.getRowsFromQuery(query, county_name)
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			return nil, apperrors.CountyNameNotFound(county_name)
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			return nil, apperrors.UnableToGetCountyName(county_name, err)
		}
	}

	return res, nil

}

func (d *DaoImpl) GetCountyDataById(county_id int) ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(COUNTY_DATA_BY_ID)

	if err != nil {
		return nil, err
	}
	logger.Info("Executing County by id query")

	res, err := d.getRowsFromQuery(query, county_id)
	fmt.Println(res)
	if err != nil {
		if err.IsKind(apperrors.DataNotFound) {
			return nil, apperrors.CountyIDNotFound(county_id)
		} else if err.IsKind(apperrors.InternalError) || err != nil {
			return nil, apperrors.UnableToGetCountyID(county_id, err)
		}
	}

	return res, nil
}

func (d *DaoImpl) GetFederalTaxData() ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(FEDERAL_TAX_DATA)

	if err != nil {
		return nil, err
	}

	logger.Info("Executing Federal tax query")
	res, err := d.getRowsFromQuery(query)
	if err != nil {
		return nil, apperrors.UnableToGetFederalTax(err)
	}

	return res, nil
}

// helper method to read given sql type in for a given query identifier
func (d *DaoImpl) readSQLFileAsString(queryId string) (string, *apperrors.AppError) {
	logger.Info("Reading in SQL for %s", queryId)
	sqlFile, ok := d.sqlMap[queryId]
	if !ok {
		return "", apperrors.NoSQLFileMappedToId(queryId, nil)
	}

	_, filename, _, ok := runtime.Caller(0)
	sqlDir := path.Dir(filename)
	fullSqlPath := sqlDir + "/sql/" + sqlFile

	logger.Info("Reading in file %s", fullSqlPath)
	b, e := ioutil.ReadFile(fullSqlPath)
	if e != nil {
		return "", apperrors.SQLFileReadError(sqlFile, e)
	}

	return string(b), nil

}

// helper method to get rows from a query result. Optionally pass filter values to apply,
// else an empty string
func (d *DaoImpl) getRowsFromQuery(query string, filterValue ...any) ([][]interface{}, *apperrors.AppError) {
	var rows *sql.Rows
	var err error
	if len(filterValue) == 0 {
		logger.Info("Executing the query with no params")
		rows, err = d.con.Query(query)
	} else {
		logger.Info("Executing the query with params")
		// ? -> $n for postgres
		paramCount := strings.Count(query, "?")
		for n := 1; n <= paramCount; n++ {
			query = strings.Replace(query, "?", "$"+strconv.Itoa(n), 1)
		}

		rows, err = d.con.Query(query, filterValue...)
	}

	if err != nil {
		return nil, apperrors.QueryExecutionError(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		return nil, apperrors.QueryExecutionError(err)
	}

	var result [][]interface{}
	pointers := make([]interface{}, len(cols))
	container := make([]interface{}, len(cols))

	for i := range pointers {
		pointers[i] = &container[i]
	}

	areRows := false
	logger.Info("Scanning the rows")
	for rows.Next() {
		if areRows == false {
			areRows = true
		}

		rows.Scan(pointers...)

		scannedRow := make([]interface{}, len(cols))
		copy(scannedRow, container)
		result = append(result, scannedRow)

	}

	if !areRows {
		return nil, apperrors.NoRows()
	}

	return result, nil
}
