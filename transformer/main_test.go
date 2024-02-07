package transformer

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"github.com/wernerdweight/filter-transformer-go/transformer/input"
	"github.com/wernerdweight/filter-transformer-go/transformer/output"
	"log"
	"reflect"
	"testing"
)

var testInputJson0, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`), &input.JsonInput{})
var testInputJson1, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"logic": "or", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "neq", "value": "val2"}]}]}`), &input.JsonInput{})
var testInputJson2, _ = contract.NewInputOutputType([]byte(`{"logic": "or", "conditions": [{"field": "key", "operator": "not-null", "value": null}]}`), &input.JsonInput{})
var testInputJson3, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "gte", "value": 123}]}`), &input.JsonInput{})
var testInputJson4, _ = contract.NewInputOutputType([]byte(`{"logic": "or", "conditions": [{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "not-empty", "value": null}]}, {"logic": "and", "conditions": [{"field": "key3", "operator": "contains", "value": "val3"}, {"field": "key4", "operator": "gt", "value": 123}]}]}`), &input.JsonInput{})
var invalidInputJson0, _ = contract.NewInputOutputType([]byte(`{"field": "key", "operator": "eq", "value": "val"}`), &input.JsonInput{})
var invalidInputJson1, _ = contract.NewInputOutputType([]byte(`"JSON string"`), &input.JsonInput{})
var invalidInputJson2, _ = contract.NewInputOutputType([]byte(`not JSON at all`), &input.JsonInput{})

var testOutputElastic0, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"term": map[string]any{"key.lowersortable": "val"}}}}}, &output.ElasticOutput{})
var testOutputElastic1, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"bool": map[string]any{"should": []map[string]any{{"term": map[string]any{"key.lowersortable": "val"}}, {"bool": map[string]any{"must_not": []map[string]any{{"term": map[string]any{"key2.lowersortable": "val2"}}}}}}}}}}}, &output.ElasticOutput{})
var testOutputElastic2, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"should": []map[string]any{{"exists": map[string]any{"field": "key"}}}}}, &output.ElasticOutput{})
var testOutputElastic3, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"range": map[string]any{"key": map[string]any{"gte": 123.0}}}}}}, &output.ElasticOutput{})
var testOutputElastic4, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"should": []map[string]any{{"bool": map[string]any{"must": []map[string]any{{"term": map[string]any{"key.lowersortable": "val"}}, {"exists": map[string]any{"field": "key2"}}}}}, {"bool": map[string]any{"must": []map[string]any{{"wildcard": map[string]any{"key3.lowersortable": "*val3*"}}, {"range": map[string]any{"key4": map[string]any{"gt": 123.0}}}}}}}}}, &output.ElasticOutput{})

var testOutputSQL0, _ = contract.NewInputOutputType(output.SQLTuple{Query: "key = $1", Params: []any{"val"}}, &output.SQLOutput{})
var testOutputSQL1, _ = contract.NewInputOutputType(output.SQLTuple{Query: "(key = $1 OR key2 != $2)", Params: []any{"val", "val2"}}, &output.SQLOutput{})
var testOutputSQL2, _ = contract.NewInputOutputType(output.SQLTuple{Query: "key IS NOT NULL", Params: nil}, &output.SQLOutput{})
var testOutputSQL3, _ = contract.NewInputOutputType(output.SQLTuple{Query: "key >= $1", Params: []any{123.0}}, &output.SQLOutput{})
var testOutputSQL4, _ = contract.NewInputOutputType(output.SQLTuple{Query: "((key = $1 AND key2 != '') OR (key3 LIKE $2 AND key4 > $3))", Params: []any{"val", "%val3%", 123.0}}, &output.SQLOutput{})

func TestBasic(t *testing.T) {
	it := input.JsonInputTransformer{}
	ot := output.ElasticOutputTransformer{}
	ft := NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot)
	jsonInput, _ := contract.NewInputOutputType[[]byte, *input.JsonInput]([]byte("test"), &input.JsonInput{})
	transformedOutput, err := ft.Transform(jsonInput)
	log.Printf("output: %+v, err: %+v", transformedOutput, err)
}

func TestFilterTransformer_TransformJsonToElastic(t1 *testing.T) {
	it := input.JsonInputTransformer{}
	ot := output.ElasticOutputTransformer{}
	ft := NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot)
	type testCase struct {
		name    string
		t       FilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput]
		input   input.JsonInput
		want    *output.ElasticOutput
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "empty",
			t:       *ft,
			input:   input.JsonInput{},
			want:    &output.ElasticOutput{},
			wantErr: false,
		},
		{
			name:    "with data",
			t:       *ft,
			input:   *testInputJson0,
			want:    testOutputElastic0,
			wantErr: false,
		},
		{
			name:    "with nested data",
			t:       *ft,
			input:   *testInputJson1,
			want:    testOutputElastic1,
			wantErr: false,
		},
		{
			name:    "with null value",
			t:       *ft,
			input:   *testInputJson2,
			want:    testOutputElastic2,
			wantErr: false,
		},
		{
			name:    "with number value",
			t:       *ft,
			input:   *testInputJson3,
			want:    testOutputElastic3,
			wantErr: false,
		},
		{
			name:    "with complex data",
			t:       *ft,
			input:   *testInputJson4,
			want:    testOutputElastic4,
			wantErr: false,
		},
		{
			name:    "invalid input - wrong structure",
			t:       *ft,
			input:   *invalidInputJson0,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid input - JSON string",
			t:       *ft,
			input:   *invalidInputJson1,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid input - not a JSON",
			t:       *ft,
			input:   *invalidInputJson2,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := tt.t.Transform(&tt.input)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Transform() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterTransformer_TransformJsonToSQL(t1 *testing.T) {
	it := input.JsonInputTransformer{}
	ot := output.SQLOutputTransformer{}
	ft := NewFilterTransformer[[]byte, output.SQLTuple, *input.JsonInput, *output.SQLOutput](&it, &ot)
	type testCase struct {
		name    string
		t       FilterTransformer[[]byte, output.SQLTuple, *input.JsonInput, *output.SQLOutput]
		input   input.JsonInput
		want    *output.SQLOutput
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "empty",
			t:       *ft,
			input:   input.JsonInput{},
			want:    &output.SQLOutput{},
			wantErr: false,
		},
		{
			name:    "with data",
			t:       *ft,
			input:   *testInputJson0,
			want:    testOutputSQL0,
			wantErr: false,
		},
		{
			name:    "with nested data",
			t:       *ft,
			input:   *testInputJson1,
			want:    testOutputSQL1,
			wantErr: false,
		},
		{
			name:    "with null value",
			t:       *ft,
			input:   *testInputJson2,
			want:    testOutputSQL2,
			wantErr: false,
		},
		{
			name:    "with number value",
			t:       *ft,
			input:   *testInputJson3,
			want:    testOutputSQL3,
			wantErr: false,
		},
		{
			name:    "with complex data",
			t:       *ft,
			input:   *testInputJson4,
			want:    testOutputSQL4,
			wantErr: false,
		},
		{
			name:    "invalid input - wrong structure",
			t:       *ft,
			input:   *invalidInputJson0,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid input - JSON string",
			t:       *ft,
			input:   *invalidInputJson1,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid input - not a JSON",
			t:       *ft,
			input:   *invalidInputJson2,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := tt.t.Transform(&tt.input)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Transform() got = %v, want %v", got, tt.want)
			}
		})
	}
}
