package test

/* Mocked dao used in testing builds */

import (
	"github.com/Matthew-Curry/re-region-api/src/apperrors"
	"github.com/Matthew-Curry/re-region-api/src/dao"
)

type DaoMock struct{}

func GetDaoMock() dao.DaoInterface{
	return &DaoMock{}
}

func (d *DaoMock) GetStateCensusData() ([][]interface{}, *apperrors.AppError) {
	res := make([][]interface{}, 0)

	ny := append(make([]interface{}, 0), 36, "New York", 18466230, 8953064, 9513166, 77578, 1381, 17)
	res = append(res, ny)

	return res, nil
}

func (d *DaoMock) GetCountyList(metric string, n int, desc bool) ([][]interface{}, *apperrors.AppError) {
	res := make([][]interface{}, 0)

	ny := append(make([]interface{}, 0), 36061, "New York County", 36, 81)
	res = append(res, ny)

	return res, nil
}

func (d *DaoMock) GetStateTax() ([][]interface{}, *apperrors.AppError) {
	res := make([][]interface{}, 0)

	f1 := append(make([]uint8, 0), 48, 46, 48, 50)
	f2 := append(make([]uint8, 0), 48, 46, 49, 50)

	a1 := append(make([]interface{}, 0), 36, "New York", 2500, 7500, 1500, 3000, 1000, f1, 0, f1, 0)
	res = append(res, a1)

	a2 := append(make([]interface{}, 0), 36, "New York", 2500, 7500, 1500, 3000, 1000, f2, 500, f2, 1000)
	res = append(res, a2)

	return res, nil
}

func (d *DaoMock) GetCountyDataByName(county_name string) ([][]interface{}, *apperrors.AppError) {
	return getMockCounty()
}

func (d *DaoMock) GetCountyDataById(county_id int) ([][]interface{}, *apperrors.AppError) {
	return getMockCounty()
}

func (d *DaoMock) GetFederalTaxData() ([][]interface{}, *apperrors.AppError) {
	res := make([][]interface{}, 0)

	f1 := append(make([]uint8, 0), 48, 46, 49, 48)
	f2 := append(make([]uint8, 0), 48, 46, 49, 50)
	f3 := append(make([]uint8, 0), 48, 46, 50, 50)

	b1 := append(make([]interface{}, 0), f1, 0, 0, 0, 12950, 25900, 19400)
	b2 := append(make([]interface{}, 0), f2, 10275, 20550, 14650, 12950, 25900, 19400)
	b3 := append(make([]interface{}, 0), f3, 41775, 83550, 55900, 12950, 25900, 19400)

	res = append(res, b1, b2, b3)

	return res, nil
}

// used by get county by name and id to return the same mock county
func getMockCounty() ([][]interface{}, *apperrors.AppError) {
	res := make([][]interface{}, 0)

	f := append(make([]uint8, 0), 48, 46, 48, 48)

	ny := append(make([]interface{}, 0), 36061, "New York County", 36, 1628706, 771278, 857428, 93651, 1753, 81, 3376, "New York City", "3.078% - 3.876%", f, f, f, f, f, "0.00%", f, f, f, f, f)
	res = append(res, ny)

	return res, nil
}
