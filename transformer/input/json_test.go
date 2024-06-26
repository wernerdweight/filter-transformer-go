package input

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"reflect"
	"testing"
)

var testInputJson0, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`), &JsonInput{})
var testInputJson1, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"logic": "or", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "neq", "value": "val2"}]}]}`), &JsonInput{})
var testInputJson2, _ = contract.NewInputOutputType([]byte(`{"logic": "or", "conditions": [{"field": "key", "operator": "not-null", "value": null}]}`), &JsonInput{})
var testInputJson3, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "gte", "value": 123}]}`), &JsonInput{})
var testInputJson4, _ = contract.NewInputOutputType([]byte(`{"logic": "or", "conditions": [{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "not-empty", "value": null}]}, {"logic": "and", "conditions": [{"field": "key3", "operator": "contains", "value": "val3"}, {"field": "key4", "operator": "gt", "value": 123}]}]}`), &JsonInput{})
var testInputJson5, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "release_year", "operator": "gte", "value": 2022}, {"logic": "or", "conditions": [{"field": "duration", "operator": "gte", "value": 120}, {"field": "track_name", "operator": "not-contains", "value": "cloud"}, {"field": "release_year", "operator": "neq", "value": 2022}]}]}`), &JsonInput{})
var testInputJson6, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val"}]}`), &JsonInput{})
var testInputJson7, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val,val2"}]}`), &JsonInput{})
var testInputJson8, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val, val2"}]}`), &JsonInput{})
var testInputJson9, _ = contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": ["val", "val2"]}]}`), &JsonInput{})
var invalidInputJson0, _ = contract.NewInputOutputType([]byte(`{"field": "key", "operator": "eq", "value": "val"}`), &JsonInput{})
var invalidInputJson1, _ = contract.NewInputOutputType([]byte(`"JSON string"`), &JsonInput{})
var invalidInputJson2, _ = contract.NewInputOutputType([]byte(`not JSON at all`), &JsonInput{})

func TestJsonInputTransformer_Transform(t1 *testing.T) {
	type args struct {
		input *JsonInput
	}
	tests := []struct {
		name    string
		args    args
		want    contract.Filters
		wantErr bool
	}{
		{
			name: "empty input",
			args: args{
				input: &JsonInput{},
			},
			want:    contract.Filters{},
			wantErr: false,
		},
		{
			name: "input with data",
			args: args{
				input: testInputJson0,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorEqual,
							Value:    "val",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with nested data",
			args: args{
				input: testInputJson1,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Filters: []contract.Filters{
						{
							Logic: contract.FilterLogicOr,
							Conditions: contract.FilterConditions{
								Conditions: []contract.FilterCondition{
									{
										Field:    "key",
										Operator: contract.FilterOperatorEqual,
										Value:    "val",
									},
									{
										Field:    "key2",
										Operator: contract.FilterOperatorNotEqual,
										Value:    "val2",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with null value",
			args: args{
				input: testInputJson2,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicOr,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorIsNotNil,
							Value:    nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with number value",
			args: args{
				input: testInputJson3,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorGreaterThanOrEqual,
							Value:    123.0,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with complex data",
			args: args{
				input: testInputJson4,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicOr,
				Conditions: contract.FilterConditions{
					Filters: []contract.Filters{
						{
							Logic: contract.FilterLogicAnd,
							Conditions: contract.FilterConditions{
								Conditions: []contract.FilterCondition{
									{
										Field:    "key",
										Operator: contract.FilterOperatorEqual,
										Value:    "val",
									},
									{
										Field:    "key2",
										Operator: contract.FilterOperatorIsNotEmpty,
										Value:    nil,
									},
								},
							},
						},
						{
							Logic: contract.FilterLogicAnd,
							Conditions: contract.FilterConditions{
								Conditions: []contract.FilterCondition{
									{
										Field:    "key3",
										Operator: contract.FilterOperatorContains,
										Value:    "val3",
									},
									{
										Field:    "key4",
										Operator: contract.FilterOperatorGreaterThan,
										Value:    123.0,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with complex nested data",
			args: args{
				input: testInputJson5,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "release_year",
							Operator: contract.FilterOperatorGreaterThanOrEqual,
							Value:    2022.0,
						},
					},
					Filters: []contract.Filters{
						{
							Logic: contract.FilterLogicOr,
							Conditions: contract.FilterConditions{
								Conditions: []contract.FilterCondition{
									{
										Field:    "duration",
										Operator: contract.FilterOperatorGreaterThanOrEqual,
										Value:    120.0,
									},
									{
										Field:    "track_name",
										Operator: contract.FilterOperatorNotContains,
										Value:    "cloud",
									},
									{
										Field:    "release_year",
										Operator: contract.FilterOperatorNotEqual,
										Value:    2022.0,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with in operator - single value",
			args: args{
				input: testInputJson6,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorIn,
							Value:    []string{"val"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with in operator - multiple string values",
			args: args{
				input: testInputJson7,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorIn,
							Value:    []string{"val", "val2"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with in operator - multiple string values with space",
			args: args{
				input: testInputJson8,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorIn,
							Value:    []string{"val", "val2"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "input with in operator - multiple string values in array",
			args: args{
				input: testInputJson9,
			},
			want: contract.Filters{
				Logic: contract.FilterLogicAnd,
				Conditions: contract.FilterConditions{
					Conditions: []contract.FilterCondition{
						{
							Field:    "key",
							Operator: contract.FilterOperatorIn,
							Value:    []any{"val", "val2"},
						},
					},
				},
			},
		},
		{
			name: "invalid input - wrong structure",
			args: args{
				input: invalidInputJson0,
			},
			want:    contract.Filters{},
			wantErr: true,
		},
		{
			name: "invalid input - JSON string",
			args: args{
				input: invalidInputJson1,
			},
			want:    contract.Filters{},
			wantErr: true,
		},
		{
			name: "invalid input - not JSON",
			args: args{
				input: invalidInputJson2,
			},
			want:    contract.Filters{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &JsonInputTransformer{}
			got, err := t.Transform(tt.args.input)
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

func TestJsonInput_GetDataJson(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput JsonInput
		want      []byte
		wantErr   bool
	}{
		{
			name:      "empty input",
			jsonInput: JsonInput{},
			want:      nil,
			wantErr:   false,
		},
		{
			name:      "input with data",
			jsonInput: *testInputJson0,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`),
			wantErr:   false,
		},
		{
			name:      "input with nested data",
			jsonInput: *testInputJson1,
			want:      []byte(`{"logic": "and", "conditions": [{"logic": "or", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "neq", "value": "val2"}]}]}`),
			wantErr:   false,
		},
		{
			name:      "input with null value",
			jsonInput: *testInputJson2,
			want:      []byte(`{"logic": "or", "conditions": [{"field": "key", "operator": "not-null", "value": null}]}`),
			wantErr:   false,
		},
		{
			name:      "input with number value",
			jsonInput: *testInputJson3,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "gte", "value": 123}]}`),
			wantErr:   false,
		},
		{
			name:      "input with complex data",
			jsonInput: *testInputJson4,
			want:      []byte(`{"logic": "or", "conditions": [{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "not-empty", "value": null}]}, {"logic": "and", "conditions": [{"field": "key3", "operator": "contains", "value": "val3"}, {"field": "key4", "operator": "gt", "value": 123}]}]}`),
			wantErr:   false,
		},
		{
			name:      "input with complex nested data",
			jsonInput: *testInputJson5,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "release_year", "operator": "gte", "value": 2022}, {"logic": "or", "conditions": [{"field": "duration", "operator": "gte", "value": 120}, {"field": "track_name", "operator": "not-contains", "value": "cloud"}, {"field": "release_year", "operator": "neq", "value": 2022}]}]}`),
			wantErr:   false,
		},
		{
			name:      "input with in operator - single value",
			jsonInput: *testInputJson6,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val"}]}`),
			wantErr:   false,
		},
		{
			name:      "input with in operator - multiple string values",
			jsonInput: *testInputJson7,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val,val2"}]}`),
			wantErr:   false,
		},
		{
			name:      "input with in operator - multiple string values with space",
			jsonInput: *testInputJson8,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val, val2"}]}`),
			wantErr:   false,
		},
		{
			name:      "input with in operator - multiple string values in array",
			jsonInput: *testInputJson9,
			want:      []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": ["val", "val2"]}]}`),
			wantErr:   false,
		},
		{
			name:      "invalid input - wrong structure",
			jsonInput: *invalidInputJson0,
			want:      []byte(`{"field": "key", "operator": "eq", "value": "val"}`),
			wantErr:   false,
		},
		{
			name:      "invalid input - JSON string",
			jsonInput: *invalidInputJson1,
			want:      []byte(`"JSON string"`),
			wantErr:   false,
		},
		{
			name:      "invalid input - not JSON",
			jsonInput: *invalidInputJson2,
			want:      []byte(`not JSON at all`),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jsonInput.GetDataJson()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonInput_GetDataString(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput JsonInput
		want      string
		wantErr   bool
	}{
		{
			name:      "empty input",
			jsonInput: JsonInput{},
			want:      "",
			wantErr:   false,
		},
		{
			name:      "input with data",
			jsonInput: *testInputJson0,
			want:      `{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`,
			wantErr:   false,
		},
		{
			name:      "input with nested data",
			jsonInput: *testInputJson1,
			want:      `{"logic": "and", "conditions": [{"logic": "or", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "neq", "value": "val2"}]}]}`,
			wantErr:   false,
		},
		{
			name:      "input with null value",
			jsonInput: *testInputJson2,
			want:      `{"logic": "or", "conditions": [{"field": "key", "operator": "not-null", "value": null}]}`,
			wantErr:   false,
		},
		{
			name:      "input with number value",
			jsonInput: *testInputJson3,
			want:      `{"logic": "and", "conditions": [{"field": "key", "operator": "gte", "value": 123}]}`,
			wantErr:   false,
		},
		{
			name:      "input with complex data",
			jsonInput: *testInputJson4,
			want:      `{"logic": "or", "conditions": [{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "not-empty", "value": null}]}, {"logic": "and", "conditions": [{"field": "key3", "operator": "contains", "value": "val3"}, {"field": "key4", "operator": "gt", "value": 123}]}]}`,
			wantErr:   false,
		},
		{
			name:      "input with complex nested data",
			jsonInput: *testInputJson5,
			want:      `{"logic": "and", "conditions": [{"field": "release_year", "operator": "gte", "value": 2022}, {"logic": "or", "conditions": [{"field": "duration", "operator": "gte", "value": 120}, {"field": "track_name", "operator": "not-contains", "value": "cloud"}, {"field": "release_year", "operator": "neq", "value": 2022}]}]}`,
			wantErr:   false,
		},
		{
			name:      "input with in operator - single value",
			jsonInput: *testInputJson6,
			want:      `{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val"}]}`,
			wantErr:   false,
		},
		{
			name:      "input with in operator - multiple string values",
			jsonInput: *testInputJson7,
			want:      `{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val,val2"}]}`,
			wantErr:   false,
		},
		{
			name:      "input with in operator - multiple string values with space",
			jsonInput: *testInputJson8,
			want:      `{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": "val, val2"}]}`,
			wantErr:   false,
		},
		{
			name:      "input with in operator - multiple string values in array",
			jsonInput: *testInputJson9,
			want:      `{"logic": "and", "conditions": [{"field": "key", "operator": "in", "value": ["val", "val2"]}]}`,
			wantErr:   false,
		},
		{
			name:      "invalid input - wrong structure",
			jsonInput: *invalidInputJson0,
			want:      `{"field": "key", "operator": "eq", "value": "val"}`,
			wantErr:   false,
		},
		{
			name:      "invalid input - JSON string",
			jsonInput: *invalidInputJson1,
			want:      `"JSON string"`,
			wantErr:   false,
		},
		{
			name:      "invalid input - not JSON",
			jsonInput: *invalidInputJson2,
			want:      `not JSON at all`,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jsonInput.GetDataString()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDataString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
