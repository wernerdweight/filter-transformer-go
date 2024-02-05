package output

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type ElasticOutput struct {
	contract.InputOutputType[map[string]any]
}

func (o *ElasticOutput) GetDataJson() ([]byte, error) {
	// TODO: implement
	return nil, nil
}

func (o *ElasticOutput) GetDataString() (string, error) {
	// TODO: implement
	return "", nil
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
