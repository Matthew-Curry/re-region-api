package test

/* Testing suite for Re-Region services */

import (
	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	"github.com/Matthew-Curry/re-region-api/src/dao"
	"github.com/Matthew-Curry/re-region-api/src/services"

	"log"
	"reflect"

	"testing"

	_ "github.com/Matthew-Curry/re-region-api/src/model"
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

func assertEqual(t *testing.T, method string, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("The %s response does not match the expectation.", method)
	}
}

func TestGetCountyById(t *testing.T) {

	res, err := countyService.GetCountyById(5, "S", true, 4, 45000)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}
	
	assertEqual(t, "GetCountyId", res, exCounty)
}


func TestGetCountyByName(t *testing.T) {
	res, err := countyService.GetCountyByName("name", "S", true, 4, 45000)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}
	
	assertEqual(t, "GetCountyName", res, exCounty)
}



func TestGetCountyList(t *testing.T){
	res, err := countyService.GetCountyList("metric", 5, true)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}

	assertEqual(t, "GetCountyList", res, exCountyList)
}

func TestGetCountyTaxListById(t *testing.T){
	res, err := countyService.GetCountyTaxListById(5)
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}

	assertEqual(t, "GetCountyTaxListById", res, exCountyTaxList)
}

func TestGetCountyTaxListByName(t *testing.T) {
	res, err := countyService.GetCountyTaxListByName("name")
	if err != nil{
		t.Error("Error recieved from the county service.", err)
	}

	assertEqual(t, "GetCountyTaxListByName", res, exCountyTaxList)
}


func TestGetStateById(t *testing.T){
	res, err := stateService.GetStateById(36, "M", 5, 45000)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	assertEqual(t, "GetStateById", res, exState)
}


func TestGetStateByName(t *testing.T) {
	res, err := stateService.GetStateByName("New York", "M", 5, 45000)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	assertEqual(t, "GetStateByName", res, exState)
}


func TestGetStateList(t *testing.T){ 
	res, err := stateService.GetStateList("commute", 1, true)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}
	// call methods to populate exStateList
	exStateList.AppendToRankedLists(mpList)
	exStateList.SetRankedList(1, true)

	assertEqual(t, "GetStateList", res, exStateList)
}


func TestGetStateTaxInfoById(t *testing.T){
	res, err := stateService.GetStateTaxInfoById(36)
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	exStateTaxInfoId.AppendToOrderedList(bracket1)
	exStateTaxInfoId.AppendToOrderedList(bracket2)

	assertEqual(t, "GetStateTaxInfoById", res, exStateTaxInfoId)
}



func TestGetStateTaxInfoByName(t *testing.T){
	res, err := stateService.GetStateTaxInfoByName("New York")
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	exStateTaxInfoName.AppendToOrderedList(bracket1)
	exStateTaxInfoName.AppendToOrderedList(bracket2)

	assertEqual(t, "GetStateTaxInfoByName", res, exStateTaxInfoName)
}


func TestGetFederalTaxInfo(t *testing.T) {
	res, err := federalService.GetFederalTaxInfo()
	if err != nil{
		t.Error("Error recieved from the state service.", err)
	}

	federalTaxInfo.AppendToOrderedList(fb1)
	federalTaxInfo.AppendToOrderedList(fb2)
	federalTaxInfo.AppendToOrderedList(fb3)

	assertEqual(t, "GetFederalTaxInfo", res, federalTaxInfo)
}