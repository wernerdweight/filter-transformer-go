package input

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type FormDataInput struct {
	contract.InputOutputType[interface{}] // TODO: multipart form data
}

func (i *FormDataInput) GetDataString() (string, error) {
	// TODO: implement
	return "", nil
}

func (i *FormDataInput) GetDataJson() ([]byte, error) {
	// TODO: implement
	return nil, nil
}

type FormDataInputTransformer struct {
}

func (t *FormDataInputTransformer) Transform(input FormDataInput) (contract.Filters, *contract.Error) {
	// TODO: implement
	return contract.Filters{}, nil
}
