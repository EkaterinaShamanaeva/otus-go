package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "aaaaaa", // incorrect
				Name:   "Ivan",
				Age:    10, // incorrect
				Email:  "ivanovyandex.ru",
				Role:   "admin",
				Phones: []string{"1111111", "222222"}, // incorrect
				meta:   nil,
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrIncorrectLenOfString,
				}, ValidationError{
					Field: "Age",
					Err:   ErrMin,
				}, ValidationError{
					Field: "Email",
					Err:   ErrIncorrectContent,
				}, ValidationError{
					Field: "Phones",
					Err:   ErrIncorrectLenOfString,
				},
			},
		},
		{
			App{
				Version: "v1.015",
			},
			ValidationErrors{ValidationError{
				Field: "Version",
				Err:   ErrIncorrectLenOfString,
			}},
		},
		{
			Token{
				Header:    []byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"),
				Payload:   []byte("eyJ1c2VySWQiOiJiMDhmODZhZi0zNWRhLTQ4ZjItOGZhYi1jZWYzOTA0NjYwYmQifQ"),
				Signature: []byte("-xN_h82PHVTCMA9vdoHrcZxH-x5mb11y1537t3rGzcM"),
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.NotEmpty(t, err)
				var validationEr ValidationErrors
				if errors.As(err, &validationEr) {
					var expErrs ValidationErrors
					require.ErrorAs(t, tt.expectedErr, &expErrs)
					for j, valErr := range validationEr {
						require.ErrorIs(t, valErr, validationEr[j])
					}
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			}
			_ = tt
		})
	}
}
