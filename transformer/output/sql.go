package output

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type SQLOutput struct {
	contract.InputOutputType[string]
}

func (o *SQLOutput) GetDataJson() ([]byte, error) {
	// TODO: implement
	return nil, nil
}

func (o *SQLOutput) GetDataString() (string, error) {
	return o.GetData()
}

type SQLOutputTransformer struct {
}

func (t *SQLOutputTransformer) Transform(filters contract.Filters) (*SQLOutput, error) {
	// TODO: implement
	return nil, nil
}
