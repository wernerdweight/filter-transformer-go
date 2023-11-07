package output

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type ElasticOutput struct {
}

func (o *ElasticOutput) SetData(data interface{}) error {
	// TODO: implement
	return nil
}

type ElasticOutputTransformer struct {
}

func (t *ElasticOutputTransformer) Transform(filters contract.Filters) (*ElasticOutput, error) {
	// TODO: implement
	return nil, nil
}
