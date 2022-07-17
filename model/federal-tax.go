package model

import (
    "sort"
	"encoding/json"

	"github.com/Matthew-Curry/re-region-api/apperrors"
)

type FederalTaxInfo struct {
	// deduciton/exemption info common across the brackets
    Single_deduction int
    Married_deduction int
    Head_deduction int
	// list of brackets and rates
	bracket_list []FederalBracket
}

type FederalBracket struct {
	Rate float64
    Single_bracket int
    Married_bracket int
	Head_bracket int
}

// constructor for a FederalTaxInfo, bracket list is private to enforce ordering
func GetFederalTaxInfo(sd, md, hd int, bracketList []FederalBracket) *FederalTaxInfo {
	return &FederalTaxInfo {
		Single_deduction :sd,
		Married_deduction :md,
		Head_deduction :hd,
        bracket_list: bracketList}
}

// public method to use private bracket list to get the single state tax liability
func (f *FederalTaxInfo) GetSingleTaxLiability(income int) int {
    i := sort.Search(len(f.bracket_list), func(i int) bool { return f.bracket_list[i].Single_bracket >= income })
	r := f.bracket_list[i].Rate
	return int(float64(income) * r)
}

// public method to use private bracket list to get the state tax liability
func (f *FederalTaxInfo) GetMarriedTaxLiability(income int) int {
    i := sort.Search(len(f.bracket_list), func(i int) bool { return f.bracket_list[i].Married_bracket >= income })
	r := f.bracket_list[i].Rate
	return int(float64(income) * r)
}

// public method to use private bracket list to get the head tax liability
func (f *FederalTaxInfo) GetHeadTaxLiability(income int) int {
    i := sort.Search(len(f.bracket_list), func(i int) bool { return f.bracket_list[i].Head_bracket >= income })
	r := f.bracket_list[i].Rate
	return int(float64(income) * r)
}

// add pairs to the orderd list
func (f *FederalTaxInfo) AppendToOrderedList(bracket FederalBracket) {
	// determine index to insert to
	i := sort.Search(len(f.bracket_list), func(i int) bool { return f.bracket_list[i].Rate >= bracket.Rate })
	// if i is the next index, can just append
	if i == len(f.bracket_list) {
		f.bracket_list = append(f.bracket_list, bracket)
	}

	// else shift elements in slice at the insertion index. Using append will not allocate extra memory
	// when cap(mpSlice) > len(mpSlice)
	f.bracket_list = append(f.bracket_list[:i+1], f.bracket_list[i:]...)

	// overwrite the duplicates value for the insert and return the result
	f.bracket_list[i] = bracket
}

// getter method for the controller to be able to marhsall private fields
func (f *FederalTaxInfo) MarshallFederalTaxInfo() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(struct{
		Single_deduction int
		Married_deduction int
		Head_deduction int
		Bracket_list []FederalBracket
    }{
        Single_deduction :f.Single_deduction,
		Married_deduction :f.Married_deduction,
		Head_deduction: f.Head_deduction,
		Bracket_list: f.bracket_list,
    })

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}