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
		// TODO: Add test cases.
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
	type fields struct {
		InputOutputType contract.InputOutputType
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ElasticOutput{
				InputOutputType: tt.fields.InputOutputType,
			}
			got, err := o.GetDataJson()
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
	type fields struct {
		InputOutputType contract.InputOutputType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ElasticOutput{
				InputOutputType: tt.fields.InputOutputType,
			}
			got, err := o.GetDataString()
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
