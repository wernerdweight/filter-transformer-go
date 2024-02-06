package output

import (
	"encoding/json"
	"fmt"
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
)

type ElasticOutput struct {
	contract.InputOutputType[map[string]any]
}

func (o *ElasticOutput) GetDataJson() ([]byte, error) {
	rawData, err := o.GetData()
	if err != nil {
		return nil, err
	}
	return json.Marshal(rawData)
}

func (o *ElasticOutput) GetDataString() (string, error) {
	rawData, err := o.GetDataJson()
	if err != nil {
		return "", err
	}
	return string(rawData), nil
}

var conditionResolvers = map[contract.FilterOperator]func(contract.FilterCondition) map[string]any{
	contract.FilterOperatorEqual: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"term": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorNotEqual: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"term": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorGreaterThan: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"range": map[string]any{
				condition.Field: map[string]any{
					"gt": condition.Value,
				},
			},
		}
	},
	contract.FilterOperatorGreaterThanOrEqual: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"range": map[string]any{
				condition.Field: map[string]any{
					"gte": condition.Value,
				},
			},
		}
	},
	contract.FilterOperatorGreaterThanOrEqualOrNil: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"range": map[string]any{
				condition.Field: map[string]any{
					"gte": condition.Value,
				},
			},
		}
	},
	contract.FilterOperatorLowerThan: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"range": map[string]any{
				condition.Field: map[string]any{
					"lt": condition.Value,
				},
			},
		}
	},
	contract.FilterOperatorLowerThanOrEqual: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"range": map[string]any{
				condition.Field: map[string]any{
					"lte": condition.Value,
				},
			},
		}
	},
	contract.FilterOperatorLowerThanOrEqualOrNil: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"range": map[string]any{
				condition.Field: map[string]any{
					"lte": condition.Value,
				},
			},
		}
	},
	contract.FilterOperatorBegins: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"prefix": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorContains: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"match": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorNotContains: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"match": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorEnds: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"suffix": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorIsNil: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"missing": map[string]any{
				"field": condition.Field,
			},
		}
	},
	contract.FilterOperatorIsNotNil: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"exists": map[string]any{
				"field": condition.Field,
			},
		}
	},
	contract.FilterOperatorIsEmpty: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"missing": map[string]any{
				"field": condition.Field,
			},
		}
	},
	contract.FilterOperatorIsNotEmpty: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"exists": map[string]any{
				"field": condition.Field,
			},
		}
	},
	contract.FilterOperatorIn: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"terms": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
	contract.FilterOperatorNotIn: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"terms": map[string]any{
				condition.Field: condition.Value,
			},
		}
	},
}

func transformCondition(condition contract.FilterCondition, positiveConditions *[]map[string]any, negativeConditions *[]map[string]any) {
	if condition.Field == "" || condition.Operator == "" {
		return
	}
	outputCondition := conditionResolvers[condition.Operator](condition)
	if condition.IsNegative() {
		*negativeConditions = append(*negativeConditions, outputCondition)
		return
	}
	*positiveConditions = append(*positiveConditions, outputCondition)
}

func transformConditions(conditions contract.FilterConditions) ([]map[string]any, []map[string]any) {
	if conditions.IsEmtpy() {
		return nil, nil
	}
	var positiveConditions []map[string]any
	var negativeConditions []map[string]any
	if conditions.Filters != nil {
		for i, filter := range conditions.Filters {
			transformFilters(filter, &positiveConditions[i])
		}
	}
	if conditions.Conditions != nil {
		for _, condition := range conditions.Conditions {
			transformCondition(condition, &positiveConditions, &negativeConditions)
		}
	}
	return positiveConditions, negativeConditions
}

func transformFilters(filters contract.Filters, target *map[string]any) {
	if filters.IsEmpty() {
		return
	}
	logic := "must"
	if filters.Logic != "" {
		logic = "should"
	}
	positiveConditions, negativeConditions := transformConditions(filters.Conditions)
	var outputFilters map[string]any
	if positiveConditions != nil {
		outputFilters[logic] = positiveConditions
	}
	if negativeConditions != nil {
		outputFilters[fmt.Sprintf("%s_not", logic)] = negativeConditions
	}
	if len(outputFilters) > 0 {
		(*target)["bool"] = outputFilters
	}
}

type ElasticOutputTransformer struct {
}

func (t *ElasticOutputTransformer) Transform(input contract.Filters) (*ElasticOutput, error) {
	var transformedData map[string]any
	transformFilters(input, &transformedData)

	var output ElasticOutput
	err := output.SetData(transformedData)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
