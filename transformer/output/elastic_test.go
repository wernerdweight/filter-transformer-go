package output

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"reflect"
	"testing"
)

var testOutput0, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"term": map[string]any{
					"key": "val",
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutput1, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"bool": map[string]any{
					"should": []map[string]any{
						{
							"term": map[string]any{
								"key": "val",
							},
						},
					},
					"should_not": []map[string]any{
						{
							"term": map[string]any{
								"key2": "val2",
							},
						},
					},
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutput2, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"should": []map[string]any{
			{
				"exists": map[string]any{
					"field": "key",
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutput3, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"range": map[string]any{
					"key": map[string]any{
						"gte": 123,
					},
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutput4, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"should": []map[string]any{
			{
				"bool": map[string]any{
					"must": []map[string]any{
						{
							"term": map[string]any{
								"key": "val",
							},
						},
						{
							"exists": map[string]any{
								"field": "key2",
							},
						},
					},
				},
			},
			{
				"bool": map[string]any{
					"must": []map[string]any{
						{
							"match": map[string]any{
								"key3": "val3",
							},
						},
						{
							"range": map[string]any{
								"key4": map[string]any{
									"gt": 123,
								},
							},
						},
					},
				},
			},
		},
	},
}, &ElasticOutput{})

func TestElasticOutputTransformer_Transform(t1 *testing.T) {
	type args struct {
		input contract.Filters
	}
	tests := []struct {
		name    string
		args    args
		want    *ElasticOutput
		wantErr bool
	}{
		{
			name: "empty filters",
			args: args{
				input: contract.Filters{},
			},
			want:    &ElasticOutput{},
			wantErr: false,
		},
		{
			name: "with data",
			args: args{
				input: contract.Filters{
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
			},
			want:    testOutput0,
			wantErr: false,
		},
		{
			name: "with nested data",
			args: args{
				input: contract.Filters{
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
			},
			want:    testOutput1,
			wantErr: false,
		},
		{
			name: "with null value",
			args: args{
				input: contract.Filters{
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
			},
			want:    testOutput2,
			wantErr: false,
		},
		{
			name: "with number value",
			args: args{
				input: contract.Filters{
					Logic: contract.FilterLogicAnd,
					Conditions: contract.FilterConditions{
						Conditions: []contract.FilterCondition{
							{
								Field:    "key",
								Operator: contract.FilterOperatorGreaterThanOrEqual,
								Value:    123,
							},
						},
					},
				},
			},
			want:    testOutput3,
			wantErr: false,
		},
		{
			"with complex data",
			args{
				input: contract.Filters{
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
											Value:    123,
										},
									},
								},
							},
						},
					},
				},
			},
			testOutput4,
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &ElasticOutputTransformer{}
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

func TestElasticOutput_GetDataJson(t *testing.T) {
	tests := []struct {
		name          string
		elasticOutput ElasticOutput
		want          []byte
		wantErr       bool
	}{
		{
			name:          "empty",
			elasticOutput: ElasticOutput{},
			want:          nil,
			wantErr:       false,
		},
		{
			name:          "with data",
			elasticOutput: *testOutput0,
			want:          []byte(`{"bool":{"must":[{"term":{"key":"val"}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with nested data",
			elasticOutput: *testOutput1,
			want:          []byte(`{"bool":{"must":[{"bool":{"should":[{"term":{"key":"val"}}],"should_not":[{"term":{"key2":"val2"}}]}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with null value",
			elasticOutput: *testOutput2,
			want:          []byte(`{"bool":{"should":[{"exists":{"field":"key"}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with number value",
			elasticOutput: *testOutput3,
			want:          []byte(`{"bool":{"must":[{"range":{"key":{"gte":123}}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with complex data",
			elasticOutput: *testOutput4,
			want:          []byte(`{"bool":{"should":[{"bool":{"must":[{"term":{"key":"val"}},{"exists":{"field":"key2"}}]}},{"bool":{"must":[{"match":{"key3":"val3"}},{"range":{"key4":{"gt":123}}}]}}]}}`),
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.elasticOutput.GetDataJson()
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

func TestElasticOutput_GetDataString(t *testing.T) {
	tests := []struct {
		name          string
		elasticOutput ElasticOutput
		want          string
		wantErr       bool
	}{
		{
			name:          "empty",
			elasticOutput: ElasticOutput{},
			want:          "",
			wantErr:       false,
		},
		{
			name:          "with data",
			elasticOutput: *testOutput0,
			want:          `{"bool":{"must":[{"term":{"key":"val"}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with nested data",
			elasticOutput: *testOutput1,
			want:          `{"bool":{"must":[{"bool":{"should":[{"term":{"key":"val"}}],"should_not":[{"term":{"key2":"val2"}}]}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with null value",
			elasticOutput: *testOutput2,
			want:          `{"bool":{"should":[{"exists":{"field":"key"}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with number value",
			elasticOutput: *testOutput3,
			want:          `{"bool":{"must":[{"range":{"key":{"gte":123}}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with complex data",
			elasticOutput: *testOutput4,
			want:          `{"bool":{"should":[{"bool":{"must":[{"term":{"key":"val"}},{"exists":{"field":"key2"}}]}},{"bool":{"must":[{"match":{"key3":"val3"}},{"range":{"key4":{"gt":123}}}]}}]}}`,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.elasticOutput.GetDataString()
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