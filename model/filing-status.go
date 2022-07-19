package model

import (
	"strings"
	"errors"
	"fmt"
)

// "enum" for filing status request param
type FilingStatus string

const (
	Head    FilingStatus = "h"
	Single               = "s"
	Married              = "m"
)

func ToFilingStatus(s string) (FilingStatus, error){
	s = strings.ToLower(s)
	if s != "h" && s != "s" && s != "m" {
		return FilingStatus("s"), errors.New(fmt.Sprintf("%s is not a valid filing status.", s))
	} 

	return FilingStatus(s), nil
}