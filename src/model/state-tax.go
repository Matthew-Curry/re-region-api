package model

import (
	"encoding/json"
	"sort"

	"github.com/Matthew-Curry/re-region-api/src/apperrors"
)

type StateTaxInfo struct {
	State_id   int
	State_name string
	// deduciton/exemption info common across the brackets
	Single_deduction    int
	Married_deduction   int
	Single_exemption    int
	Married_exemption   int
	Dependent_exemption int
	// list of brackets and rates
	bracket_list []StateBracket
}

type StateBracket struct {
	Single_rate     float64
	Single_bracket  int
	Married_rate    float64
	Married_bracket int
}

// constructor for StateTaxInfo, bracket list is private to enforce ordering
func GetStateTaxInfo(si int, sn string, sd, md, se, me, de int) *StateTaxInfo {
	return &StateTaxInfo{State_id: si,
		State_name:          sn,
		Single_deduction:    sd,
		Married_deduction:   md,
		Single_exemption:    se,
		Married_exemption:   me,
		Dependent_exemption: de,
		bracket_list:        []StateBracket{}}
}

// public method to use private bracket list to get the single state tax liability
func (s *StateTaxInfo) GetSingleTaxLiability(income int) int {
	i := sort.Search(len(s.bracket_list), func(i int) bool { return s.bracket_list[i].Single_bracket >= income })
	// subtract one if in highest bracket
	if i == len(s.bracket_list) {
		i = i - 1
	}
	r := s.bracket_list[i].Single_rate

	return int(float64(income) * r)
}

// public method to use private bracket list to get the state tax liability
func (s *StateTaxInfo) GetMarriedTaxLiability(income int) int {
	i := sort.Search(len(s.bracket_list), func(i int) bool { return s.bracket_list[i].Married_bracket >= income })
	// subtract one if in highest bracket
	if i == len(s.bracket_list) {
		i = i - 1
	}
	r := s.bracket_list[i].Married_rate
	return int(float64(income) * r)
}

// add pairs to the orderd list
func (s *StateTaxInfo) AppendToOrderedList(bracket StateBracket) {
	// determine index to insert to
	i := sort.Search(len(s.bracket_list), func(i int) bool { return s.bracket_list[i].Single_rate >= bracket.Single_rate })
	// if i is the next index, can just append
	if i == len(s.bracket_list) {
		s.bracket_list = append(s.bracket_list, bracket)
		return
	}

	// else shift elements in slice at the insertion index. Using append will not allocate extra memory
	s.bracket_list = append(s.bracket_list[:i+1], s.bracket_list[i:]...)

	// overwrite the duplicates value for the insert and return the result
	s.bracket_list[i] = bracket
}

// marshaller for the controller to be able to marhsall private fields
func (s *StateTaxInfo) MarshallStateTaxInfo() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(struct {
		State_id            int
		State_name          string
		Single_deduction    int
		Married_deduction   int
		Single_exemption    int
		Married_exemption   int
		Dependent_exemption int
		Bracket_list        []StateBracket
	}{
		State_id:            s.State_id,
		State_name:          s.State_name,
		Single_deduction:    s.Single_deduction,
		Married_deduction:   s.Married_deduction,
		Single_exemption:    s.Single_exemption,
		Married_exemption:   s.Married_exemption,
		Dependent_exemption: s.Dependent_exemption,
		Bracket_list:        s.bracket_list,
	})

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}
