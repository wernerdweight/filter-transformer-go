package transformer

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"github.com/wernerdweight/filter-transformer-go/transformer/input"
	"github.com/wernerdweight/filter-transformer-go/transformer/output"
)

type FilterTransformer[IDT any, ODT any, IT contract.InputOutputInterface[IDT], OT contract.InputOutputInterface[ODT]] struct {
	inputTransformer  contract.InputTransformerInterface[IDT, IT]
	outputTransformer contract.OutputTransformerInterface[ODT, OT]
	validationFunc    *contract.ValidationFunc
}

func (t *FilterTransformer[IDT, ODT, IT, OT]) Transform(input IT, failOnEmpty bool) (o OT, err *contract.Error) {
	filter, err := t.inputTransformer.Transform(input)
	if err != nil {
		return
	}
	validationErrors := filter.Validate(t.validationFunc, failOnEmpty)
	if len(validationErrors) > 0 {
		err = contract.NewError(contract.InvalidFiltersStructure, validationErrors)
		return
	}
	o, err = t.outputTransformer.Transform(filter)
	return
}

func (t *FilterTransformer[IDT, ODT, IT, OT]) WithValidationFunc(validationFunc contract.ValidationFunc) *FilterTransformer[IDT, ODT, IT, OT] {
	t.validationFunc = &validationFunc
	return t
}

func NewFilterTransformer[IDT any, ODT any, IT contract.InputOutputInterface[IDT], OT contract.InputOutputInterface[ODT]](
	inputTransformer contract.InputTransformerInterface[IDT, IT],
	outputTransformer contract.OutputTransformerInterface[ODT, OT],
	validationFunc *contract.ValidationFunc,
) *FilterTransformer[IDT, ODT, IT, OT] {
	return &FilterTransformer[IDT, ODT, IT, OT]{
		inputTransformer:  inputTransformer,
		outputTransformer: outputTransformer,
		validationFunc:    validationFunc,
	}
}

func NewJsonToElasticFilterTransformer() *FilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput] {
	it := input.JsonInputTransformer{}
	ot := output.ElasticOutputTransformer{}
	return NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot, nil)
}

func NewJsonToSQLFilterTransformer() *FilterTransformer[[]byte, output.SQLTuple, *input.JsonInput, *output.SQLOutput] {
	it := input.JsonInputTransformer{}
	ot := output.SQLOutputTransformer{}
	return NewFilterTransformer[[]byte, output.SQLTuple, *input.JsonInput, *output.SQLOutput](&it, &ot, nil)
}

// TODO: NewFormDataToElasticFilterTransformer
// TODO: NewFormDataToSQLFilterTransformer
