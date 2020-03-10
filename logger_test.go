package micro

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/server"
	"github.com/rs/zerolog"
)

func Test_logHandlerWrapper(t *testing.T) {
	type args struct {
		fn server.HandlerFunc
	}

	tests := []struct {
		name string
		args args
		want server.HandlerFunc
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := logHandlerWrapper(test.args.fn); !reflect.DeepEqual(got, test.want) {
				t.Errorf("logHandlerWrapper() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_logClientWrapper(t *testing.T) {
	tests := []struct {
		name string
		args client.Client
		want client.Client
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := logClientWrapper(test.args); !reflect.DeepEqual(got, test.want) {
				t.Errorf("logClientWrapper() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_logClientWrapper_Call(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  client.Request
		rsp  interface{}
		opts []client.CallOption
	}

	tests := []struct {
		name    string
		w       *_logClientWrapper
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if err := test.w.Call(test.args.ctx, test.args.req, test.args.rsp, test.args.opts...); (err != nil) != test.wantErr {
				t.Errorf("_logClientWrapper.Call() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_logWrapper(t *testing.T) {
	type args struct {
		ctx context.Context
		req request
	}

	tests := []struct {
		name  string
		args  args
		want  *zerolog.Event
		want1 time.Time
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			got, got1 := logWrapper(test.args.ctx, test.args.req)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("logWrapper() got = %v, want %v", got, test.want)
			}
			if !reflect.DeepEqual(got1, test.want1) {
				t.Errorf("logWrapper() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}
