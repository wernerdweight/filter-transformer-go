package input

import (
	"encoding/json"
	"errors"
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

func (t *JsonInputTransformer) Transform(input *JsonInput) (contract.Filters, error) {
	var filters contract.Filters
	rawData, err := input.GetData()
	if rawData == nil {
		return filters, nil
	}
	if err != nil {
		// TODO: custom error - unreadable data
		return filters, err
	}
	err = json.Unmarshal(rawData, &filters)
	if err != nil {
		// TODO: custom error - structure invalid
		return filters, err
	}
	if len(rawData) > 0 && filters.IsEmpty() {
		// not empty, but not a valid filter conditions
		return filters, errors.New("invalid filters structure")
	}
	return filters, nil
}
