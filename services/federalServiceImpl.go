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
	logger.Info("Getting federal tax data from data access layer")
	federalTaxList, err := daoImpl.GetFederalTaxData()
	if err != nil {
		return nil, err
	}

	logger.Info("Caching the response")
	federalTaxInfo := buildCachedResponse(federalTaxList)
	logger.Info("Federal tax cache created")

	return &FederalServiceImpl{federalTaxInfo: federalTaxInfo}, nil

}

// public method to get overall federal tax information
func buildCachedResponse(federalTaxList [][]interface{}) *model.FederalTaxInfo {
	// get constant attributes from first record
	sd := readAsInt(federalTaxList[0][FEDERAL_STANDARD_DEDUCTION])
	md := readAsInt(federalTaxList[0][FEDERAL_MARRIED_DEDUCTION])
	hd := readAsInt(federalTaxList[0][FEDERAL_HEAD_DEDUCTION])

	// pass over brackets to form the bracket list
	bracketList := []model.FederalBracket{}
	for _, row := range federalTaxList {
		// convert rate to float
		r := readAsFloat(row[FEDERAL_RATE])
		// other fields
		sb := readAsInt(row[FEDERAL_SINGLE_BRACKET])
		mb := readAsInt(row[FEDERAL_MARRIED_BRACKET])
		hb := readAsInt(row[FEDERAL_HEAD_BRACKET])
		b := model.FederalBracket{Rate: r, Single_bracket: sb, Married_bracket: mb, Head_bracket: hb}
		bracketList = append(bracketList, b)
	}

	// form complete response
	return model.GetFederalTaxInfo(sd, md, hd, bracketList)
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
		income = getTaxableIncome(income, f.federalTaxInfo.Single_deduction, 0, 0)
		return f.federalTaxInfo.GetHeadTaxLiability(income)
	case model.Single:
		income = getTaxableIncome(income, f.federalTaxInfo.Head_deduction, 0, 0)
		return f.federalTaxInfo.GetSingleTaxLiability(income)
	case model.Married:
		income = getTaxableIncome(income, f.federalTaxInfo.Married_deduction, 0, 0)
		return f.federalTaxInfo.GetMarriedTaxLiability(income)
	}

	return 0

}
