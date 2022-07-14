package dao

/* Postgres implementation of a DaoInterface */

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"database/sql"

	"github.com/Matthew-Curry/re-region-api/apperrors"
)

const (
	// constants defining name of environment vars storing DB info
	RE_REGION_API_USER     string = "RE_REGION_API_USER"
	RE_REGION_API_PASSWORD string = "RE_REGION_API_PASSWORD"
	RE_REGION_DB           string = "RE_REGION_DB"

	// DB connection string constants
	host string = "localhost"
	port int    = 5432

	// identifiers to sql query files
	COUNTY_DATA_BY_ID   string = "COUNTY_DATA_BY_ID"
	COUNTY_DATA_BY_NAME string = "COUNTY_DATA_BY_NAME"
	FEDERAL_TAX_DATA    string = "FEDERAL_TAX_DATA"
	STATE_CENSUS_DATA   string = "STATE_CENSUS_DATA"
	STATE_TAX_DATA      string = "STATE_TAX_DATA"
	COUNTY_LIST_DATA    string = "COUNTY_LIST_DATA"

	// sql queries
	COUNTY_DATA_BY_ID_QUERY   string = "county_data_by_id.sql"
	COUNTY_DATA_BY_NAME_QUERY string = "county_data_by_name.sql"
	FEDERAL_TAX_DATA_QUERY    string = "federal_tax_data.sql"
	STATE_CENSUS_DATA_QUERY   string = "state_census_data.sql"
	STATE_TAX_DATA_QUERY      string = "state_tax_data.sql"
	COUNTY_LIST_DATA_QUERY    string = "county_list_data.sql"
)

type DaoImpl struct {
	// the database connection
	con *sql.DB
	// map of identifiers to SQL queries to load in to pull data
	sqlMap map[string]string
}

// public constructor to return the postgres impl of the dao
func GetPostgresDao() (DaoInterface, error) {
	// map of identifiers to sql files
	sqlMap := map[string]string{
		"COUNTY_DATA_BY_ID":   COUNTY_DATA_BY_ID_QUERY,
		"COUNTY_DATA_BY_NAME": COUNTY_DATA_BY_NAME_QUERY,
		"FEDERAL_TAX_DATA":    FEDERAL_TAX_DATA_QUERY,
		"STATE_CENSUS_DATA":   STATE_CENSUS_DATA_QUERY,
		"STATE_TAX_DATA":      STATE_TAX_DATA_QUERY,
		"COUNTY_LIST_DATA":    COUNTY_LIST_DATA_QUERY,
	}
	// read in creds + database from environment
	user := os.Getenv(RE_REGION_API_USER)
	password := os.Getenv(RE_REGION_API_PASSWORD)
	dbname := os.Getenv(RE_REGION_DB)
	// use these and constants for the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	d, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, apperrors.DBConnectionError(err)
	}

	return &DaoImpl{con: d, sqlMap: sqlMap}, nil

}

func (d *DaoImpl) GetStateCensusData() ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(STATE_CENSUS_DATA)

	if err != nil {
		return nil, err
	}

	res, err := d.getRowsFromQuery(query, "")
	if err != nil {
		return nil, apperrors.StateCensusNotFound(err)
	}

	return res, nil
}

func (d *DaoImpl) GetCountyList(metric string, n int) ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(COUNTY_LIST_DATA)

	if err != nil {
		return nil, err
	}

	res, err := d.getRowsFromQuery(query, metric, metric, strconv.Itoa(n))
	if err != nil {
		return nil, apperrors.CountyListNotFound(err)
	}

	return res, nil
}

func (d *DaoImpl) GetStateTax() ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(STATE_TAX_DATA)

	if err != nil {
		return nil, err
	}

	res, err := d.getRowsFromQuery(query, "")
	if err != nil {
		return nil, apperrors.StateTaxNotFound(err)
	}

	return res, nil
}

func (d *DaoImpl) GetCountyDataByName(county_name string) ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(COUNTY_DATA_BY_NAME)

	if err != nil {
		return nil, err
	}

	res, err := d.getRowsFromQuery(query, county_name)
	if err != nil {
		return nil, apperrors.CountyNameNotFound(county_name, err)
	}

	return res, nil

}

func (d *DaoImpl) GetCountyDataById(county_id int) ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(COUNTY_DATA_BY_ID)

	if err != nil {
		return nil, err
	}

	res, err := d.getRowsFromQuery(query, strconv.Itoa(county_id))
	if err != nil {
		return nil, apperrors.CountyIDNotFound(county_id, err)
	}

	return res, nil
}

func (d *DaoImpl) GetFederalTaxData() ([][]interface{}, *apperrors.AppError) {
	query, err := d.readSQLFileAsString(FEDERAL_TAX_DATA)

	if err != nil {
		return nil, err
	}

	res, err := d.getRowsFromQuery(query, "")
	if err != nil {
		return nil, apperrors.FederalTaxNotFound(err)
	}

	return res, nil
}

// helper method to read given sql type in for a given query identifier
func (d *DaoImpl) readSQLFileAsString(queryId string) (string, *apperrors.AppError) {
	sqlFile, ok := d.sqlMap[queryId]
	if !ok {
		return "", apperrors.NoSQLFileMappedToId(queryId, nil)
	}

	_, filename, _, ok := runtime.Caller(0)
	sqlDir := path.Dir(filename)

	b, e := ioutil.ReadFile(sqlDir + "/sql/" + sqlFile)
	if e != nil {
		return "", apperrors.SQLFileReadError(sqlFile, e)
	}

	return string(b), nil

}

// helper method to get rows from a query result. Optionally pass filter values to apply,
// else an empty string
func (d *DaoImpl) getRowsFromQuery(query string, filterValue ...string) ([][]interface{}, *apperrors.AppError) {
	var rows *sql.Rows
	var err error
	if len(filterValue) == 0 {
		rows, err = d.con.Query(query)
	} else {
		rows, err = d.con.Query(query, filterValue)
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
	for rows.Next() {
		if areRows == false {
			areRows = true
		}
		rows.Scan(pointers...)
		result = append(result, container)

	}

	if !areRows {
		return nil, apperrors.NoRows()
	}

	return result, nil
}
