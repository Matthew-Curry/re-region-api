package services

/* Implementation of the Re-Region API federal service */

import (
	"github.com/Matthew-Curry/re-region-api/apperrors"
	"github.com/Matthew-Curry/re-region-api/dao"
	"github.com/Matthew-Curry/re-region-api/model"
)

// indexes from the federal response
const (
	FEDERAL_RATE = iota          
	FEDERAL_SINGLE_BRACKET   
	FEDERAL_MARRIED_BRACKET  
	FEDERAL_HEAD_BRACKET      
	FEDERAL_STANDARD_DEDUCTION 
	FEDERAL_MARRIED_DEDUCTION 
	FEDERAL_HEAD_DEDUCTION   
)

type FederalServiceImpl struct {
	federalTaxInfo *model.FederalTaxInfo
}

// constructor to return this implementation of the federal service
func GetFederalServiceImpl(daoImpl dao.DaoInterface) (FederalServiceInterface, *apperrors.AppError) {
	// use dao to retrieve the federaltaxlist
	federalTaxList, err := daoImpl.GetFederalTaxData()
	if err != nil {
		return nil, err
	}

	federalTaxInfo := buildCachedResponse(federalTaxList)

	return &FederalServiceImpl{federalTaxInfo: federalTaxInfo}, nil

}

// public method to get overall federal tax information
func buildCachedResponse(federalTaxList [][]interface{}) *model.FederalTaxInfo {
	// get constant attributes from first record
	sd := federalTaxList[0][FEDERAL_STANDARD_DEDUCTION]
	md := federalTaxList[0][FEDERAL_MARRIED_DEDUCTION]
	hd := federalTaxList[0][FEDERAL_HEAD_DEDUCTION]

	// pass over brackets to form the bracket list
	bracketList := []model.FederalBracket{}
	for _, row := range federalTaxList {
		r := row[FEDERAL_RATE]
		sb := row[FEDERAL_SINGLE_BRACKET]
		mb := row[FEDERAL_MARRIED_BRACKET]
		hb := row[FEDERAL_HEAD_BRACKET]
		b := model.FederalBracket{Rate: r.(float64), Single_bracket: sb.(int), Married_bracket: mb.(int), Head_bracket: hb.(int)}
		bracketList = append(bracketList, b)
	}

	// form complete response
	return model.GetFederalTaxInfo(sd.(int), md.(int), hd.(int))
}

// public method to return the federal tax information
func (f *FederalServiceImpl) GetFederalTaxInfo() (*model.FederalTaxInfo, *apperrors.AppError) {
	if f.federalTaxInfo == nil {
		return nil, apperrors.EmptyFederalCache()
	}
	return f.federalTaxInfo, nil
}

// public method to get overall federal tax liability
func (f *FederalServiceImpl) GetFederalLiability(filingStatus model.FilingStatus, dependents int, income int) int {
	// use filing status to determine state deduction and exemption
	switch filingStatus {
	case model.Head:
		income = income - f.federalTaxInfo.Single_deduction
		return f.federalTaxInfo.GetHeadTaxLiability(income)
	case model.Single:
		income = income - f.federalTaxInfo.Head_deduction
		return f.federalTaxInfo.GetSingleTaxLiability(income)
	case model.Married:
		income = income - f.federalTaxInfo.Married_deduction
		return f.federalTaxInfo.GetMarriedTaxLiability(income)
	}

	return 0

}
