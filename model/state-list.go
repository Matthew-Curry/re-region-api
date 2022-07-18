package model

import (
	"sort"
	"encoding/json"

	"github.com/Matthew-Curry/re-region-api/apperrors"
)

type StateList struct {
	Metric_name string
	ranked_list []StateMetricPair
}

type StateMetricPair struct {
	State_id int
	State_name string
	Metric_value int
}

// constructor for a metric state list, ranked list is private to enforce ordering
func GetMetricStateList(metric string) *StateList {
	return &StateList {Metric_name: metric, ranked_list: []StateMetricPair{}}
}

// add pairs to the ranked list in order
func (s *StateList) AppendToRankedList(metricPair StateMetricPair) {
	// determine index to insert to
	i := sort.Search(len(s.ranked_list), func(i int) bool { return s.ranked_list[i].Metric_value >= metricPair.Metric_value })
	// if i is the next index, can just append
	if i == len(s.ranked_list) {
		s.ranked_list = append(s.ranked_list, metricPair)
		return
	}

	// else shift elements in slice at the insertion index. Using append will not allocate extra memory
	// when cap(s.ranked_list) > len(s.ranked_list)
	s.ranked_list = append(s.ranked_list[:i+1], s.ranked_list[i:]...)

	// overwrite the duplicates value for the insert and return the result
	s.ranked_list[i] = metricPair
}

// getter method for the controller to be able to marhsall private fields
func (s *StateList) MarshallStateList() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(struct{
        Metric_name string
		Ranked_list []StateMetricPair
    }{
        Metric_name: s.Metric_name,
		Ranked_list :s.ranked_list,
    })

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}

