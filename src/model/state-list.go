package model

import (
	"encoding/json"
	"sort"

	"github.com/Matthew-Curry/re-region-api/src/apperrors"
)

type StateList struct {
	Metric_name string
	// store an ascending and descending list to be prepared for either request
	asc_list  []StateMetricPair
	desc_list []StateMetricPair
	// the ranked list returned
	ranked_list []StateMetricPair
}

type StateMetricPair struct {
	State_id     int
	State_name   string
	Metric_value int
}

// constructor for a metric state list, ranked list is private to enforce ordering
func GetMetricStateList(metric string) *StateList {
	return &StateList{Metric_name: metric, ranked_list: []StateMetricPair{}}
}

// method called to set the ranked list and its length
func (s *StateList) SetRankedList(n int, desc bool) {
	if desc {
		s.ranked_list = s.desc_list[:n]
	} else {
		s.ranked_list = s.asc_list[:n]
	}
}

// add pairs to the ranked list in order
func (s *StateList) AppendToRankedLists(metricPair StateMetricPair) {
	// determine indexes to insert to
	ia := sort.Search(len(s.asc_list), func(i int) bool { return s.asc_list[i].Metric_value >= metricPair.Metric_value })
	id := sort.Search(len(s.desc_list), func(i int) bool { return s.desc_list[i].Metric_value <= metricPair.Metric_value })

	s.asc_list = s.insertAtIndex(ia, metricPair, s.asc_list)
	s.desc_list = s.insertAtIndex(id, metricPair, s.desc_list)

}

// insert at both ascending and descending lists
func (s *StateList) insertAtIndex(i int, metricPair StateMetricPair, list []StateMetricPair) []StateMetricPair {
	// if i is the next index, can just append
	if i == len(list) {
		list = append(list, metricPair)
		return list
	}

	// else shift elements in slice at the insertion index. Using append will not allocate extra memory
	list = append(list[:i+1], list[i:]...)

	// overwrite the duplicates value for the insert and return the result
	list[i] = metricPair

	return list
}

// getter method for the controller to be able to marhsall private fields
func (s *StateList) MarshallStateList() ([]byte, *apperrors.AppError) {
	r, err := json.Marshal(struct {
		Metric_name string
		Ranked_list []StateMetricPair
	}{
		Metric_name: s.Metric_name,
		Ranked_list: s.ranked_list,
	})

	if err != nil {
		return nil, apperrors.UnableToMarshall(err)
	}

	return r, nil
}
