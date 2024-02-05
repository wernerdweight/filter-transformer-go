package input

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type JsonInput struct {
	contract.InputOutputType[[]byte]
}

func (i *JsonInput) GetDataString() (string, error) {
	// TODO: implement
	return "", nil
}

func (i *JsonInput) GetDataJson() ([]byte, error) {
	return i.GetData()
}

type JsonInputTransformer struct {
}

func (t *JsonInputTransformer) Transform(input *JsonInput) (contract.Filters, error) {
	var filters contract.Filters
	// TODO: implement
	filters.Logic = contract.FilterLogicAnd
	filters.Conditions = contract.FilterConditions{
		Conditions: []contract.FilterCondition{
			{
				Field:    "field",
				Operator: contract.FilterOperatorEqual,
				Value:    "value",
			},
		},
	}
	return filters, nil
}
