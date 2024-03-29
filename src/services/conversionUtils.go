package services

import (
	"strconv"
)

/* Hold functions used across services to convert data retrieved from the dao to types needed for the model structs */

// return var as int
func readAsInt(i interface{}) int {
	switch i.(type) {
	case int64:
		return int(i.(int64))
	case int:
		return i.(int)
	}

	return 0
}

// return var as string
func readAsString(i interface{}) string {
	return i.(string)
}

// return var as float
func readAsFloat(i interface{}) float64 {
	s := string(i.([]uint8))
	f, _ := strconv.ParseFloat(s, 10)

	return f
}
