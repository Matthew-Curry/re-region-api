package model

import (
	"encoding/json"
)

type State struct {
	State_id int
	State_name string
	// state level census metrics
	Pop int
    Male_pop int
    Female_pop int
    Median_income int
    Average_rent int
    Commute int
    // the tax metrics
    Total_tax int
    State_tax int
	Federal_tax int
}

// marshaller for controller
func (s *State) MarshallState() ([]byte, error) {
	return json.Marshal(s)
}