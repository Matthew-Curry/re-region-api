package model

import (
	"sort"
	"encoding/json"
)

type CountyList struct {
	Metric_name string
	ranked_list []CountyMetricPair
}

type CountyMetricPair struct {
	County_id int
	County_name string
	State_id int
	State_name string
	Metric_value int
}

// constructor for a metric county list, ranked list is private to enforce ordering
func GetMetricCountyList(metric string) *CountyList {
	return &CountyList {Metric_name: metric, ranked_list: []CountyMetricPair{}}
}

// add pairs to the ranked list in order
func (c *CountyList) AppendToRankedList(metricPair CountyMetricPair) {
	// determine index to insert to
	i := sort.Search(len(c.ranked_list), func(i int) bool { return c.ranked_list[i].Metric_value >= metricPair.Metric_value })
	// if i is the next index, can just append
	if i == len(c.ranked_list) {
		c.ranked_list = append(c.ranked_list, metricPair)
	}

	// else shift elements in slice at the insertion index. Using append will not allocate extra memory
	// when cap(c.ranked_list) > len(c.ranked_list)
	c.ranked_list = append(c.ranked_list[:i+1], c.ranked_list[i:]...)

	// overwrite the duplicates value for the insert and return the result
	c.ranked_list[i] = metricPair
}

// getter method for the controller to be able to marhsall private fields
func (c *CountyList) MarshallCountyList() ([]byte, error) {
	return json.Marshal(struct{
        Metric_name string
		Ranked_list []CountyMetricPair
    }{
        Metric_name: c.Metric_name,
		Ranked_list :c.ranked_list,
    })
}