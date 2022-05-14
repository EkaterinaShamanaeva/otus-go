package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	rs := []rune(inputString)
	// add '\n' to work in a loop with the last char in inputString
	rs = append(rs, '\n')
	var times int
	var res strings.Builder
	// check if inputString is correct
	if unicode.IsDigit(rs[0]) {
		return "", ErrInvalidString
	}
	// if the string is empty
	if len(rs) == 0 {
		return "", nil
	}
	for i := 1; i < len(rs); i++ {
		// if the current char is a digit, not equal zero
		if unicode.IsDigit(rs[i]) && rs[i] != '0' {
			// if the next char is not a digit, print the previous char <times> times
			if !unicode.IsDigit(rs[i+1]) {
				times, _ = strconv.Atoi(string(rs[i]))
				res.WriteString(strings.Repeat(string(rs[i-1]), times))
				// else - return error
			} else {
				return "", ErrInvalidString
			}
		} else if !unicode.IsDigit(rs[i-1]) && rs[i] != '0' {
			res.WriteString(string(rs[i-1]))
		}
	}
	return res.String(), nil
}
