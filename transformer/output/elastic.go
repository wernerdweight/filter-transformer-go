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
	if rawData == nil {
		return nil, nil
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

func getElasticFieldVariantByValueType(condition contract.FilterCondition) string {
	field := fmt.Sprintf("%s.lowersortable", condition.Field)
	value := condition.Value
	slice, isSlice := value.([]interface{})
	if isSlice && len(slice) > 0 {
		value = slice[0]
	}
	_, isInt := value.(int)
	_, isFloat := value.(float64)
	_, isBool := value.(bool)
	if isInt || isFloat || isBool {
		field = condition.Field
	}
	return field
}

var conditionResolversElastic = map[contract.FilterOperator]func(contract.FilterCondition) map[string]any{
	contract.FilterOperatorEqual: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"term": map[string]any{
				getElasticFieldVariantByValueType(condition): condition.Value,
			},
		}
	},
	contract.FilterOperatorNotEqual: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"term": map[string]any{
				getElasticFieldVariantByValueType(condition): condition.Value,
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
			"bool": map[string]any{
				"should": []map[string]any{
					{
						"range": map[string]any{
							condition.Field: map[string]any{
								"gte": condition.Value,
							},
						},
					},
					{
						"bool": map[string]any{
							"must_not": []map[string]any{
								{
									"exists": map[string]any{
										"field": condition.Field,
									},
								},
							},
						},
					},
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
			"bool": map[string]any{
				"should": []map[string]any{
					{
						"range": map[string]any{
							condition.Field: map[string]any{
								"lte": condition.Value,
							},
						},
					},
					{
						"bool": map[string]any{
							"must_not": []map[string]any{
								{
									"exists": map[string]any{
										"field": condition.Field,
									},
								},
							},
						},
					},
				},
			},
		}
	},
	contract.FilterOperatorBegins: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"prefix": map[string]any{
				fmt.Sprintf("%s.lowersortable", condition.Field): condition.Value,
			},
		}
	},
	contract.FilterOperatorContains: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"wildcard": map[string]any{
				fmt.Sprintf("%s.lowersortable", condition.Field): fmt.Sprintf("*%s*", condition.Value),
			},
		}
	},
	contract.FilterOperatorNotContains: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"wildcard": map[string]any{
				fmt.Sprintf("%s.lowersortable", condition.Field): fmt.Sprintf("*%s*", condition.Value),
			},
		}
	},
	contract.FilterOperatorEnds: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"wildcard": map[string]any{
				fmt.Sprintf("%s.lowersortable", condition.Field): fmt.Sprintf("*%s", condition.Value),
			},
		}
	},
	contract.FilterOperatorIsNil: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"exists": map[string]any{
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
			"exists": map[string]any{
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
				getElasticFieldVariantByValueType(condition): condition.Value,
			},
		}
	},
	contract.FilterOperatorNotIn: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"terms": map[string]any{
				getElasticFieldVariantByValueType(condition): condition.Value,
			},
		}
	},
	contract.FilterOperatorMatchPhrase: func(condition contract.FilterCondition) map[string]any {
		return map[string]any{
			"match_phrase": map[string]any{
				condition.Field: fmt.Sprint(condition.Value),
			},
		}
	},
}

func transformConditionElastic(condition contract.FilterCondition, positiveConditions *[]map[string]any, negativeConditions *[]map[string]any) {
	if condition.Field == "" || condition.Operator == "" {
		return
	}
	outputCondition := conditionResolversElastic[condition.Operator](condition)
	if condition.IsNegative() {
		*negativeConditions = append(*negativeConditions, outputCondition)
		return
	}
	*positiveConditions = append(*positiveConditions, outputCondition)
}

func transformConditionsElastic(conditions contract.FilterConditions) ([]map[string]any, []map[string]any) {
	if conditions.IsEmpty() {
		return nil, nil
	}
	var positiveConditions []map[string]any
	var negativeConditions []map[string]any
	if conditions.Filters != nil {
		for _, filter := range conditions.Filters {
			var condition = make(map[string]any)
			transformFiltersElastic(filter, &condition)
			positiveConditions = append(positiveConditions, condition)
		}
	}
	if conditions.Conditions != nil {
		for _, condition := range conditions.Conditions {
			transformConditionElastic(condition, &positiveConditions, &negativeConditions)
		}
	}
	return positiveConditions, negativeConditions
}

func transformFiltersElastic(filters contract.Filters, target *map[string]any) {
	if filters.IsEmpty() {
		return
	}
	var outputFilters = make(map[string]any)
	logic := "must"
	if filters.Logic == contract.FilterLogicOr {
		logic = "should"
		outputFilters["minimum_should_match"] = 1
	}
	positiveConditions, negativeConditions := transformConditionsElastic(filters.Conditions)
	if positiveConditions != nil {
		outputFilters[logic] = positiveConditions
	}
	if negativeConditions != nil {
		if logic == "should" {
			negativeShouldConditions := map[string]any{
				"bool": map[string]any{
					"must_not": negativeConditions,
				},
			}
			outputFilters[logic] = append(outputFilters[logic].([]map[string]any), negativeShouldConditions)

		} else {
			outputFilters["must_not"] = negativeConditions
		}
	}
	if len(outputFilters) > 0 {
		(*target)["bool"] = outputFilters
	}
}

type ElasticOutputTransformer struct {
}

func (t *ElasticOutputTransformer) Transform(input contract.Filters) (*ElasticOutput, *contract.Error) {
	var transformedData = make(map[string]any)
	transformFiltersElastic(input, &transformedData)

	var output ElasticOutput
	if len(transformedData) == 0 {
		return &output, nil
	}

	err := output.SetData(transformedData)
	if err != nil {
		return nil, contract.NewError(contract.NonWriteableOutputData, err.Error())
	}
	return &output, nil
}
