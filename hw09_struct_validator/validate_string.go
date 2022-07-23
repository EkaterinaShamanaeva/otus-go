package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	regExp   = "regexp:.+"
	lenField = "len:\\d+"
	inField  = "in:.+"
	value    = "\\w+"
)

var (
	reExp = regexp.MustCompile(regExp)
	reLen = regexp.MustCompile(lenField)
	reIn  = regexp.MustCompile(inField)
	re    = regexp.MustCompile(value)
)

func stringValidate(str string, tag string) []error {
	errSlice := make([]error, 0)

	// check reg exp
	if reExp.FindString(tag) != "" {
		exp := reExp.FindString(tag)[7:]
		// if there is a second condition (...|len:4)
		if i := strings.Index(exp, "|"); i != -1 {
			exp = exp[:i]
		}
		if matched, _ := regexp.MatchString(exp, str); !matched {
			errSlice = append(errSlice, ErrIncorrectContent)
		}
	}

	// check len
	if reLen.FindString(tag) != "" {
		lenRequired, _ := strconv.Atoi(reLen.FindString(tag)[4:])
		if len(str) != lenRequired {
			errSlice = append(errSlice, ErrIncorrectLenOfString)
		}
	}

	// check in
	if reIn.FindString(tag) != "" {
		s := reIn.FindString(tag)[3:]
		if i := strings.Index(s, "|"); i != -1 {
			s = s[:i]
		}
		res := re.FindAllStringSubmatch(s, -1)
		flag := false
		for _, match := range res {
			for _, word := range match {
				if str == word {
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
