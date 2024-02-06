package contract

import (
	"encoding/json"
	"errors"
	"slices"
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
)

type FilterCondition struct {
	Field    string
	Operator FilterOperator
	Value    interface{}
}

func (c FilterCondition) IsNegative() bool {
	return slices.Contains([]FilterOperator{
		FilterOperatorNotEqual,
		FilterOperatorNotContains,
		FilterOperatorIsEmpty,
		FilterOperatorIsNil,
		FilterOperatorNotIn,
	}, c.Operator)
}

type FilterConditions struct {
	Conditions []FilterCondition
	Filters    []Filters
}

func (fc *FilterConditions) IsEmtpy() bool {
	return len(fc.Conditions) == 0 && len(fc.Filters) == 0
}

func (fc *FilterConditions) UnmarshalJSON(data []byte) error {
	var conditions []FilterCondition
	err := json.Unmarshal(data, &conditions)
	if err == nil && len(conditions) > 0 && conditions[0].Field != "" {
		fc.Conditions = conditions
		return nil
	}
	var filters []Filters
	err = json.Unmarshal(data, &filters)
	if err == nil && len(filters) > 0 && !filters[0].Conditions.IsEmtpy() {
		fc.Filters = filters
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
	return f.Logic == "" && f.Conditions.IsEmtpy()
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
	// TODO: add optional validators for fields (filters.conditions.conditions[i].field)
	// TODO: also for values (filters.conditions.conditions[i].value)
	Transform(input IOT) (Filters, error)
}

type OutputTransformerInterface[T any, IOT InputOutputInterface[T]] interface {
	Transform(input Filters) (IOT, error)
}
