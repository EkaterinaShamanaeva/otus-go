package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, error: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

var (
	ErrNotStruct            = errors.New("interface{} is not a struct")
	ErrIncorrectLenOfString = errors.New("incorrect len of the string")
	ErrIncorrectContent     = errors.New("does not match regular expression")
	ErrIncorrectInCond      = errors.New("not included in the set")
	ErrMax                  = errors.New("value is bigger then max allowable value")
	ErrMin                  = errors.New("value is lower then min allowable value")
)

func (v ValidationErrors) Error() string {
	if len(v) < 1 {
		return "there are not validation errors"
	}
	buf := bytes.NewBufferString("")
	for i := 0; i < len(v); i++ {
		text := fmt.Sprintf("Field: %s; Error: %s",
			v[i].Field, v[i].Err)
		buf.WriteString(text)
		buf.WriteString("\n")
	}
	return buf.String()
}

func Validate(v interface{}) error {
	errAll := make(ValidationErrors, 0)

	vType := reflect.TypeOf(v)
	if vType.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	vValue := reflect.ValueOf(v)

	for i := 0; i < vType.NumField(); i++ {
		f := vType.Field(i)
		if f.Tag != "" {
			if f.Tag.Get("validate") != "" {
				tag := f.Tag.Get("validate")
				switch { // switch f.Type().Kind() -> switch (corrected after lint)
				case f.Type.Kind() == reflect.String:
					err := stringPrepareAndValidate(f, vValue.Field(i), tag)
					errAll = append(errAll, err...)
				case f.Type.Kind() == reflect.Int:
					err := intPrepareAndValidate(f, vValue.Field(i), tag)
					errAll = append(errAll, err...)
				case f.Type.Kind() == reflect.Slice:
					err := sliceValidate(f, vValue.Field(i), tag)
					errAll = append(errAll, err...)
				default:
				}
			}
		}
	}
	e := errAll.Error()
	fmt.Println(e)
	if len(errAll) > 0 {
		return errAll
	}
	return nil
}

func stringPrepareAndValidate(f reflect.StructField, fieldValue reflect.Value, tag string) ValidationErrors {
	errsString := make(ValidationErrors, 0)
	str := fieldValue.String()
	if err := stringValidate(str, tag); err != nil {
		for k := range err {
			errsString = append(errsString, ValidationError{f.Name, err[k]})
		}
	}
	return errsString
}

func intPrepareAndValidate(f reflect.StructField, fieldValue reflect.Value, tag string) ValidationErrors {
	errsInt := make(ValidationErrors, 0)
	n := fieldValue.Int()
	if err := intValidate(int(n), tag); err != nil {
		for k := range err {
			errsInt = append(errsInt, ValidationError{f.Name, err[k]})
		}
	}
	return errsInt
}

func sliceValidate(f reflect.StructField, sl reflect.Value, tag string) ValidationErrors {
	errsSlice := make(ValidationErrors, 0)

	for j := 0; j < sl.Len(); j++ {
		elem := sl.Index(j)
		if elem.Kind() == reflect.String {
			if err := stringValidate(elem.String(), tag); err != nil {
				for k := range err {
					errsSlice = append(errsSlice, ValidationError{f.Name, err[k]})
				}
			}
		} else if elem.Kind() == reflect.Int {
			if err := intValidate(int(elem.Int()), tag); err != nil {
				for k := range err {
					errsSlice = append(errsSlice, ValidationError{f.Name, err[k]})
				}
			}
		}
	}
	return errsSlice
}
