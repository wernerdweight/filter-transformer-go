package transformer

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"github.com/wernerdweight/filter-transformer-go/transformer/input"
	"github.com/wernerdweight/filter-transformer-go/transformer/output"
	"log"
	"testing"
)

func TestBasic(t *testing.T) {
	it := input.JsonInputTransformer{}
	ot := output.ElasticOutputTransformer{}
	ft := NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot)
	jsonInput, _ := contract.NewInputOutputType[[]byte, *input.JsonInput]([]byte("test"), &input.JsonInput{})
	transformedOutput, err := ft.Transform(jsonInput)
	log.Printf("output: %+v, err: %+v", transformedOutput, err)
}
