package services

/* Implementation of the Re-Region API state service */

import (
	"github.com/Matthew-Curry/re-region-api/apperrors"
	"github.com/Matthew-Curry/re-region-api/dao"
	"github.com/Matthew-Curry/re-region-api/model"
)

const (
	// indexes for id, name in the census response
	CENSUS_STATE_ID   int = 0
	CENSUS_STATE_NAME int = 1

	// indexes for the metrics in the census response, used to process the state result
	STATE_POP           int = 2
	STATE_MALE_POP      int = 3
	STATE_FEMALE_POP    int = 4
	STATE_MEDIAN_INCOME int = 5
	STATE_AVERAGE_RENT  int = 6
	STATE_COMMUTE       int = 7
)

// indexes for the state tax information will be very static, so constants can reliably be defined
const (
	TAX_STATE_ID = iota
	TAX_STATE_NAME
	SINGLE_DEDUCTION
	MARRIED_DEDUCTION
	SINGLE_EXEMPTION
	MARRIED_EXEMPTION
	DEPENDENT_EXEMPTION
	SINGLE_RATE
	SINGLE_BRACKET
	MARRIED_RATE
	MARRIED_BRACKET
)

// metrics from the state data response mapped to the index they will be read in to
var metrics = map[string]int{"pop": STATE_POP,
	"male_pop":      STATE_MALE_POP,
	"female_pop":    STATE_FEMALE_POP,
	"median_income": STATE_MEDIAN_INCOME,
	"average_rent":  STATE_AVERAGE_RENT,
	"commute":       STATE_COMMUTE}

type StateServiceImpl struct {
	// maps for the get state endpoint. Map identifiers to base state attributes and
	// calculate tax estimates by request
	stateIdMp   map[int][]interface{}
	stateNameMp map[string][]interface{}

	// map of metrics to ranked lists of state according to the metric
	metricListMp map[string]*model.StateList

	// maps for tax info endpoint
	stateTaxNameMp map[string]*model.StateTaxInfo
	stateTaxIdMp   map[int]*model.StateTaxInfo

	// use provided impl of federal service to access federal tax information
	federalService FederalServiceInterface
}

// constructor to return this implementation of the state service
func GetStateServiceImpl(daoImpl dao.DaoInterface, federalService FederalServiceInterface) (StateServiceInterface, *apperrors.AppError) {
	// use the dao to retrieve data needed for caches
	// fetch data required for the caches from the database
	logger.Info("Getting state census data from the data access layer")
	stateCensusData, err := daoImpl.GetStateCensusData()

	if err != nil {
		return nil, err
	}

	logger.Info("Getting state tax data from the data access layer")
	stateTaxData, err := daoImpl.GetStateTax()

	if err != nil {
		return nil, err
	}

	// build caches
	logger.Info("Building caches")
	stateIdMp, stateNameMp := buildStateCaches(stateCensusData)
	logger.Info("State cache created")

	metricListMp := buildStateListCaches(stateCensusData)
	logger.Info("State list cache created")

	stateTaxIdMp, stateTaxNameMp := buildStateTaxCaches(stateTaxData)
	logger.Info("State tax cache created")
	logger.Info("All state caches are now created")

	// return the constructed service
	return &StateServiceImpl{stateNameMp: stateNameMp,
		stateIdMp:      stateIdMp,
		metricListMp:   metricListMp,
		stateTaxNameMp: stateTaxNameMp,
		stateTaxIdMp:   stateTaxIdMp,
		federalService: federalService}, nil
}

// constructor helper method, builds state caches
func buildStateCaches(stateCensusData [][]interface{}) (map[int][]interface{}, map[string][]interface{}) {
	idMp := make(map[int][]interface{})
	nameMp := make(map[string][]interface{})
	for _, state := range stateCensusData {
		// initialize for each map with opposing identifier so variadic append can be used for remaining metrics
		si := []interface{}{state[CENSUS_STATE_NAME]}
		sn := []interface{}{state[CENSUS_STATE_ID]}
		// use variadic append to adaptively append metrics
		idMp[readAsInt(state[CENSUS_STATE_ID])] = append(si, []interface{}{state[2:]}...)
		nameMp[readAsString(state[CENSUS_STATE_NAME])] = append(sn, []interface{}{state[2:]}...)
	}

	return idMp, nameMp

}

// constructor helper method, builds listing caches
func buildStateListCaches(stateCensusData [][]interface{}) map[string]*model.StateList {
	// initialize map with metric keys defined in this module
	mp := make(map[string]*model.StateList)
	for m := range metrics {
		mp[m] = model.GetMetricStateList(m)
	}
	// pass over states and insert rows to the appropriate cache in the correct order
	for m, i := range metrics {
		for _, state := range stateCensusData {
			// initialize the metric pair for this row
			metricPair := model.StateMetricPair{State_id: readAsInt(state[CENSUS_STATE_ID]), 
												State_name: readAsString(state[CENSUS_STATE_NAME]), 
												Metric_value: readAsInt(state[i])}
			// insert the metric pair into the appropriate slice in order
			mp[m].AppendToRankedList(metricPair)
		}
	}

	return mp

}

// constructor helper method, builds the static tax info caches
func buildStateTaxCaches(stateTaxData [][]interface{}) (map[int]*model.StateTaxInfo, map[string]*model.StateTaxInfo) {
	idMp := make(map[int]*model.StateTaxInfo)
	nameMp := make(map[string]*model.StateTaxInfo)
	for _, row := range stateTaxData {
		// if state tax info is not in id map, create for both maps
		si := readAsInt(row[TAX_STATE_ID])
		sn := readAsString(row[TAX_STATE_NAME])
		if _, ok := idMp[si]; !ok {
			stateTaxInfo := model.GetStateTaxInfo(si, sn, readAsInt(row[SINGLE_DEDUCTION]), readAsInt(row[MARRIED_DEDUCTION]),
			readAsInt(row[SINGLE_EXEMPTION]), readAsInt(row[MARRIED_EXEMPTION]), readAsInt(row[DEPENDENT_EXEMPTION]))

			idMp[si] = stateTaxInfo
			nameMp[sn] = stateTaxInfo

		}
		// append bracket information to the tax info at this id and name position in the respective maps
		sb := model.StateBracket{
			Single_rate:     readAsFloat(row[SINGLE_RATE]),
			Single_bracket:  readAsInt(row[SINGLE_BRACKET]),
			Married_rate:    readAsFloat(row[MARRIED_RATE]),
			Married_bracket: readAsInt(row[MARRIED_BRACKET]),
		}
		// add bracket to list in order of the single rate so they ascend properly. Call for just one map 
		// because they each point to the same tax info
		idMp[si].AppendToOrderedList(sb)

	}

	return idMp, nameMp
}

// get census and tax information by ID
func (s *StateServiceImpl) GetStateById(id int, fs model.FilingStatus, dependents int, income int) (*model.State, *apperrors.AppError) {
	// retrieve state census information using the given id
	sc, ok := s.stateIdMp[id]
	if !ok {
		logger.Warn("State id %s not in the cache", id)
		return nil, apperrors.StateIDNotFound(id)
	}
	// process the yearly tax estimate given this income
	logger.Info("Processing the tax liability for %s", id)
	t, st, ft := s.processTaxLiabilityById(id, fs, dependents, income)

	return s.buildState(sc, t, st, ft), nil
}

// helper method to construct state for given args
func (s *StateServiceImpl) buildState(sc []interface{}, t, st, ft int) *model.State {
	return &model.State{
		State_id:   readAsInt(sc[CENSUS_STATE_ID]),
		State_name: readAsString(sc[CENSUS_STATE_NAME]),
		// state level census metrics
		Pop:           readAsInt(sc[STATE_POP]),
		Male_pop:      readAsInt(sc[STATE_MALE_POP]),
		Female_pop:    readAsInt(sc[STATE_FEMALE_POP]),
		Median_income: readAsInt(sc[STATE_MEDIAN_INCOME]),
		Average_rent:  readAsInt(sc[STATE_AVERAGE_RENT]),
		Commute:       readAsInt(sc[STATE_COMMUTE]),
		// the tax metrics
		Total_tax:   t,
		State_tax:   st,
		Federal_tax: ft,
	}
}

// process state tax liability for a given id
func (s *StateServiceImpl) processTaxLiabilityById(id int, fs model.FilingStatus, dependents int, income int) (int, int, int) {
	ti := s.stateTaxIdMp[id]
	return s.processTaxLiability(fs, dependents, income, ti)
}

// process state tax liability for a given name
func (s *StateServiceImpl) processTaxLiabilityByName(name string, fs model.FilingStatus, dependents int, income int) (int, int, int) {
	ti := s.stateTaxNameMp[name]
	return s.processTaxLiability(fs, dependents, income, ti)
}

// core logic to process state tax liability
func (s *StateServiceImpl) processTaxLiability(fs model.FilingStatus, dependents int, income int, ti *model.StateTaxInfo) (int, int, int) {
	// use filing status to determine state deduction and exemption
	logger.Info("Processing state liability")
	stateTax := 0
	switch fs {
	case model.Head, model.Single:
		income = income - ti.Single_deduction - dependents*ti.Single_exemption
		stateTax = ti.GetSingleTaxLiability(income)
	case model.Married:
		income = income - ti.Married_deduction - dependents*ti.Married_exemption
		stateTax = ti.GetMarriedTaxLiability(income)
	}
	// apply the rate to the income for the state tax amount
	logger.Info("Processing federal liability")
	federalTax := s.federalService.GetFederalLiability(fs, dependents, income)

	return stateTax + federalTax, stateTax, federalTax
}

// get census and tax information by name
func (s *StateServiceImpl) GetStateByName(name string, fs model.FilingStatus, dependents int, income int) (*model.State, *apperrors.AppError) {
	// retrieve state census information using the given name
	sc, ok := s.stateNameMp[name]
	if !ok {
		logger.Warn("State %s not in the cache", name)
		return nil, apperrors.StateNameNotFound(name)
	}
	// process the yearly tax estimate given this income
	logger.Info("Processing the tax liability for %s", name)
	t, st, ft := s.processTaxLiabilityByName(name, fs, dependents, income)

	return s.buildState(sc, t, st, ft), nil
}

// get state for given metric and size
func (s *StateServiceImpl) GetStateList(metricName string, n int) (*model.StateList, *apperrors.AppError) {
	res, ok := s.metricListMp[metricName]
	if !ok {
		logger.Warn("Metric %s not found in the state list cache", metricName)
		return nil, apperrors.StateListNotFound(metricName)
	}

	return res, nil

}

// get state tax info by id
func (s *StateServiceImpl) GetStateTaxInfoById(id int) (*model.StateTaxInfo, *apperrors.AppError) {
	res, ok := s.stateTaxIdMp[id]
	if !ok {
		logger.Warn("State %s not found in the state tax cache", id)
		return nil, apperrors.StateIDNotInTaxCache(id)
	}
	return res, nil
}

// get state tax info by name
func (s *StateServiceImpl) GetStateTaxInfoByName(name string) (*model.StateTaxInfo, *apperrors.AppError) {

	res, ok := s.stateTaxNameMp[name]
	if !ok {
		logger.Warn("State %s not found in the state tax cache", name)
		return nil, apperrors.StateNameNotInTaxCache(name)
	}

	return res, nil

}

// get the state name associated with an id
func (s *StateServiceImpl) getStateNameById(id int) (string, *apperrors.AppError) {

	res, ok := s.stateTaxIdMp[id]
	if !ok {
		logger.Warn("State %s not found in the state tax cache", id)
		return "", apperrors.StateIDNotFound(id)
	}

	return res.State_name, nil
}
