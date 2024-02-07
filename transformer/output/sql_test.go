package output

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"reflect"
	"testing"
)

var testOutputSQL0, _ = contract.NewInputOutputType(SQLTuple{
	Query:  "key = $1",
	Params: []any{"val"},
}, &SQLOutput{})
var testOutputSQL1, _ = contract.NewInputOutputType(SQLTuple{
	Query:  "(key = $1 OR key2 != $2)",
	Params: []any{"val", "val2"},
}, &SQLOutput{})
var testOutputSQL2, _ = contract.NewInputOutputType(SQLTuple{
	Query:  "key IS NOT NULL",
	Params: nil,
}, &SQLOutput{})
var testOutputSQL3, _ = contract.NewInputOutputType(SQLTuple{
	Query:  "key >= $1",
	Params: []any{123},
}, &SQLOutput{})
var testOutputSQL4, _ = contract.NewInputOutputType(SQLTuple{
	Query:  "((key = $1 AND key2 != '') OR (key3 LIKE $2 AND key4 > $3))",
	Params: []any{"val", "%val3%", 123},
}, &SQLOutput{})

func TestSQLOutputTransformer_Transform(t1 *testing.T) {
	type args struct {
		input contract.Filters
	}
	tests := []struct {
		name    string
		args    args
		want    *SQLOutput
		wantErr bool
	}{
		{
			name: "empty filters",
			args: args{
				input: contract.Filters{},
			},
			want:    &SQLOutput{},
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
			want:    testOutputSQL0,
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
			want:    testOutputSQL1,
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
			want:    testOutputSQL2,
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
			want:    testOutputSQL3,
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
			}, testOutputSQL4,
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &SQLOutputTransformer{}
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

func TestSQLOutput_GetDataJson(t *testing.T) {
	tests := []struct {
		name          string
		elasticOutput SQLOutput
		want          []byte
		wantErr       bool
	}{
		{
			name:          "empty",
			elasticOutput: SQLOutput{},
			want:          nil,
			wantErr:       false,
		},
		{
			name:          "with data",
			elasticOutput: *testOutputSQL0,
			want:          []byte(`{"query":"key = $1","params":["val"]}`),
			wantErr:       false,
		},
		{
			name:          "with nested data",
			elasticOutput: *testOutputSQL1,
			want:          []byte(`{"query":"(key = $1 OR key2 != $2)","params":["val","val2"]}`),
			wantErr:       false,
		},
		{
			name:          "with null value",
			elasticOutput: *testOutputSQL2,
			want:          []byte(`{"query":"key IS NOT NULL","params":null}`),
			wantErr:       false,
		},
		{
			name:          "with number value",
			elasticOutput: *testOutputSQL3,
			want:          []byte(`{"query":"key \u003e= $1","params":[123]}`),
			wantErr:       false,
		},
		{
			name:          "with complex data",
			elasticOutput: *testOutputSQL4,
			want:          []byte(`{"query":"((key = $1 AND key2 != '') OR (key3 LIKE $2 AND key4 \u003e $3))","params":["val","%val3%",123]}`),
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

func TestSQLOutput_GetDataString(t *testing.T) {
	tests := []struct {
		name          string
		elasticOutput SQLOutput
		want          string
		wantErr       bool
	}{
		{
			name:          "empty",
			elasticOutput: SQLOutput{},
			want:          "",
			wantErr:       false,
		},
		{
			name:          "with data",
			elasticOutput: *testOutputSQL0,
			want:          `key = 'val'`,
			wantErr:       false,
		},
		{
			name:          "with nested data",
			elasticOutput: *testOutputSQL1,
			want:          `(key = 'val' OR key2 != 'val2')`,
			wantErr:       false,
		},
		{
			name:          "with null value",
			elasticOutput: *testOutputSQL2,
			want:          `key IS NOT NULL`,
			wantErr:       false,
		},
		{
			name:          "with number value",
			elasticOutput: *testOutputSQL3,
			want:          `key >= '123'`,
			wantErr:       false,
		},
		{
			name:          "with complex data",
			elasticOutput: *testOutputSQL4,
			want:          `((key = 'val' AND key2 != '') OR (key3 LIKE '%val3%' AND key4 > '123'))`,
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

func Test_transformConditionSQL(t *testing.T) {
	type args struct {
		condition        contract.FilterCondition
		outputConditions *[]string
		params           *[]any
	}
	tests := []struct {
		name           string
		args           args
		wantConditions *[]string
		wantParams     *[]any
	}{
		{
			name: "equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorEqual,
					Value:    "val",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key = $1",
			},
			wantParams: &[]any{
				"val",
			},
		},
		{
			name: "equal number",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorEqual,
					Value:    123,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key = $1",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "not equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotEqual,
					Value:    "val",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key != $1",
			},
			wantParams: &[]any{
				"val",
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
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key != $1",
			},
			wantParams: &[]any{
				123,
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
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key > $1",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "greater than date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThan,
					Value:    "2021-01-01",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key > $1",
			},
			wantParams: &[]any{
				"2021-01-01",
			},
		},
		{
			name: "greater than or equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqual,
					Value:    123,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key >= $1",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "greater than or equal date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqual,
					Value:    "2021-01-01",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key >= $1",
			},
			wantParams: &[]any{
				"2021-01-01",
			},
		},
		{
			name: "greater than or equal or nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqualOrNil,
					Value:    123,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"(key >= $1 OR key IS NULL)",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "greater than or equal or nil date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorGreaterThanOrEqualOrNil,
					Value:    "2021-01-01",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"(key >= $1 OR key IS NULL)",
			},
			wantParams: &[]any{
				"2021-01-01",
			},
		},
		{
			name: "lower than",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThan,
					Value:    123,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key < $1",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "lower than date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThan,
					Value:    "2021-01-01",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key < $1",
			},
			wantParams: &[]any{
				"2021-01-01",
			},
		},
		{
			name: "lower than or equal",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqual,
					Value:    123,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key <= $1",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "lower than or equal date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqual,
					Value:    "2021-01-01",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key <= $1",
			},
			wantParams: &[]any{
				"2021-01-01",
			},
		},
		{
			name: "lower than or equal or nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqualOrNil,
					Value:    123,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"(key <= $1 OR key IS NULL)",
			},
			wantParams: &[]any{
				123,
			},
		},
		{
			name: "lower than or equal or nil date",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorLowerThanOrEqualOrNil,
					Value:    "2021-01-01",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"(key <= $1 OR key IS NULL)",
			},
			wantParams: &[]any{
				"2021-01-01",
			},
		},
		{
			name: "begins",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorBegins,
					Value:    "val",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key LIKE $1",
			},
			wantParams: &[]any{
				"val%",
			},
		},
		{
			name: "contains",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorContains,
					Value:    "val",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key LIKE $1",
			},
			wantParams: &[]any{
				"%val%",
			},
		},
		{
			name: "not contains",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotContains,
					Value:    "val",
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key NOT LIKE $1",
			},
			wantParams: &[]any{
				"%val%",
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
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key LIKE $1",
			},
			wantParams: &[]any{
				"%val",
			},
		},
		{
			name: "is nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsNil,
					Value:    nil,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key IS NULL",
			},
			wantParams: &[]any{},
		},
		{
			name: "is not nil",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsNotNil,
					Value:    nil,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key IS NOT NULL",
			},
			wantParams: &[]any{},
		},
		{
			name: "is empty",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsEmpty,
					Value:    nil,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key = ''",
			},
			wantParams: &[]any{},
		},
		{
			name: "is not empty",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIsNotEmpty,
					Value:    nil,
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key != ''",
			},
			wantParams: &[]any{},
		},
		{
			name: "in",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIn,
					Value:    []interface{}{"val1", "val2"},
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key IN ($1, $2)",
			},
			wantParams: &[]any{
				"val1", "val2",
			},
		},
		{
			name: "in numbers",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorIn,
					Value:    []interface{}{123, 456},
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key IN ($1, $2)",
			},
			wantParams: &[]any{
				123, 456,
			},
		},
		{
			name: "not in",
			args: args{
				condition: contract.FilterCondition{
					Field:    "key",
					Operator: contract.FilterOperatorNotIn,
					Value:    []interface{}{"val1", "val2"},
				},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key NOT IN ($1, $2)",
			},
			wantParams: &[]any{
				"val1", "val2",
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
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{
				"key NOT IN ($1, $2)",
			},
			wantParams: &[]any{
				123, 456,
			},
		},
		{
			name: "empty condition",
			args: args{
				condition:        contract.FilterCondition{},
				outputConditions: &[]string{},
				params:           &[]any{},
			},
			wantConditions: &[]string{},
			wantParams:     &[]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformConditionSQL(tt.args.condition, tt.args.outputConditions, tt.args.params)
			if !reflect.DeepEqual(tt.args.outputConditions, tt.wantConditions) {
				t.Errorf("transformConditionSQL() conditions: got = %v, want %v", tt.args.outputConditions, tt.wantConditions)
			}
			if !reflect.DeepEqual(tt.args.params, tt.wantParams) {
				t.Errorf("transformConditionSQL() params: got = %v, want %v", tt.args.params, tt.wantParams)
			}
		})
	}
}
