package transformer

import (
	"fmt"
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
var testInputJson5, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "release_year", "operator": "gte", "value": 2022}, {"logic": "or", "conditions": [{"field": "duration", "operator": "gte", "value": 120}, {"field": "track_name", "operator": "not-contains", "value": "cloud"}, {"field": "release_year", "operator": "neq", "value": 2022}]}]}`), &input.JsonInput{})
var testInputJson6, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val"}]}`), &input.JsonInput{})
var testInputJson7, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val,val2"}]}`), &input.JsonInput{})
var testInputJson8, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val, val2"}]}`), &input.JsonInput{})
var testInputJson9, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": ["val", "val2"]}]}`), &input.JsonInput{})
var invalidInputJson0, _ = contract.NewInputOutputType([]byte(`{"field": "key", "operator": "eq", "value": "val"}`), &input.JsonInput{})
var invalidInputJson1, _ = contract.NewInputOutputType([]byte(`"JSON string"`), &input.JsonInput{})
var invalidInputJson2, _ = contract.NewInputOutputType([]byte(`not JSON at all`), &input.JsonInput{})
var invalidInputJson3, _ = contract.NewInputOutputType([]byte(`{"logic": "test", "conditions": [{"field": "test", "operator": "ss", "value": "test"}]}`), &input.JsonInput{})
var invalidInputJson4, _ = contract.NewInputOutputType([]byte(`{"logic": "and","conditions": [{"field": "test", "operator": "ss", "value": "test"}]}`), &input.JsonInput{})
var invalidInputJson5, _ = contract.NewInputOutputType([]byte(`{"logic": "or","conditions": [{"field": "test", "oooooperator": "eq", "value": "test"}]}`), &input.JsonInput{})

var customInvalidInputJson0, _ = contract.NewInputOutputType([]byte(`{"conditions": [{"field": "test", "operator": "eq", "value": 1}]}`), &input.JsonInput{})
var customInvalidInputJson1, _ = contract.NewInputOutputType([]byte(`{"conditions": [{"field": "key", "operator": "neq", "value": 1}]}`), &input.JsonInput{})
var customInvalidInputJson2, _ = contract.NewInputOutputType([]byte(`{"conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`), &input.JsonInput{})
var customInvalidInputJson3, _ = contract.NewInputOutputType([]byte(`{"conditions": [{"field": "key", "operator": "eq", "value": -1}]}`), &input.JsonInput{})
var customInvalidInputJson4, _ = contract.NewInputOutputType([]byte(`{"conditions": [{"field": "key", "operator": "eq", "value": 1}]}`), &input.JsonInput{})

var testOutputElastic0, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"term": map[string]any{"key.lowersortable": "val"}}}}}, &output.ElasticOutput{})
var testOutputElastic1, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"bool": map[string]any{"should": []map[string]any{{"term": map[string]any{"key.lowersortable": "val"}}, {"bool": map[string]any{"must_not": []map[string]any{{"term": map[string]any{"key2.lowersortable": "val2"}}}}}}, "minimum_should_match": 1}}}}}, &output.ElasticOutput{})
var testOutputElastic2, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"should": []map[string]any{{"exists": map[string]any{"field": "key"}}}, "minimum_should_match": 1}}, &output.ElasticOutput{})
var testOutputElastic3, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"range": map[string]any{"key": map[string]any{"gte": 123.0}}}}}}, &output.ElasticOutput{})
var testOutputElastic4, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"should": []map[string]any{{"bool": map[string]any{"must": []map[string]any{{"term": map[string]any{"key.lowersortable": "val"}}, {"exists": map[string]any{"field": "key2"}}}}}, {"bool": map[string]any{"must": []map[string]any{{"wildcard": map[string]any{"key3.lowersortable": "*val3*"}}, {"range": map[string]any{"key4": map[string]any{"gt": 123.0}}}}}}}, "minimum_should_match": 1}}, &output.ElasticOutput{})
var testOutputElastic5, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"bool": map[string]any{"should": []map[string]any{{"range": map[string]any{"duration": map[string]any{"gte": 120.0}}}, {"bool": map[string]any{"must_not": []map[string]any{{"wildcard": map[string]any{"track_name.lowersortable": "*cloud*"}}, {"term": map[string]any{"release_year": 2022.0}}}}}}, "minimum_should_match": 1}}, {"range": map[string]any{"release_year": map[string]any{"gte": 2022.0}}}}}}, &output.ElasticOutput{})
var testOutputElastic6, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"terms": map[string]any{"key.lowersortable": []string{"val"}}}}}}, &output.ElasticOutput{})
var testOutputElastic7, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"terms": map[string]any{"key.lowersortable": []string{"val", "val2"}}}}}}, &output.ElasticOutput{})
var testOutputElastic8, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"terms": map[string]any{"key.lowersortable": []string{"val", "val2"}}}}}}, &output.ElasticOutput{})
var testOutputElastic9, _ = contract.NewInputOutputType(map[string]any{"bool": map[string]any{"must": []map[string]any{{"terms": map[string]any{"key.lowersortable": []any{"val", "val2"}}}}}}, &output.ElasticOutput{})

var testOutputSQL0, _ = contract.NewInputOutputType(output.SQLTuple{Query: "key = $1", Params: []any{"val"}}, &output.SQLOutput{})
var testOutputSQL1, _ = contract.NewInputOutputType(output.SQLTuple{Query: "(key = $1 OR key2 != $2)", Params: []any{"val", "val2"}}, &output.SQLOutput{})
var testOutputSQL2, _ = contract.NewInputOutputType(output.SQLTuple{Query: "key IS NOT NULL", Params: nil}, &output.SQLOutput{})
var testOutputSQL3, _ = contract.NewInputOutputType(output.SQLTuple{Query: "key >= $1", Params: []any{123.0}}, &output.SQLOutput{})
var testOutputSQL4, _ = contract.NewInputOutputType(output.SQLTuple{Query: "((key = $1 AND key2 != '') OR (key3 LIKE $2 AND key4 > $3))", Params: []any{"val", "%val3%", 123.0}}, &output.SQLOutput{})

func TestBasic(t *testing.T) {
	it := input.JsonInputTransformer{}
	ot := output.ElasticOutputTransformer{}
	ft := NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot, nil)
	jsonInput, _ := contract.NewInputOutputType[[]byte, *input.JsonInput]([]byte("test"), &input.JsonInput{})
	transformedOutput, err := ft.Transform(jsonInput)
	log.Printf("output: %+v, err: %+v", transformedOutput, err)
}

func TestFilterTransformer_TransformJsonToElastic(t1 *testing.T) {
	ft := NewJsonToElasticFilterTransformer()
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
			want:    nil,
			wantErr: true,
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
			name:    "with complex nested data ",
			t:       *ft,
			input:   *testInputJson5,
			want:    testOutputElastic5,
			wantErr: false,
		},
		{
			name:  "with in operator",
			t:     *ft,
			input: *testInputJson6,
			want:  testOutputElastic6,
		},
		{
			name:  "with in operator - multiple values",
			t:     *ft,
			input: *testInputJson7,
			want:  testOutputElastic7,
		},
		{
			name:  "with in operator - multiple values with space",
			t:     *ft,
			input: *testInputJson8,
			want:  testOutputElastic8,
		},
		{
			name:  "with in operator - array",
			t:     *ft,
			input: *testInputJson9,
			want:  testOutputElastic9,
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
	ft := NewJsonToSQLFilterTransformer()
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
			want:    nil,
			wantErr: true,
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

func TestFilterTransformer_Validate(t *testing.T) {
	ft := NewJsonToElasticFilterTransformer()
	type testCase struct {
		name  string
		t     FilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput]
		input input.JsonInput
		want  *[]contract.ValidationError
	}
	tests := []testCase{
		{
			name:  "empty",
			t:     *ft,
			input: input.JsonInput{},
			want: &[]contract.ValidationError{
				{
					Path:  "root",
					Error: contract.ValidationErrorEmpty,
				},
			},
		},
		{
			name:  "invalid input - missing logic",
			t:     *ft,
			input: *invalidInputJson3,
			want: &[]contract.ValidationError{
				{
					Path:    "root.logic",
					Error:   contract.ValidationErrorInvalidOperator,
					Field:   "logic",
					Payload: "test",
				},
				{
					Path:    "root.conditions.0.operator",
					Error:   contract.ValidationErrorInvalidOperator,
					Field:   "operator",
					Payload: "ss",
				},
			},
		},
		{
			name:  "invalid input - invalid operator",
			t:     *ft,
			input: *invalidInputJson4,
			want: &[]contract.ValidationError{
				{
					Path:    "root.conditions.0.operator",
					Error:   contract.ValidationErrorInvalidOperator,
					Field:   "operator",
					Payload: "ss",
				},
			},
		},
		{
			name:  "invalid input - invalid operator key",
			t:     *ft,
			input: *invalidInputJson5,
			want: &[]contract.ValidationError{
				{
					Path:  "root.conditions.0.operator",
					Error: contract.ValidationErrorEmpty,
					Field: "operator",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := ft.Transform(&tt.input)
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("Transform().Error got = %v, want %v", got, tt.want)
				return
			}
			if got != nil && !reflect.DeepEqual(got.Payload, *tt.want) {
				t.Errorf("Transform().Error got = %v, want %v", got.Payload, *tt.want)
			}
		})
	}
}

func TestFilterTransformer_CustomValidation(t *testing.T) {
	ft := NewJsonToElasticFilterTransformer().WithValidationFunc(
		func(filterCondition contract.FilterCondition, path string, validationErrors *[]contract.ValidationError) {
			// example: narrow the supported fields
			if filterCondition.Field != "key" {
				*validationErrors = append(*validationErrors, contract.ValidationError{
					Path:    fmt.Sprintf("%s.field", path),
					Error:   "unsupported field",
					Field:   "field",
					Payload: filterCondition.Field,
				})
			}
			// example: only allow certain operators with certain fields
			if filterCondition.Field == "key" && filterCondition.Operator != "eq" {
				*validationErrors = append(*validationErrors, contract.ValidationError{
					Path:  fmt.Sprintf("%s.operator", path),
					Error: "unsupported field operator",
					Field: "operator",
					Payload: map[string]string{
						"operator": string(filterCondition.Operator),
						"field":    filterCondition.Field,
					},
				})
			}
			// example: only allow certain value types with certain fields
			if filterCondition.Field == "key" {
				value, ok := filterCondition.Value.(float64)
				if !ok {
					*validationErrors = append(*validationErrors, contract.ValidationError{
						Path:  fmt.Sprintf("%s.value", path),
						Error: "unsupported field value type",
						Field: "value",
						Payload: map[string]string{
							"value":    fmt.Sprintf("%v", filterCondition.Value),
							"field":    filterCondition.Field,
							"requires": "float64",
						},
					})
					return
				}
				if value < 0 {
					*validationErrors = append(*validationErrors, contract.ValidationError{
						Path:  fmt.Sprintf("%s.value", path),
						Error: "unsupported field value",
						Field: "value",
						Payload: map[string]string{
							"value":  fmt.Sprintf("%v", filterCondition.Value),
							"field":  filterCondition.Field,
							"reason": "negative",
						},
					})
				}
			}
		},
	)
	type testCase struct {
		name  string
		t     FilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput]
		input input.JsonInput
		want  *[]contract.ValidationError
	}
	tests := []testCase{
		{
			name:  "invalid input - unsupported field",
			t:     *ft,
			input: *customInvalidInputJson0,
			want: &[]contract.ValidationError{
				{
					Path:    "root.conditions.0.field",
					Error:   "unsupported field",
					Field:   "field",
					Payload: "test",
				},
			},
		},
		{
			name:  "invalid input - unsupported field operator",
			t:     *ft,
			input: *customInvalidInputJson1,
			want: &[]contract.ValidationError{
				{
					Path:  "root.conditions.0.operator",
					Error: "unsupported field operator",
					Field: "operator",
					Payload: map[string]string{
						"operator": "neq",
						"field":    "key",
					},
				},
			},
		},
		{
			name:  "invalid input - unsupported field value type",
			t:     *ft,
			input: *customInvalidInputJson2,
			want: &[]contract.ValidationError{
				{
					Path:  "root.conditions.0.value",
					Error: "unsupported field value type",
					Field: "value",
					Payload: map[string]string{
						"value":    "val",
						"field":    "key",
						"requires": "float64",
					},
				},
			},
		},
		{
			name:  "invalid input - unsupported field value",
			t:     *ft,
			input: *customInvalidInputJson3,
			want: &[]contract.ValidationError{
				{
					Path:  "root.conditions.0.value",
					Error: "unsupported field value",
					Field: "value",
					Payload: map[string]string{
						"value":  "-1",
						"field":  "key",
						"reason": "negative",
					},
				},
			},
		},
		{
			name:  "valid input",
			t:     *ft,
			input: *customInvalidInputJson4,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := ft.Transform(&tt.input)
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("Transform().Error got = %v, want %v", got, tt.want)
				return
			}
			if got != nil && !reflect.DeepEqual(got.Payload, *tt.want) {
				t.Errorf("Transform().Error got = %v, want %v", got.Payload, *tt.want)
			}
		})
	}
}
