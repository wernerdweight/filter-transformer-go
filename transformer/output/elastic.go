package output

import (
	"encoding/json"
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

type ElasticOutputTransformer struct {
}

func (t *ElasticOutputTransformer) Transform(input contract.Filters) (*ElasticOutput, error) {
	var output ElasticOutput
	// TODO: implement
	err := output.SetData(map[string]any{"test": "test"})
	if err != nil {
		return nil, err
	}
	return &output, nil
}
