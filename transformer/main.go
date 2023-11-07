package transformer

import "github.com/wernerdweight/filter-transformer-go/transformer/contract"

type FilterTransformer[IT contract.InputType, OT contract.OutputType] struct {
	inputTransformer  contract.InputTransformerInterface[IT]
	outputTransformer contract.OutputTransformerInterface[OT]
}

func (t *FilterTransformer[IT, OT]) Transform(input IT) (*OT, error) {
	filter, err := t.inputTransformer.Transform(input)
	if err != nil {
		return nil, err
	}

	return t.outputTransformer.Transform(filter)
}

func NewFilterTransformer[IT contract.InputType, OT contract.OutputType](
	inputTransformer contract.InputTransformerInterface[IT],
	outputTransformer contract.OutputTransformerInterface[OT],
) *FilterTransformer[IT, OT] {
	return &FilterTransformer[IT, OT]{
		inputTransformer:  inputTransformer,
		outputTransformer: outputTransformer,
	}
}
