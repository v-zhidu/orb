package elastic

import (
	"reflect"
	"testing"

	"github.com/olivere/elastic"
)

func TestNewElasticClient(t *testing.T) {
	type args struct {
		esConfig *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *elastic.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewElasticClient(tt.args.esConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewElasticClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElasticClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
