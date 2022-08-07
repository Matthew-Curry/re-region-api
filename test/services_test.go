package test 

/* Testing suite for Re-Region services */

import (
	"github.com/Matthew-Curry/re-region-api/src/dao"
	"github.com/Matthew-Curry/re-region-api/src/services"
	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	
	"log"
	"reflect"

	_"github.com/Matthew-Curry/re-region-api/src/model"
	"testing"
)

var daoMock dao.DaoInterface
var federalService services.FederalServiceInterface
var stateService services.StateServiceInterface
var countyService services.CountyServiceInterface


func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func setup() {
	var err *apperrors.AppError = nil
	daoMock = GetDaoMock()

	federalService, err = services.GetFederalServiceImpl(daoMock)
	if err != nil {
		log.Panic("Could not initialize federal service", err)
	}

	stateService, err = services.GetStateServiceImpl(daoMock, federalService)
	if err != nil {
		log.Panic("Could not initialize state service", err)
	}

	countyService, err = services.GetCountyServiceImpl(daoMock, stateService)
	if err != nil {
		log.Panic("Could not initialize county service", err)
	}
}

func TestGetCountyById(t *testing.T) {

	res, err := countyService.GetCountyById(5, "S", true, 4, 45000)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}
	
	if !reflect.DeepEqual(exCounty, res) {
		t.Error("The GetCountyId response does not match the expectation.")
	}
}


func TestGetCountyByName(t *testing.T) {
	res, err := countyService.GetCountyByName("name", "S", true, 4, 45000)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}
	
	if !reflect.DeepEqual(exCounty, res) {
		t.Error("The GetCountyName response does not match the expectation.")
	}
}



func TestGetCountyList(t *testing.T){
	res, err := countyService.GetCountyList("metric", 5, true)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}

	if !reflect.DeepEqual(exCountyList, res) {
		t.Error("The GetCountyList response does not match the expectation.")
	}
}

func TestGetCountyTaxListById(t *testing.T){
	res, err := countyService.GetCountyTaxListById(5)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}

	if !reflect.DeepEqual(exCountyTaxList, res) {
		t.Error("The GetCountyTaxListById response does not match the expectation.")
	}
}

func TestGetCountyTaxListByName(t *testing.T) {
	res, err := countyService.GetCountyTaxListByName("name")
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}

	if !reflect.DeepEqual(exCountyTaxList, res) {
		t.Error("The GetCountyTaxListById response does not match the expectation.")
	}
}


func TestGetStateById(t *testing.T){
	res, err := stateService.GetStateById(36, "M", 5, 45000)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	if !reflect.DeepEqual(exState, res) {
		t.Error("The GetStateById response does not match the expectation.")
	}
}


func TestGetStateByName(t *testing.T) {
	res, err := stateService.GetStateByName("New York", "M", 5, 45000)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	if !reflect.DeepEqual(exState, res) {
		t.Error("The GetStateById response does not match the expectation.")
	}
}


func TestGetStateList(t *testing.T){ 
	res, err := stateService.GetStateList("commute", 1, true)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}
	// call methods to populate exStateList
	exStateList.AppendToRankedLists(mpList)
	exStateList.SetRankedList(1, true)

	if !reflect.DeepEqual(exStateList, res) {
		t.Error("The GetStateList response does not match the expectation.")
	}
}


func TestGetStateTaxInfoById(t *testing.T){
	res, err := stateService.GetStateTaxInfoById(36)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	exStateTaxInfoId.AppendToOrderedList(bracket1)
	exStateTaxInfoId.AppendToOrderedList(bracket2)

	if !reflect.DeepEqual(exStateTaxInfoId, res) {
		t.Error("The StateTaxInfoById response does not match the expectation.")
	}
}



func TestGetStateTaxInfoByName(t *testing.T){
	res, err := stateService.GetStateTaxInfoByName("New York")
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	exStateTaxInfoName.AppendToOrderedList(bracket1)
	exStateTaxInfoName.AppendToOrderedList(bracket2)

	if !reflect.DeepEqual(exStateTaxInfoName, res) {
		t.Error("The StateTaxInfoByName response does not match the expectation.")
	}
}


func TestGetFederalTaxInfo(t *testing.T) {
	res, err := federalService.GetFederalTaxInfo()
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	federalTaxInfo.AppendToOrderedList(fb1)
	federalTaxInfo.AppendToOrderedList(fb2)
	federalTaxInfo.AppendToOrderedList(fb3)

	if !reflect.DeepEqual(federalTaxInfo, res) {
		t.Error("The GetFederalTaxInfo response does not match the expectation.")
	}
}