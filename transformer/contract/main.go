package contract

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
	FilterOperatorIsNotNil                FilterOperator = "not-nil"
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

type FilterConditions struct {
	Conditions []FilterCondition
	Filters    []Filters
}

type Filters struct {
	Logic      FilterLogic
	Conditions FilterConditions
}

type InputType interface {
	GetData() (interface{}, error)
}

type OutputType interface {
	SetData(data interface{}) error
}

type InputTransformerInterface[T InputType] interface {
	Transform(input T) (Filters, error)
}

type OutputTransformerInterface[T OutputType] interface {
	Transform(input Filters) (*T, error)
}
