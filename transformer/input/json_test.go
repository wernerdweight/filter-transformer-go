package input

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"reflect"
	"testing"
)

func TestJsonInputTransformer_Transform(t1 *testing.T) {
	testInput0, _ := contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`), &JsonInput{})
	testInput1, _ := contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"logic": "or", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "neq", "value": "val2"}]}]}`), &JsonInput{})
	testInput2, _ := contract.NewInputOutputType([]byte(`{"logic": "or", "conditions": [{"field": "key", "operator": "not-null", "value": null}]}`), &JsonInput{})
	testInput3, _ := contract.NewInputOutputType([]byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "gte", "value": 123}]}`), &JsonInput{})
	testInput4, _ := contract.NewInputOutputType([]byte(`{"logic": "or", "conditions": [{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}, {"field": "key2", "operator": "not-empty", "value": null}]}, {"logic": "and", "conditions": [{"field": "key3", "operator": "contains", "value": "val3"}, {"field": "key4", "operator": "gt", "value": 123}]}]}`), &JsonInput{})
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
				input: testInput0,
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
				input: testInput1,
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
				input: testInput2,
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
				input: testInput3,
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
			name: "input with nested data",
			args: args{
				input: testInput4,
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
