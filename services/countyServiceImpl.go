package services

import (
	"github.com/Matthew-Curry/re-region-api/dao"
	"github.com/Matthew-Curry/re-region-api/model"
	"github.com/Matthew-Curry/re-region-api/apperrors"

	"strings"
)

// for core county response
const (
	COUNTY_ID = iota
	COUNTY_NAME
	COUNTY_STATE_ID
	COUNTY_POP
	COUNTY_MALE_POP
	COUNTY_FEMALE_POP
	COUNTY_MEDIAN_INCOME
	COUNTY_AVERAGE_RENT
	COUNTY_COMMUTE
	COUNTY_TAX_JURISDICTION_ID
	COUNTY_TAX_JURISDICTION_NAME
	COUNTY_RESIDENT_DESC
	COUNTY_RESIDENT_RATE
	COUNTY_RESIDENT_MONTH_FEE
	COUNTY_RESIDENT_YEAR_FEE
	COUNTY_RESIDENT_PAY_PERIOD_FEE
	COUNTY_RESIDENT_STATE_RATE
	COUNTY_NONRESIDENT_DESC
	COUNTY_NONRESIDENT_RATE
	COUNTY_NONRESIDENT_MONTH_FEE
	COUNTY_NONRESIDENT_YEAR_FEE
	COUNTY_NONRESIDENT_PAY_PERIOD_FEE
	COUNTY_NONRESIDENT_STATE_RATE
)

// for the county list response
const (
	COUNTY_LIST_ID = iota
	COUNTY_LIST_NAME
    COUNTY_LIST_STATE_ID
    COUNTY_LIST_METRIC_VALUE
)

/* Implementation of the Re-Region API county service */

type CountyServiceImpl struct {
	// maps for the get county endpoint. Map identifiers to base state attributes and
	// calculate tax estimates by request. Populates as requests to database are made
	countyIdMp   map[int]*model.County
	countyNameMp map[string]*model.County

	// maps for tax info endpoint. Populated when requests for counties are made to the database
	countyTaxNameMp map[string]*model.CountyTaxList
	countyTaxIdMp   map[int]*model.CountyTaxList

	// use provided impl of state service to access state + federal tax information
	stateService StateServiceInterface
	// use provided implementation of dao service to make requests to the database
	daoImpl      dao.DaoInterface
}

// constructor to return this implementation of the county service
func GetCountyServiceImpl(daoImpl dao.DaoInterface, stateService StateServiceInterface) (CountyServiceInterface, *apperrors.AppError) {
	// initialize implementation with empty caches. Caches will be populated as records are requested
	return &CountyServiceImpl{countyIdMp: map[int]*model.County{},
		countyNameMp:    map[string]*model.County{},
		countyTaxNameMp: map[string]*model.CountyTaxList{},
		countyTaxIdMp:   map[int]*model.CountyTaxList{},
		stateService:    stateService,
		daoImpl:         daoImpl}, nil
}

func (c *CountyServiceImpl) GetCountyById(id int, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *apperrors.AppError) {
	// check if id in map, if not get from db
	county, ok := c.countyIdMp[id]
	if ok {
		// populate the tax information
		logger.Info("County %v found in cache", id)
		countyTaxInfo := c.countyTaxIdMp[id]

		return c.appendLocalTaxToCounty(county, countyTaxInfo, fs, resident, dependents, income), nil
	}
	logger.Info("County %v not found in cache, querying data access layer", id)
	countyData, err := c.daoImpl.GetCountyDataById(id)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the county
	logger.Info("Recieved response, placing data into the appropriate caches")
	county, _, _ = c.placeCountyDataInMaps(countyData, fs, resident, dependents, income)

	return county, nil
}

// helper method with core logic to update caches and return responses
func (c *CountyServiceImpl) placeCountyDataInMaps(countyData [][]interface{}, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *model.CountyTaxList, *apperrors.AppError) {
	// county name and id
	countyName := readAsString(countyData[0][COUNTY_NAME])
	countyId := readAsInt(countyData[0][COUNTY_ID])

	stateId := readAsInt(countyData[0][COUNTY_STATE_ID])
	stateName, err := c.stateService.getStateNameById(stateId)
	if err != nil {
		return nil, nil, err
	}

	// process the local tax info for each row
	var taxLocaleInfos []model.TaxLocaleInfo
	var taxLocales []model.TaxLocale
	for _, row := range countyData {
		// attributes of the locality read in from the row
		tli := readAsInt(row[COUNTY_TAX_JURISDICTION_ID])
		tln := readAsString(row[COUNTY_TAX_JURISDICTION_NAME])
		resDesc := readAsString(row[COUNTY_RESIDENT_DESC])
		resRate := readAsFloat(row[COUNTY_RESIDENT_RATE])
		resMonthFee := readAsFloat(row[COUNTY_RESIDENT_MONTH_FEE])
		resYearFee := readAsFloat(row[COUNTY_RESIDENT_YEAR_FEE])
		resPayPeriod := readAsFloat(row[COUNTY_RESIDENT_PAY_PERIOD_FEE])
		resStateRate := readAsFloat(row[COUNTY_RESIDENT_STATE_RATE])
		nonResDesc := readAsString(row[COUNTY_NONRESIDENT_DESC])
		nonResRate := readAsFloat(row[COUNTY_NONRESIDENT_RATE])
		nonResMonthFee := readAsFloat(row[COUNTY_NONRESIDENT_MONTH_FEE])
		nonResYearFee := readAsFloat(row[COUNTY_NONRESIDENT_YEAR_FEE])
		nonResPayPeriod := readAsFloat(row[COUNTY_NONRESIDENT_PAY_PERIOD_FEE])
		nonResStateRate := readAsFloat(row[COUNTY_NONRESIDENT_STATE_RATE])

		// the different levels of tax liability to populate
		var tl, fl, sl, ll int
		// append static info for a locality
		taxLocaleInfos = append(taxLocaleInfos, model.TaxLocaleInfo{
			Locale_id:              tli,
			Local_name:             tln,
			Resident_desc:              resDesc,
			Resident_rate:              resRate,
			Resident_month_fee:         resMonthFee,
			Resident_year_fee:          resYearFee,
			Resident_pay_period_fee:    resPayPeriod,
			Resident_state_rate:        resStateRate,
			Nonresident_desc:           nonResDesc,
			Nonresident_rate:           nonResRate,
			Nonresident_month_fee:      nonResMonthFee,
			Nonresident_year_fee:       nonResYearFee,
			Nonresident_pay_period_fee: nonResPayPeriod,
			Nonresident_state_rate:     nonResStateRate,
		})

		// process tax liabilities for the given parameters
		if resident {
			logger.Info("Getting resident tax liability")
			tl, fl, sl, ll = c.getTaxLiability(tli, tln, stateId, fs, dependents, income, resStateRate,
				resMonthFee, resYearFee, resPayPeriod, resStateRate)
		} else {
			logger.Info("Getting non-resident tax liability")
			tl, fl, sl, ll = c.getTaxLiability(tli, tln, stateId, fs, dependents, income, nonResStateRate,
				nonResMonthFee, nonResYearFee, nonResPayPeriod, nonResStateRate)
		}

		// append the formed tax locale
		taxLocales = append(taxLocales, model.TaxLocale{
			Locale_id:   tli,
			Locale_name: tln,
			Total_tax:   tl,
			Federal_tax: fl,
			State_tax:   sl,
			Locale_tax:  ll,
		})
	}

	// append to maps an return the county and tax list
	taxList := &model.CountyTaxList{
		County_name: countyName,
		County_id:   countyId,
		Tax_locales: taxLocaleInfos,
	}

	// lowercase and trim the county name for the maps
	lowerCountyName := strings.TrimSpace(strings.ToLower(countyName))

	c.countyTaxIdMp[countyId] = taxList
	c.countyTaxNameMp[lowerCountyName] = taxList

	respCounty := c.buildCounty(countyId, countyName, stateId, stateName, countyData[0], taxLocales)
	// cache the county information with an empty tax local, will use tax info + request info to calculate tax attributes when request arrives
	cacheCounty := c.buildCounty(countyId, countyName, stateId, stateName, countyData[0], []model.TaxLocale{})

	c.countyIdMp[countyId] = cacheCounty
	c.countyNameMp[lowerCountyName] = cacheCounty

	return respCounty, taxList, nil

}

// helper method to get the local tax liability
func (c *CountyServiceImpl) getTaxLiability(tli int, tln string, stateId int, fs model.FilingStatus, dep int, income int, rate, monthFee, yearFee, payPeriodFee, stateRate float64) (int, int, int, int) {
	tl, sl, fl := c.stateService.processTaxLiabilityById(stateId, fs, dep, income)
	logger.Info("Getting county liability")
	ll := int(float64(income)*rate) + int(12*monthFee) + int(yearFee) + int(payPeriodFee*26) + sl*int(stateRate)
	tl = tl + ll

	return tl, fl, sl, ll
}

// helper method to build a county
func (c *CountyServiceImpl) buildCounty(countyId int, countyName string, stateId int, stateName string, countyDataRow []interface{}, taxLocales []model.TaxLocale) *model.County {
	return &model.County{
		County_id:     countyId,
		County_name:   countyName,
		State_id:      stateId,
		State_name:    stateName,
		Pop:           readAsInt(countyDataRow[COUNTY_STATE_ID]),
		Male_pop:      readAsInt(countyDataRow[COUNTY_MALE_POP]),
		Female_pop:    readAsInt(countyDataRow[COUNTY_FEMALE_POP]),
		Median_income: readAsInt(countyDataRow[COUNTY_MEDIAN_INCOME]),
		Average_rent:  readAsInt(countyDataRow[COUNTY_AVERAGE_RENT]),
		Commute:       readAsInt(countyDataRow[COUNTY_COMMUTE]),
		Tax_locale:    taxLocales,
	}
}

// logic to populate tax locales for a given county, tax information, and inputs to tax calculation
func (c *CountyServiceImpl) appendLocalTaxToCounty(county *model.County, countyTaxInfo *model.CountyTaxList, fs model.FilingStatus, resident bool, dependents int, income int) *model.County {
	for _, taxLocale := range countyTaxInfo.Tax_locales {
		tl, fl, sl, ll := c.getTaxLiability(taxLocale.Locale_id, taxLocale.Local_name, county.State_id, fs, dependents, income, taxLocale.Resident_rate,
			taxLocale.Resident_month_fee, taxLocale.Resident_year_fee, taxLocale.Resident_pay_period_fee, taxLocale.Resident_state_rate)

		county.Tax_locale = append(county.Tax_locale, model.TaxLocale{
			Locale_id:   taxLocale.Locale_id,
			Locale_name: taxLocale.Local_name,
			Total_tax:   tl,
			Federal_tax: fl,
			State_tax:   sl,
			Locale_tax:  ll,
		})
	}
	return county
}

func (c *CountyServiceImpl) GetCountyByName(name string, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *apperrors.AppError) {
	// check if name in map, if not get from db
	name = formatCountyInput(name)
	county, ok := c.countyNameMp[name]
	if ok {
		// populate the tax information
		logger.Info("County %s found in cache", name)
		countyTaxInfo := c.countyTaxNameMp[name]

		return c.appendLocalTaxToCounty(county, countyTaxInfo, fs, resident, dependents, income), nil
	}
	logger.Info("County %s not found in cache, querying data access layer", name)
	countyData, err := c.daoImpl.GetCountyDataByName(name)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the county
	logger.Info("Recieved response, placing data into the appropriate caches")
	county, _, _ = c.placeCountyDataInMaps(countyData, fs, resident, dependents, income)

	return county, nil
}

func (c *CountyServiceImpl) GetCountyList(metricName string, n int, desc bool) (*model.CountyList, *apperrors.AppError) {
	// request list from dao
	logger.Info("Querying data access layer for list of counties ranked by metric %s", metricName)
	countyListData, err := c.daoImpl.GetCountyList(metricName, n, desc)
	if err != nil {
		return nil, err
	}

	// pass over rows and append to ranked list
	logger.Info("Processing the response")
	countyList := model.GetMetricCountyList(metricName)
	for _, countyData := range countyListData {
		stateId := readAsInt(countyData[COUNTY_LIST_STATE_ID])
		stateName, err := c.stateService.getStateNameById(stateId)
		if err != nil {
			return nil, err
		}
		
		cmp := model.CountyMetricPair{
			County_id : readAsInt(countyData[COUNTY_LIST_ID]), 
			County_name : readAsString(countyData[COUNTY_LIST_NAME]),
			State_id : stateId,
			State_name : stateName,
			Metric_value : readAsInt(countyData[COUNTY_LIST_METRIC_VALUE]),
		}

		// append to list, order is enforced by query
		countyList.Ranked_list = append(countyList.Ranked_list, cmp)
	}

	return countyList, nil
}

func (c *CountyServiceImpl) GetCountyTaxListById(id int) (*model.CountyTaxList, *apperrors.AppError) {
	// check if id in map, if not get from db
	countyTax, ok := c.countyTaxIdMp[id]
	if ok {
		logger.Info("Found county %v in the tax cache", id)
		return countyTax, nil
	}
	logger.Info("Did not find county %v in the tax cache. Querying data access layer", id)
	countyData, err := c.daoImpl.GetCountyDataById(id)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the tax information list
	logger.Info("Placing county %v data in the correct maps", id)
	_, countyTax, _ = c.placeCountyDataInMaps(countyData, "H", false, 0, 0)

	return countyTax, nil
}

func (c *CountyServiceImpl) GetCountyTaxListByName(name string) (*model.CountyTaxList, *apperrors.AppError) {
	// check if name in map, if not get from db
	name = formatCountyInput(name)
	countyTax, ok := c.countyTaxNameMp[name]
	if ok {
		logger.Info("Found county %s in the tax cache", name)
		return countyTax, nil
	}
	logger.Info("Did not find county %s in the tax cache. Querying data access layer", name)
	countyData, err := c.daoImpl.GetCountyDataByName(name)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the tax information list
	logger.Info("Placing county %s data in the correct maps", name)
	_, countyTax, _ = c.placeCountyDataInMaps(countyData, "H", false, 0, 0)

	return countyTax, nil

}

// Helper method with logic to format a county name input to the public methods.
// Input should be trimmed and lowercased, and if there is no " county" in the name, 
// apended at the end
func formatCountyInput(county string) string {
	county = strings.TrimSpace(strings.ToLower(county))
	if !strings.Contains(county, " county") {
		county = county + " county"
	}

	return county

}
