package output

import (
	"encoding/json"
	"fmt"
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"log"
	"regexp"
	"strings"
)

type SQLTuple struct {
	Query  string `json:"query"`
	Params []any  `json:"params"`
}

type SQLOutput struct {
	contract.InputOutputType[SQLTuple]
}

func (o *SQLOutput) GetDataJson() ([]byte, error) {
	rawData, err := o.GetData()
	if err != nil {
		return nil, err
	}
	if rawData.Query == "" {
		return nil, nil
	}
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (o *SQLOutput) GetDataString() (string, error) {
	log.Printf("A call to GetDataString on SQLOutput detected - this is not safe! Only meant for debugging purposes!")
	rawData, err := o.GetData()
	if err != nil {
		return "", err
	}
	template := regexp.MustCompile(`\$([0-9]+)`).ReplaceAllString(rawData.Query, "'%[$1]v'")
	return fmt.Sprintf(template, rawData.Params...), nil
}

type SQLOutputTransformer struct {
}

func addToParams(params *[]any, value any) int {
	*params = append(*params, value)
	return len(*params)
}

var conditionResolversSQL = map[contract.FilterOperator]func(contract.FilterCondition, *[]any) string{
	contract.FilterOperatorEqual: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("%s = $%d", condition.Field, index)
	},
	contract.FilterOperatorNotEqual: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("%s != $%d", condition.Field, index)
	},
	contract.FilterOperatorGreaterThan: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("%s > $%d", condition.Field, index)
	},
	contract.FilterOperatorGreaterThanOrEqual: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("%s >= $%d", condition.Field, index)
	},
	contract.FilterOperatorGreaterThanOrEqualOrNil: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("(%s >= $%d OR %s IS NULL)", condition.Field, index, condition.Field)
	},
	contract.FilterOperatorLowerThan: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("%s < $%d", condition.Field, index)
	},
	contract.FilterOperatorLowerThanOrEqual: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("%s <= $%d", condition.Field, index)
	},
	contract.FilterOperatorLowerThanOrEqualOrNil: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, condition.Value)
		return fmt.Sprintf("(%s <= $%d OR %s IS NULL)", condition.Field, index, condition.Field)
	},
	contract.FilterOperatorBegins: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, fmt.Sprintf("%s%%", condition.Value))
		return fmt.Sprintf("%s LIKE $%d", condition.Field, index)
	},
	contract.FilterOperatorContains: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, fmt.Sprintf("%%%s%%", condition.Value))
		return fmt.Sprintf("%s LIKE $%d", condition.Field, index)
	},
	contract.FilterOperatorNotContains: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, fmt.Sprintf("%%%s%%", condition.Value))
		return fmt.Sprintf("%s NOT LIKE $%d", condition.Field, index)
	},
	contract.FilterOperatorEnds: func(condition contract.FilterCondition, params *[]any) string {
		index := addToParams(params, fmt.Sprintf("%%%s", condition.Value))
		return fmt.Sprintf("%s LIKE $%d", condition.Field, index)
	},
	contract.FilterOperatorIsNil: func(condition contract.FilterCondition, params *[]any) string {
		return fmt.Sprintf("%s IS NULL", condition.Field)
	},
	contract.FilterOperatorIsNotNil: func(condition contract.FilterCondition, params *[]any) string {
		return fmt.Sprintf("%s IS NOT NULL", condition.Field)
	},
	contract.FilterOperatorIsEmpty: func(condition contract.FilterCondition, params *[]any) string {
		return fmt.Sprintf("%s = ''", condition.Field)
	},
	contract.FilterOperatorIsNotEmpty: func(condition contract.FilterCondition, params *[]any) string {
		return fmt.Sprintf("%s != ''", condition.Field)
	},
	contract.FilterOperatorIn: func(condition contract.FilterCondition, params *[]any) string {
		var indices []string
		for _, value := range condition.Value.([]any) {
			index := addToParams(params, value)
			indices = append(indices, fmt.Sprintf("$%d", index))
		}
		return fmt.Sprintf("%s IN (%s)", condition.Field, strings.Join(indices, ", "))
	},
	contract.FilterOperatorNotIn: func(condition contract.FilterCondition, params *[]any) string {
		var indices []string
		for _, value := range condition.Value.([]any) {
			index := addToParams(params, value)
			indices = append(indices, fmt.Sprintf("$%d", index))
		}
		return fmt.Sprintf("%s NOT IN (%s)", condition.Field, strings.Join(indices, ", "))
	},
}

func transformConditionSQL(condition contract.FilterCondition, outputConditions *[]string, params *[]any) {
	if condition.Field == "" || condition.Operator == "" {
		return
	}
	outputCondition := conditionResolversSQL[condition.Operator](condition, params)
	*outputConditions = append(*outputConditions, outputCondition)
}

func transformConditionsSQL(conditions contract.FilterConditions, params *[]any) []string {
	if conditions.IsEmtpy() {
		return nil
	}
	var outputConditions []string
	if conditions.Filters != nil {
		for _, filter := range conditions.Filters {
			var condition string
			transformFiltersSQL(filter, &condition, params)
			outputConditions = append(outputConditions, condition)
		}
	}
	if conditions.Conditions != nil {
		for _, condition := range conditions.Conditions {
			transformConditionSQL(condition, &outputConditions, params)
		}
	}
	return outputConditions
}

func transformFiltersSQL(filters contract.Filters, target *string, params *[]any) {
	if filters.IsEmpty() {
		return
	}
	conditions := transformConditionsSQL(filters.Conditions, params)
	if len(conditions) == 0 {
		return
	}
	if len(conditions) == 1 {
		*target = conditions[0]
		return
	}
	*target = fmt.Sprintf("(%s)", strings.Join(conditions, fmt.Sprintf(" %s ", strings.ToUpper(string(filters.Logic)))))
}

func (t *SQLOutputTransformer) Transform(input contract.Filters) (*SQLOutput, error) {
	var sql string
	var params []any
	transformFiltersSQL(input, &sql, &params)

	var output SQLOutput
	if sql == "" {
		return &output, nil
	}

	err := output.SetData(SQLTuple{
		Query:  sql,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	return &output, nil
}
