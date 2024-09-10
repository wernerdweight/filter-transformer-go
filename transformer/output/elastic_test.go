package output

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"reflect"
	"testing"
)

var testOutputElastic0, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"term": map[string]any{
					"key.lowersortable": "val",
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutputElastic1, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"bool": map[string]any{
					"should": []map[string]any{
						{
							"term": map[string]any{
								"key.lowersortable": "val",
							},
						},
						{
							"bool": map[string]any{
								"must_not": []map[string]any{
									{
										"term": map[string]any{
											"key2.lowersortable": "val2",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutputElastic2, _ = contract.NewInputOutputType(map[string]any{
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
var testOutputElastic3, _ = contract.NewInputOutputType(map[string]any{
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
var testOutputElastic4, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"should": []map[string]any{
			{
				"bool": map[string]any{
					"must": []map[string]any{
						{
							"term": map[string]any{
								"key.lowersortable": "val",
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
							"wildcard": map[string]any{
								"key3.lowersortable": "*val3*",
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
var testOutputElastic5, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"bool": map[string]any{
					"should": []map[string]any{
						{
							"range": map[string]any{
								"duration": map[string]any{"gte": 120.0},
							},
						}, {
							"bool": map[string]any{
								"must_not": []map[string]any{
									{
										"wildcard": map[string]any{"track_name.lowersortable": "*cloud*"},
									}, {
										"term": map[string]any{"release_year": 2022.0},
									},
								},
							},
						},
					},
				},
			},
			{
				"range": map[string]any{
					"release_year": map[string]any{"gte": 2022.0},
				},
			},
		},
	},
}, &ElasticOutput{})
var testOutputElastic6, _ = contract.NewInputOutputType(map[string]any{
	"bool": map[string]any{
		"must": []map[string]any{
			{
				"terms": map[string]any{
					"key.lowersortable": []string{"val", "val2"},
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
			want:    testOutputElastic0,
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
			want:    testOutputElastic1,
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
			want:    testOutputElastic2,
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
			want:    testOutputElastic3,
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
			testOutputElastic4,
			false,
		},
		{
			"with complex nested data",
			args{
				input: contract.Filters{
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
			},
			testOutputElastic5,
			false,
		},
		{
			"with in operator and array of values",
			args{
				input: contract.Filters{
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
			},
			testOutputElastic6,
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
			elasticOutput: *testOutputElastic0,
			want:          []byte(`{"bool":{"must":[{"term":{"key.lowersortable":"val"}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with nested data",
			elasticOutput: *testOutputElastic1,
			want:          []byte(`{"bool":{"must":[{"bool":{"should":[{"term":{"key.lowersortable":"val"}},{"bool":{"must_not":[{"term":{"key2.lowersortable":"val2"}}]}}]}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with null value",
			elasticOutput: *testOutputElastic2,
			want:          []byte(`{"bool":{"should":[{"exists":{"field":"key"}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with number value",
			elasticOutput: *testOutputElastic3,
			want:          []byte(`{"bool":{"must":[{"range":{"key":{"gte":123}}}]}}`),
			wantErr:       false,
		},
		{
			name:          "with complex data",
			elasticOutput: *testOutputElastic4,
			want:          []byte(`{"bool":{"should":[{"bool":{"must":[{"term":{"key.lowersortable":"val"}},{"exists":{"field":"key2"}}]}},{"bool":{"must":[{"wildcard":{"key3.lowersortable":"*val3*"}},{"range":{"key4":{"gt":123}}}]}}]}}`),
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
			elasticOutput: *testOutputElastic0,
			want:          `{"bool":{"must":[{"term":{"key.lowersortable":"val"}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with nested data",
			elasticOutput: *testOutputElastic1,
			want:          `{"bool":{"must":[{"bool":{"should":[{"term":{"key.lowersortable":"val"}},{"bool":{"must_not":[{"term":{"key2.lowersortable":"val2"}}]}}]}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with null value",
			elasticOutput: *testOutputElastic2,
			want:          `{"bool":{"should":[{"exists":{"field":"key"}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with number value",
			elasticOutput: *testOutputElastic3,
			want:          `{"bool":{"must":[{"range":{"key":{"gte":123}}}]}}`,
			wantErr:       false,
		},
		{
			name:          "with complex data",
			elasticOutput: *testOutputElastic4,
			want:          `{"bool":{"should":[{"bool":{"must":[{"term":{"key.lowersortable":"val"}},{"exists":{"field":"key2"}}]}},{"bool":{"must":[{"wildcard":{"key3.lowersortable":"*val3*"}},{"range":{"key4":{"gt":123}}}]}}]}}`,
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

func Test_transformConditionElastic(t *testing.T) {
	type args struct {
		condition          contract.FilterCondition
		positiveConditions *[]map[string]any
		negativeConditions *[]map[string]any
	}
	tests := []struct {
		name         string
		args         args
		wantPositive *[]map[string]any
		wantNegative *[]map[string]any
	}{
		{
			name: "equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorEqual,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"term": map[string]any{
						"key.lowersortable": "val",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "equal number",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorEqual,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"term": map[string]any{
						"key": 123,
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "not equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotEqual,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"term": map[string]any{
						"key.lowersortable": "val",
					},
				},
			},
		},
		{
			name: "not equal number",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotEqual,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"term": map[string]any{
						"key": 123,
					},
				},
			},
		},
		{
			name: "greater than",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThan,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"gt": 123,
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "greater than date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThan,
					Value:    "2021-01-01",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"gt": "2021-01-01",
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "greater than or equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqual,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"gte": 123,
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "greater than or equal date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqual,
					Value:    "2021-01-01",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"gte": "2021-01-01",
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "greater than or equal or nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqualOrNil,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"bool": map[string]any{
						"should": []map[string]any{
							{
								"range": map[string]any{
									"key": map[string]any{
										"gte": 123,
									},
								},
							},
							{
								"bool": map[string]any{
									"must_not": []map[string]any{
										{
											"exists": map[string]any{
												"field": "key",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "greater than or equal or nil date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqualOrNil,
					Value:    "2021-01-01",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"bool": map[string]any{
						"should": []map[string]any{
							{
								"range": map[string]any{
									"key": map[string]any{
										"gte": "2021-01-01",
									},
								},
							},
							{
								"bool": map[string]any{
									"must_not": []map[string]any{
										{
											"exists": map[string]any{
												"field": "key",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "lower than",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThan,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"lt": 123,
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "lower than date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThan,
					Value:    "2021-01-01",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"lt": "2021-01-01",
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "lower than or equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqual,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"lte": 123,
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "lower than or equal date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqual,
					Value:    "2021-01-01",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"range": map[string]any{
						"key": map[string]any{
							"lte": "2021-01-01",
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "lower than or equal or nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqualOrNil,
					Value:    123,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"bool": map[string]any{
						"should": []map[string]any{
							{
								"range": map[string]any{
									"key": map[string]any{
										"lte": 123,
									},
								},
							},
							{
								"bool": map[string]any{
									"must_not": []map[string]any{
										{
											"exists": map[string]any{
												"field": "key",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "lower than or equal or nil date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqualOrNil,
					Value:    "2021-01-01",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"bool": map[string]any{
						"should": []map[string]any{
							{
								"range": map[string]any{
									"key": map[string]any{
										"lte": "2021-01-01",
									},
								},
							},
							{
								"bool": map[string]any{
									"must_not": []map[string]any{
										{
											"exists": map[string]any{
												"field": "key",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "begins",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorBegins,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"prefix": map[string]any{
						"key.lowersortable": "val",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "contains",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorContains,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"wildcard": map[string]any{
						"key.lowersortable": "*val*",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "not contains",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotContains,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"wildcard": map[string]any{
						"key.lowersortable": "*val*",
					},
				},
			},
		},
		{
			name: "ends",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorEnds,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"wildcard": map[string]any{
						"key.lowersortable": "*val",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "is nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsNil,
					Value:    nil,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"exists": map[string]any{
						"field": "key",
					},
				},
			},
		},
		{
			name: "is not nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsNotNil,
					Value:    nil,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"exists": map[string]any{
						"field": "key",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "is empty",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsEmpty,
					Value:    nil,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"exists": map[string]any{
						"field": "key",
					},
				},
			},
		},
		{
			name: "is not empty",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsNotEmpty,
					Value:    nil,
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"exists": map[string]any{
						"field": "key",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "in",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIn,
					Value:    []interface{}{"val1", "val2"},
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"terms": map[string]any{
						"key.lowersortable": []interface{}{"val1", "val2"},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "in numbers",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIn,
					Value:    []interface{}{123, 456},
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"terms": map[string]any{
						"key": []interface{}{123, 456},
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "not in",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotIn,
					Value:    []interface{}{"val1", "val2"},
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"terms": map[string]any{
						"key.lowersortable": []interface{}{"val1", "val2"},
					},
				},
			},
		},
		{
			name: "not in numbers",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotIn,
					Value:    []interface{}{123, 456},
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{
				{
					"terms": map[string]any{
						"key": []interface{}{123, 456},
					},
				},
			},
		},
		{
			name: "match-phrase",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorMatchPhrase,
					Value:    "val",
				},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{
				{
					"match_phrase": map[string]any{
						"key": "val",
					},
				},
			},
			wantNegative: &[]map[string]any{},
		},
		{
			name: "empty condition",
			args: args{
				condition:          contract.FilterCondition{},
				positiveConditions: &[]map[string]any{},
				negativeConditions: &[]map[string]any{},
			},
			wantPositive: &[]map[string]any{},
			wantNegative: &[]map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformConditionElastic(tt.args.condition, tt.args.positiveConditions, tt.args.negativeConditions)
			if !reflect.DeepEqual(tt.args.positiveConditions, tt.wantPositive) {
				t.Errorf("transformConditionElastic() positive: got = %v, want %v", tt.args.positiveConditions, tt.wantPositive)
			}
			if !reflect.DeepEqual(tt.args.negativeConditions, tt.wantNegative) {
				t.Errorf("transformConditionElastic() negative: got = %v, want %v", tt.args.negativeConditions, tt.wantNegative)
			}
		})
	}
}
