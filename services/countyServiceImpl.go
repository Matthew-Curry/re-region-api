package services

import (
	"github.com/Matthew-Curry/re-region-api/dao"
	"github.com/Matthew-Curry/re-region-api/model"
	"github.com/Matthew-Curry/re-region-api/apperrors"
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
    COUNTY_LIST_POP
    COUNTY_LIST_MALE_POP
    COUNTY_LIST_FEMALE_POP
    COUNTY_LIST_MEDIAN_INCOME
    COUNTY_LIST_AVERAGE_RENT
    COUNTY_LIST_COMMUTE
)

// mapping of metric name inputs to indexes in the metric list response
var metricIndexMap = map[string]int{"pop": COUNTY_LIST_POP,
									"male_pop": COUNTY_LIST_MALE_POP,
									"female_pop": COUNTY_LIST_FEMALE_POP,
									"median_income": COUNTY_MEDIAN_INCOME,
									"average_rent": COUNTY_LIST_AVERAGE_RENT,
									"commute": COUNTY_LIST_COMMUTE}

/* Implementation of the Re-Region API county service */

type CountyServiceImpl struct {
	// maps for the get county endpoint. Map identifiers to base state attributes and
	// calculate tax estimates by request. Populates as requests to database are made
	countyIdMp   map[int]*model.County
	countyNameMp map[string]*model.County

	// maps for tax info endpoint. Populated when requests for counties are made to the database
	countyTaxNameMp map[string]*model.CountyTaxList
	countyTaxIdMp   map[int]*model.CountyTaxList

	// map of metric names to the indexes in the county list response
	metricIndexMap map[string]int

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
		metricIndexMap:  metricIndexMap,
		stateService:    stateService,
		daoImpl:         daoImpl}, nil
}

func (c *CountyServiceImpl) GetCountyById(id int, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *apperrors.AppError) {
	// check if id in map, if not get from db
	county, ok := c.countyIdMp[id]
	if ok {
		// populate the tax information
		// TODO: RAISE CUSTOM ERROR HER IF ID IS NOT IN THE INFO MAP
		countyTaxInfo := c.countyTaxIdMp[id]

		return c.appendLocalTaxToCounty(county, countyTaxInfo, fs, resident, dependents, income), nil
	}

	countyData, err := c.daoImpl.GetCountyDataById(id)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the county
	county, _, _ = c.placeCountyDataInMaps(countyData, fs, resident, dependents, income)

	return county, nil
}

// helper method with core logic to update caches and return responses
func (c *CountyServiceImpl) placeCountyDataInMaps(countyData [][]interface{}, fs model.FilingStatus, resident bool, dependents int, income int) (*model.County, *model.CountyTaxList, *apperrors.AppError) {
	// county name and id
	countyName := countyData[0][COUNTY_NAME].(string)
	countyId := countyData[0][COUNTY_ID].(int)

	stateId := countyData[0][COUNTY_STATE_ID].(int)
	stateName, err := c.stateService.getStateNameById(stateId)
	if err != nil {
		return nil, nil, err
	}

	// process the local tax info for each row
	var taxLocaleInfos []model.TaxLocaleInfo
	var taxLocales []model.TaxLocale
	for _, row := range countyData {
		// attributes of the locality read in from the row
		tli := row[COUNTY_TAX_JURISDICTION_ID].(int)
		tln := row[COUNTY_TAX_JURISDICTION_NAME].(string)
		resDesc := row[COUNTY_RESIDENT_DESC].(string)
		resRate := row[COUNTY_RESIDENT_RATE].(float64)
		resMonthFee := row[COUNTY_RESIDENT_MONTH_FEE].(float64)
		resYearFee := row[COUNTY_RESIDENT_YEAR_FEE].(float64)
		resPayPeriod := row[COUNTY_RESIDENT_PAY_PERIOD_FEE].(float64)
		resStateRate := row[COUNTY_RESIDENT_STATE_RATE].(float64)
		nonResDesc := row[COUNTY_NONRESIDENT_DESC].(string)
		nonResRate := row[COUNTY_NONRESIDENT_RATE].(float64)
		nonResMonthFee := row[COUNTY_NONRESIDENT_MONTH_FEE].(float64)
		nonResYearFee := row[COUNTY_NONRESIDENT_YEAR_FEE].(float64)
		nonResPayPeriod := row[COUNTY_NONRESIDENT_PAY_PERIOD_FEE].(float64)
		nonResStateRate := row[COUNTY_NONRESIDENT_STATE_RATE].(float64)

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
			tl, fl, sl, ll = c.getTaxLiability(tli, tln, stateId, fs, dependents, income, resStateRate,
				resMonthFee, resYearFee, resPayPeriod, resStateRate)
		} else {
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

	c.countyTaxIdMp[countyId] = taxList
	c.countyTaxNameMp[countyName] = taxList

	respCounty := c.buildCounty(countyId, countyName, stateId, stateName, countyData[0], taxLocales)
	// cache the county information with an empty tax local, will use tax info + request info to calculate tax attributes when request arrives
	cacheCounty := c.buildCounty(countyId, countyName, stateId, stateName, countyData[0], []model.TaxLocale{})

	c.countyIdMp[countyId] = cacheCounty
	c.countyNameMp[countyName] = cacheCounty

	return respCounty, taxList, nil

}

// helper method to get the local tax liability
func (c *CountyServiceImpl) getTaxLiability(tli int, tln string, stateId int, fs model.FilingStatus, dep int, income int, rate, monthFee, yearFee, payPeriodFee, stateRate float64) (int, int, int, int) {
	tl, sl, fl := c.stateService.processTaxLiabilityById(stateId, fs, dep, income)
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
		Pop:           countyDataRow[COUNTY_STATE_ID].(int),
		Male_pop:      countyDataRow[COUNTY_MALE_POP].(int),
		Female_pop:    countyDataRow[COUNTY_FEMALE_POP].(int),
		Median_income: countyDataRow[COUNTY_MEDIAN_INCOME].(int),
		Average_rent:  countyDataRow[COUNTY_AVERAGE_RENT].(int),
		Commute:       countyDataRow[COUNTY_COMMUTE].(int),
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
	// check if id in map, if not get from db
	county, ok := c.countyNameMp[name]
	if ok {
		// populate the tax information
		// TODO: RAISE CUSTOM ERROR HER IF ID IS NOT IN THE INFO MAP
		countyTaxInfo := c.countyTaxNameMp[name]

		return c.appendLocalTaxToCounty(county, countyTaxInfo, fs, resident, dependents, income), nil
	}

	countyData, err := c.daoImpl.GetCountyDataByName(name)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the county
	county, _, _ = c.placeCountyDataInMaps(countyData, fs, resident, dependents, income)

	return county, nil
}

func (c *CountyServiceImpl) GetCountyList(metricName string, n int) (*model.CountyList, *apperrors.AppError) {
	// request list from dao
	countyListData, err := c.daoImpl.GetCountyList(metricName, n)
	if err != nil {
		return nil, err
	}

	// pass over rows and append to ranked list
	countyList := model.GetMetricCountyList(metricName)
	for _, countyData := range countyListData {
		stateId := countyData[COUNTY_LIST_STATE_ID].(int)
		stateName, err := c.stateService.getStateNameById(stateId)
		if err != nil {
			return nil, err
		}
		metricValue := c.metricIndexMap[metricName]
		cmp := model.CountyMetricPair{
			County_id : countyData[COUNTY_LIST_ID].(int), 
			County_name : countyData[COUNTY_LIST_NAME].(string),
			State_id : stateId,
			State_name : stateName,
			Metric_value : metricValue,
		}

		// append pair in order
		countyList.AppendToRankedList(cmp)
	}

	return countyList, nil
}

func (c *CountyServiceImpl) GetCountyTaxListById(id int) (*model.CountyTaxList, *apperrors.AppError) {
	// check if name in map, if not get from db
	countyTax, ok := c.countyTaxIdMp[id]
	if ok {
		return countyTax, nil
	}

	countyData, err := c.daoImpl.GetCountyDataById(id)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the tax information list
	_, countyTax, _ = c.placeCountyDataInMaps(countyData, "H", false, 0, 0)

	return countyTax, nil
}

func (c *CountyServiceImpl) GetCountyTaxListByName(name string) (*model.CountyTaxList, *apperrors.AppError) {
	// check if name in map, if not get from db
	countyTax, ok := c.countyTaxNameMp[name]
	if ok {
		return countyTax, nil
	}

	countyData, err := c.daoImpl.GetCountyDataByName(name)
	if err != nil {
		return nil, err
	}

	// place the data in the maps and return the tax information list
	_, countyTax, _ = c.placeCountyDataInMaps(countyData, "H", false, 0, 0)

	return countyTax, nil

}
