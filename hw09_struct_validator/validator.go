package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("ValidationError: field \"%s\" not valid %s\n", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for _, err := range v {
		sb.WriteString(err.Error())
	}
	return sb.String()
}

var (
	ErrProgramNotStructure         = errors.New("argument is not a struct")
	ErrProgramUnsupportedFieldKind = errors.New("unsupported kind present")
	ErrProgramUnsupportedValidator = errors.New("unsupported validator present")
	ErrProgramValidatorError       = errors.New("ProgramError")
	ErrValidationIntMin            = errors.New("must be greater than")
	ErrValidationIntMax            = errors.New("must be less than")
	ErrValidationIntIn             = errors.New("must be any of")
	ErrValidationStringLen         = errors.New("must be length of")
	ErrValidationStringRegexp      = errors.New("should match")
	ErrValidationStringIn          = errors.New("should be any of")
)

const tagName = "validate"

func Validate(v interface{}) error {
	if v == nil {
		return fmt.Errorf("%w", ErrProgramNotStructure)
	}
	inputValue := reflect.ValueOf(v)
	inputType := inputValue.Type()
	if inputType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %w", ErrProgramValidatorError, ErrProgramNotStructure)
	}
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < inputType.NumField(); i++ {
		fieldType := inputType.Field(i)
		fieldTag := fieldType.Tag.Get(tagName)
		if len(fieldTag) == 0 {
			continue
		}
		fieldName := fieldType.Name
		fieldValue := inputValue.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}
		tags := *splitValidatorTags(fieldTag)
		if len(tags) == 0 {
			continue
		}
		var (
			vE []error
			pE error
		)
		switch fieldType.Type.Kind() { //nolint:exhaustive
		case reflect.Int:
			vE, pE = ValidateInt(fieldValue, tags)
		case reflect.String:
			vE, pE = ValidateString(fieldValue, tags)
		case reflect.Slice:
			vE, pE = ValidateSlice(fieldType.Type.Elem().Kind(), fieldValue, tags)
		default:
			pE = ErrProgramUnsupportedFieldKind
		}
		if pE != nil {
			return fmt.Errorf("%w: %w", ErrProgramValidatorError, pE)
		} else if len(vE) != 0 {
			for _, fieldE := range vE {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fieldE,
				})
			}
		}
	}
	return validationErrors
}

func splitValidatorTags(fT string) *[]string {
	r := strings.Split(fT, "|")
	return &r
}

func splitValidatorTagValue(v string) (*string, *string, error) {
	s := strings.Split(v, ":")
	if len(s) != 2 {
		return nil, nil, errors.New("validator must have one value")
	}
	return &s[0], &s[1], nil
}

func ValidateInt(fieldValue reflect.Value, tags []string) (validationErrors []error, programError error) {
	fV := int(fieldValue.Int())
	for _, v := range tags {
		vN, vV, err := splitValidatorTagValue(v)
		if err != nil {
			return nil, fmt.Errorf("for %s: %w", v, err)
		}
		switch *vN {
		case "min":
			val, err := strconv.Atoi(*vV)
			if err != nil {
				return nil, fmt.Errorf("for %s: %w", *vV, err)
			}
			if fV < val {
				validationErrors = append(validationErrors, fmt.Errorf(
					"by \"%s\": value \"%d\" %w \"%s\"", *vN, fV, ErrValidationIntMin, *vV,
				),
				)
			}

		case "max":
			vVInt, err := strconv.Atoi(*vV)
			if err != nil {
				return nil, fmt.Errorf("for %s: %w", *vV, err)
			}
			if fV > vVInt {
				validationErrors = append(validationErrors, fmt.Errorf(
					"by \"%s\": value \"%d\" %w \"%s\"", *vN, fV, ErrValidationIntMax, *vV,
				),
				)
			}
		case "in":
			vVIn := strings.Split(*vV, ",")
			var in bool
			for _, inS := range vVIn {
				inV, err := strconv.Atoi(inS)
				if err != nil {
					return nil, fmt.Errorf("for %s: %w", *vV, err)
				}
				in = in || inV == fV
				// continue to search for error parsing
			}
			if !in {
				validationErrors = append(validationErrors, fmt.Errorf(
					"by \"%s\": value \"%d\" %w \"%s\"", *vN, fV, ErrValidationIntIn, *vV,
				),
				)
			}
		default:
			return nil, fmt.Errorf("for %s: %w", v, ErrProgramUnsupportedValidator)
		}
	}
	return validationErrors, nil
}

func ValidateString(fieldValue reflect.Value, tags []string) (validationErrors []error, programError error) {
	fV := fieldValue.String()
	for _, v := range tags {
		vN, vV, err := splitValidatorTagValue(v)
		if err != nil {
			return nil, fmt.Errorf("for %s: %w", v, err)
		}
		switch *vN {
		case "len":
			val, err := strconv.Atoi(*vV)
			if err != nil {
				return nil, fmt.Errorf("for %s: %w", *vV, err)
			}
			if len(fV) != val {
				validationErrors = append(validationErrors, fmt.Errorf(
					"by \"%s\": value \"%s\" %w \"%s\"", *vN, fV, ErrValidationStringLen, *vV,
				),
				)
			}
		case "regexp":
			r, err := regexp.Compile(*vV)
			if err != nil {
				return nil, fmt.Errorf("for %s: %w", *vV, err)
			}
			if !r.MatchString(fV) {
				validationErrors = append(validationErrors, fmt.Errorf(
					"by \"%s\": value \"%s\" %w \"%s\"", *vN, fV, ErrValidationStringRegexp, *vV,
				),
				)
			}
		case "in":
			vVIn := strings.Split(*vV, ",")
			var in bool
			for _, inS := range vVIn {
				if in = inS == fV; in {
					break
				}
			}
			if !in {
				validationErrors = append(validationErrors, fmt.Errorf(
					"by \"%s\": value \"%s\" %w \"%s\"", *vN, fV, ErrValidationStringIn, *vV,
				),
				)
			}
		default:
			return nil, fmt.Errorf("for %s: %w", v, ErrProgramUnsupportedValidator)
		}
	}
	return validationErrors, nil
}

func ValidateSlice(
	fieldKind reflect.Kind,
	fieldValue reflect.Value,
	tags []string,
) (
	validationErrors []error,
	programError error,
) {
	var validator func(reflect.Value, []string) ([]error, error)
	switch fieldKind { //nolint:exhaustive
	case reflect.Int:
		validator = ValidateInt
	case reflect.String:
		validator = ValidateString
	default:
		return nil, ErrProgramUnsupportedFieldKind
	}
	for i := 0; i < fieldValue.Len(); i++ {
		vE, pE := validator(fieldValue.Index(i), tags)
		if pE != nil {
			return nil, pE
		}
		validationErrors = append(validationErrors, vE...)
	}
	return validationErrors, nil
}
