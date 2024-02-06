package output

import (
	"github.com/wernerdweight/filter-transformer-go/transformer/contract"
	"reflect"
	"testing"
)

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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
