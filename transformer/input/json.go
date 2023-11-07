package input

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type JsonInput struct {
}

func (i *JsonInput) GetData() (interface{}, error) {
	// TODO: implement
	return nil, nil
}

type JsonInputTransformer struct {
}

func (t *JsonInputTransformer) Transform(input JsonInput) (contract.Filters, error) {
	// TODO: implement
	return contract.Filters{}, nil
}
