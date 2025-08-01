package input

import (
	"encoding/json"

	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
)

type JsonInput struct {
	contract.InputOutputType[[]byte]
}

func (i *JsonInput) GetDataString() (string, error) {
	rawData, err := i.GetData()
	if err != nil {
		return "", err
	}
	return string(rawData), nil
}

func (i *JsonInput) GetDataJson() ([]byte, error) {
	return i.GetData()
}

type JsonInputTransformer struct {
}

func (t *JsonInputTransformer) Transform(input *JsonInput) (contract.Filters, *contract.Error) {
	var filters contract.Filters
	rawData, err := input.GetData()
	if rawData == nil {
		return filters, nil
	}
	if err != nil {
		return filters, contract.NewError(contract.UnreadableInputData, err.Error())
	}
	err = json.Unmarshal(rawData, &filters)
	if err != nil {
		return filters, contract.NewError(contract.InvalidInputDataStructure, err.Error())
	}

	return filters, nil
}
