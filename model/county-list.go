package model

import (
	"encoding/json"

	"github.com/Matthew-Curry/re-region-api/apperrors"
)

type CountyList struct {
	Metric_name string
	Ranked_list []CountyMetricPair
}

type CountyMetricPair struct {
	County_id int
	County_name string
	State_id int
	State_name string
	Metric_value int
}

// constructor for a metric county list
func GetMetricCountyList(metric string) *CountyList {
	return &CountyList {Metric_name: metric, Ranked_list: []CountyMetricPair{}}
}

// getter method for the controller to be able to marhsall private fields
func (c *CountyList) MarshallCountyList() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(c)

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}