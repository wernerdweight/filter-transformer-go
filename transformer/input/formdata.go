package input

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type FormDataInput struct {
}

func (i *FormDataInput) GetData() (interface{}, error) {
	// TODO: implement
	return nil, nil
}

type FormDataInputTransformer struct {
}

func (t *FormDataInputTransformer) Transform(input FormDataInput) (contract.Filters, error) {
	// TODO: implement
	return contract.Filters{}, nil
}
