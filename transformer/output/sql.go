package output

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type SQLOutput struct {
}

func (o *SQLOutput) SetData(data interface{}) error {
	// TODO: implement
	return nil
}

type SQLOutputTransformer struct {
}

func (t *SQLOutputTransformer) Transform(filters contract.Filters) (*SQLOutput, error) {
	// TODO: implement
	return nil, nil
}
