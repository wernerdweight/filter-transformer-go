package transformer

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"github.com/wernerdweight/filter-transformer-go/transformer/input"
	"github.com/wernerdweight/filter-transformer-go/transformer/output"
)

type FilterTransformer[IDT any, ODT any, IT contract.InputOutputInterface[IDT], OT contract.InputOutputInterface[ODT]] struct {
	inputTransformer  contract.InputTransformerInterface[IDT, IT]
	outputTransformer contract.OutputTransformerInterface[ODT, OT]
}

func (t *FilterTransformer[IDT, ODT, IT, OT]) Transform(input IT) (o OT, err *contract.Error) {
	filter, err := t.inputTransformer.Transform(input)
	if err != nil {
		return
	}
	o, err = t.outputTransformer.Transform(filter)
	return
}

func NewFilterTransformer[IDT any, ODT any, IT contract.InputOutputInterface[IDT], OT contract.InputOutputInterface[ODT]](
	inputTransformer contract.InputTransformerInterface[IDT, IT],
	outputTransformer contract.OutputTransformerInterface[ODT, OT],
) *FilterTransformer[IDT, ODT, IT, OT] {
	return &FilterTransformer[IDT, ODT, IT, OT]{
		inputTransformer:  inputTransformer,
		outputTransformer: outputTransformer,
	}
}

func NewJsonToElasticFilterTransformer() *FilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput] {
	it := input.JsonInputTransformer{}
	ot := output.ElasticOutputTransformer{}
	return NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot)
}

func NewJsonToSQLFilterTransformer() *FilterTransformer[[]byte, output.SQLTuple, *input.JsonInput, *output.SQLOutput] {
	it := input.JsonInputTransformer{}
	ot := output.SQLOutputTransformer{}
	return NewFilterTransformer[[]byte, output.SQLTuple, *input.JsonInput, *output.SQLOutput](&it, &ot)
}
