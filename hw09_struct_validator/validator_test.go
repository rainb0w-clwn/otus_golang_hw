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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
	Slices struct {
		Strings []string `validate:"len:5|in:first,12345"`
		Ints    []int    `validate:"min:1|max:50|in:25,30,50"`
	}
	PrivateWithTags struct {
		field1 string `validate:"in:private,life"`
		field2 int    `validate:"min:18|max:50"`
	}
	WithMap struct {
		Map map[string]string `validate:"len:5"`
	}
	UnknownValidator struct {
		Field string `validate:"in:private|unknown:123"`
	}
	BadRegexp struct {
		Field string `validate:"regexp:("`
	}
	BadLen struct {
		Field string `validate:"len:a"`
	}
	BadMin struct {
		Field int `validate:"min:a"`
	}
	BadMax struct {
		Field int `validate:"max:a"`
	}
	BadIntIn struct {
		Field []int `validate:"in:a"`
	}
	BadTag struct {
		Field int `validate:"in:abc:abc"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in             interface{}
		expectedFields []string
		expectedErrs   []error
	}{
		{
			in: User{Age: 14, Email: "test@test"},
			expectedFields: []string{
				"ID",
				"Age",
				"Email",
				"Role",
			},
			expectedErrs: []error{
				ErrValidationStringLen,
				ErrValidationIntMin,
				ErrValidationStringRegexp,
				ErrValidationStringIn,
			},
		},
		{in: App{Version: "12345"}},
		{
			in:             App{Version: "3.3.5a"},
			expectedFields: []string{"Version"},
			expectedErrs:   []error{ErrValidationStringLen},
		},
		{in: Token{}},
		{in: Response{Code: 404, Body: "Not found"}},
		{in: Slices{Strings: []string{"12345", "first"}, Ints: []int{25, 30, 50}}},
		{
			in: Slices{
				Strings: []string{"1234567", "first"},
				Ints:    []int{100, 40},
			},
			expectedFields: []string{
				"Strings",
				"Strings",
				"Ints",
				"Ints",
				"Ints",
			},
			expectedErrs: []error{
				ErrValidationStringLen,
				ErrValidationStringIn,
				ErrValidationIntMax,
				ErrValidationIntIn,
				ErrValidationIntIn,
			},
		},
		{in: PrivateWithTags{field1: "LoLKek", field2: 100500}},
		{in: nil, expectedErrs: []error{ErrProgramNotStructure}},
		{in: "1", expectedErrs: []error{ErrProgramNotStructure}},
		{in: WithMap{Map: map[string]string{}}, expectedErrs: []error{ErrProgramUnsupportedFieldKind}},
		{in: UnknownValidator{Field: "private"}, expectedErrs: []error{ErrProgramUnsupportedValidator}},
		{in: BadRegexp{Field: "value"}, expectedErrs: []error{ErrProgramValidatorError}},
		{in: BadLen{Field: "value"}, expectedErrs: []error{ErrProgramValidatorError}},
		{in: BadMin{Field: 5}, expectedErrs: []error{ErrProgramValidatorError}},
		{in: BadMax{Field: 5}, expectedErrs: []error{ErrProgramValidatorError}},
		{in: BadIntIn{Field: []int{5}}, expectedErrs: []error{ErrProgramValidatorError}},
		{in: BadTag{Field: 5}, expectedErrs: []error{ErrProgramValidatorError}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if ok, errs := unwrapErrors(err); ok {
				require.Len(t, errs, len(tt.expectedErrs))
				if len(errs) > 0 {
					for i, err := range errs {
						require.Equal(t, err.Field, tt.expectedFields[i])
						require.ErrorIs(t, err.Err, tt.expectedErrs[i])
					}
				}
			} else {
				for _, e := range tt.expectedErrs {
					require.ErrorIs(t, err, e)
				}
			}
		})
	}
}

func unwrapErrors(err error) (bool, ValidationErrors) {
	var vErrs ValidationErrors
	if !errors.As(err, &vErrs) {
		return false, nil
	}
	return true, vErrs
}
