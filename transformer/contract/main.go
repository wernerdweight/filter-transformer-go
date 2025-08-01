package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type FilterLogic string

type FilterOperator string

const (
	FilterLogicAnd FilterLogic = "and"
	FilterLogicOr  FilterLogic = "or"

	FilterOperatorEqual                   FilterOperator = "eq"
	FilterOperatorNotEqual                FilterOperator = "neq"
	FilterOperatorGreaterThan             FilterOperator = "gt"
	FilterOperatorGreaterThanOrEqual      FilterOperator = "gte"
	FilterOperatorGreaterThanOrEqualOrNil FilterOperator = "gten"
	FilterOperatorLowerThan               FilterOperator = "lt"
	FilterOperatorLowerThanOrEqual        FilterOperator = "lte"
	FilterOperatorLowerThanOrEqualOrNil   FilterOperator = "lten"
	FilterOperatorBegins                  FilterOperator = "begins"
	FilterOperatorContains                FilterOperator = "contains"
	FilterOperatorNotContains             FilterOperator = "not-contains"
	FilterOperatorEnds                    FilterOperator = "ends"
	FilterOperatorIsNil                   FilterOperator = "nil"
	FilterOperatorIsNotNil                FilterOperator = "not-null"
	FilterOperatorIsEmpty                 FilterOperator = "empty"
	FilterOperatorIsNotEmpty              FilterOperator = "not-empty"
	FilterOperatorIn                      FilterOperator = "in"
	FilterOperatorNotIn                   FilterOperator = "not-in"
	FilterOperatorMatchPhrase             FilterOperator = "match-phrase"

	ValidationErrorEmpty           = "empty value"
	ValidationErrorInvalidOperator = "invalid operator"
)

var supportedOperators = []FilterOperator{
	FilterOperatorEqual,
	FilterOperatorNotEqual,
	FilterOperatorGreaterThan,
	FilterOperatorGreaterThanOrEqual,
	FilterOperatorGreaterThanOrEqualOrNil,
	FilterOperatorLowerThan,
	FilterOperatorLowerThanOrEqual,
	FilterOperatorLowerThanOrEqualOrNil,
	FilterOperatorBegins,
	FilterOperatorContains,
	FilterOperatorNotContains,
	FilterOperatorEnds,
	FilterOperatorIsNil,
	FilterOperatorIsNotNil,
	FilterOperatorIsEmpty,
	FilterOperatorIsNotEmpty,
	FilterOperatorIn,
	FilterOperatorNotIn,
	FilterOperatorMatchPhrase,
}

type ValidationError struct {
	Path    string
	Error   string
	Field   string
	Payload any
}

type ValidationFunc func(filterCondition FilterCondition, path string, validationErrors *[]ValidationError)

type FilterCondition struct {
	Field    string
	Operator FilterOperator
	Value    interface{}
}

func (c *FilterCondition) IsNegative() bool {
	return slices.Contains([]FilterOperator{
		FilterOperatorNotEqual,
		FilterOperatorNotContains,
		FilterOperatorIsEmpty,
		FilterOperatorIsNil,
		FilterOperatorNotIn,
	}, c.Operator)
}

func (c *FilterCondition) UnmarshalJSON(data []byte) error {
	var condition struct {
		Field    string
		Operator FilterOperator
		Value    interface{}
	}
	err := json.Unmarshal(data, &condition)
	if err != nil {
		return err
	}
	c.Field = condition.Field
	c.Operator = condition.Operator
	c.Value = condition.Value
	if value, ok := c.Value.(string); ok && c.expectsArray() && c.Value != nil {
		array := strings.Split(value, ",")
		for index, item := range array {
			array[index] = strings.TrimSpace(item)
		}
		c.Value = array
	}
	return nil
}

func (c *FilterCondition) validate(validationErrors *[]ValidationError, path string, validationFunc *ValidationFunc) {
	if c.Field == "" {
		*validationErrors = append(*validationErrors, ValidationError{
			Path:  fmt.Sprintf("%s.field", path),
			Error: ValidationErrorEmpty,
			Field: "field",
		})
	}
	if c.Operator == "" {
		*validationErrors = append(*validationErrors, ValidationError{
			Path:  fmt.Sprintf("%s.operator", path),
			Error: ValidationErrorEmpty,
			Field: "operator",
		})
	}
	if c.Operator != "" && !slices.Contains(supportedOperators, c.Operator) {
		*validationErrors = append(*validationErrors, ValidationError{
			Path:    fmt.Sprintf("%s.operator", path),
			Error:   ValidationErrorInvalidOperator,
			Field:   "operator",
			Payload: string(c.Operator),
		})
	}
	if validationFunc != nil {
		(*validationFunc)(*c, path, validationErrors)
	}
}

func (c *FilterCondition) expectsArray() bool {
	return slices.Contains([]FilterOperator{
		FilterOperatorIn,
		FilterOperatorNotIn,
	}, c.Operator)
}

type FilterConditions struct {
	Conditions []FilterCondition
	Filters    []Filters
}

func (fc *FilterConditions) IsEmpty() bool {
	return len(fc.Conditions) == 0 && len(fc.Filters) == 0
}

func (fc *FilterConditions) UnmarshalJSON(data []byte) error {
	var conditions []FilterCondition
	err := json.Unmarshal(data, &conditions)
	if err == nil && len(conditions) > 0 {
		var filteredConditions []FilterCondition
		for _, condition := range conditions {
			if condition.Field != "" {
				filteredConditions = append(filteredConditions, condition)
			}
		}
		if len(filteredConditions) > 0 {
			fc.Conditions = filteredConditions
		}
	}

	var filters []Filters
	err = json.Unmarshal(data, &filters)
	if err == nil && len(filters) > 0 {
		var filteredFilters []Filters
		for _, filter := range filters {
			if !filter.Conditions.IsEmpty() {
				filteredFilters = append(filteredFilters, filter)
			}
		}
		if len(filteredFilters) > 0 {
			fc.Filters = filteredFilters
		}
	}

	if len(fc.Conditions) > 0 || len(fc.Filters) > 0 {
		return nil
	}

	if len(data) > 0 {
		// not empty, but not a valid filter conditions
		return errors.New("invalid filter conditions")
	}
	return err
}

type Filters struct {
	Logic      FilterLogic
	Conditions FilterConditions
}

func (f *Filters) IsEmpty() bool {
	return f.Logic == "" && f.Conditions.IsEmpty()
}

func (f *Filters) validate(validationErrors *[]ValidationError, path string, validationFunc *ValidationFunc, faIlOnEmpty bool) {
	if f.IsEmpty() && faIlOnEmpty {
		*validationErrors = append(*validationErrors, ValidationError{
			Path:  path,
			Error: ValidationErrorEmpty,
		})
		return
	}
	if f.Logic != "" && f.Logic != FilterLogicAnd && f.Logic != FilterLogicOr {
		*validationErrors = append(*validationErrors, ValidationError{
			Path:    fmt.Sprintf("%s.logic", path),
			Error:   ValidationErrorInvalidOperator,
			Field:   "logic",
			Payload: string(f.Logic),
		})
	}
	if f.Conditions.IsEmpty() && faIlOnEmpty {
		*validationErrors = append(*validationErrors, ValidationError{
			Path:  fmt.Sprintf("%s.conditions", path),
			Error: ValidationErrorEmpty,
			Field: "conditions",
		})
	}
	for index, condition := range f.Conditions.Conditions {
		condition.validate(validationErrors, fmt.Sprintf("%s.conditions.%d", path, index), validationFunc)
	}
	for index, filter := range f.Conditions.Filters {
		filter.validate(validationErrors, fmt.Sprintf("%s.conditions.%d", path, index), validationFunc, faIlOnEmpty)
	}
}

func (f *Filters) Validate(validationFunc *ValidationFunc, faIlOnEmpty bool) []ValidationError {
	var validationErrors []ValidationError
	f.validate(&validationErrors, "root", validationFunc, faIlOnEmpty)
	return validationErrors
}

type InputOutputInterface[T any] interface {
	SetData(data T) error
	GetData() (T, error)
	GetDataString() (string, error)
	GetDataJson() ([]byte, error)
}

type InputOutputType[T any] struct {
	data T
}

func (i *InputOutputType[T]) SetData(data T) error {
	i.data = data
	return nil
}

func (i *InputOutputType[T]) GetData() (T, error) {
	return i.data, nil
}

func NewInputOutputType[T any, IOT InputOutputInterface[T]](data T, i IOT) (IOT, error) {
	err := i.SetData(data)
	if err != nil {
		return i, err
	}
	return i, nil
}

type InputTransformerInterface[T any, IOT InputOutputInterface[T]] interface {
	Transform(input IOT) (Filters, *Error)
}

type OutputTransformerInterface[T any, IOT InputOutputInterface[T]] interface {
	Transform(input Filters) (IOT, *Error)
}
