package micro

import (
	"reflect"
	"testing"

	"github.com/micro/go-micro/v2/config/source"
)

func TestNewService(t *testing.T) {
	type args struct {
		name    string
		conf    interface{}
		sources []source.Source
	}

	tests := []struct {
		name    string
		args    args
		want    *Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			got, err := NewService(test.args.name, test.args.conf, test.args.sources...)
			if (err != nil) != test.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewService() = %v, want %v", got, test.want)
			}
		})
	}
}
