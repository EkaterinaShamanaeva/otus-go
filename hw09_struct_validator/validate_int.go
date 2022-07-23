package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	maxValue = "max:\\d+"
	minValue = "min:\\d+"
	inLimit  = "in:\\d+,.+"
	valueInt = "\\d+"
)

var (
	reMax     = regexp.MustCompile(maxValue)
	reMin     = regexp.MustCompile(minValue)
	reInLimit = regexp.MustCompile(inLimit)
	reInt     = regexp.MustCompile(valueInt)
)

func intValidate(n int, tag string) []error {
	errSlice := make([]error, 0)

	// compare with max
	if reMax.FindString(tag) != "" {
		maxRequired, _ := strconv.Atoi(reMax.FindString(tag)[4:])
		if n >= maxRequired {
			errSlice = append(errSlice, ErrMax)
		}
	}

	// compare with min
	if reMin.FindString(tag) != "" {
		minRequired, _ := strconv.Atoi(reMin.FindString(tag)[4:])
		if n <= minRequired {
			errSlice = append(errSlice, ErrMin)
		}
	}

	// check in
	if reInLimit.FindString(tag) != "" {
		s := reInLimit.FindString(tag)[3:]
		if i := strings.Index(s, "|"); i != -1 {
			s = s[:i]
		}

		res := reInt.FindAllStringSubmatch(s, -1)

		flag := false
		for _, match := range res {
			for _, val := range match {
				v, _ := strconv.Atoi(val)
				if n == v {
					flag = true
				}
			}
		}
		if !flag {
			errSlice = append(errSlice, ErrIncorrectInCond)
		}
	}

	return errSlice
}
