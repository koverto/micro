package micro

import (
	"context"
	"reflect"
	"testing"

	"github.com/koverto/micro/v2/mock_client"

	"github.com/golang/mock/gomock"
	"github.com/koverto/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
)

func TestContextWithRequestID(t *testing.T) {
	type args struct {
		ctx context.Context
		rid *uuid.UUID
	}

	rid := uuid.New()
	tests := []struct {
		name string
		args args
		want *uuid.UUID
	}{
		{"With a UUID", args{context.Background(), rid}, rid},
		{"Without a UUID", args{context.Background(), nil}, nil},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ctx := ContextWithRequestID(test.args.ctx, test.args.rid)
			if got := ctx.Value(ridContextKey{}); !reflect.DeepEqual(got, test.want) {
				t.Errorf("ContextWithRequestID() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRequestIDFromContext(t *testing.T) {
	rid := uuid.New()
	tests := []struct {
		name  string
		args  context.Context
		want  *uuid.UUID
		want1 bool
	}{
		{"With a request ID", context.WithValue(context.Background(), ridContextKey{}, rid), rid, true},
		{"Without a request ID", context.Background(), nil, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			got, got1 := RequestIDFromContext(test.args)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("RequestIDFromContext() got = %v, want %v", got, test.want)
			}
			if got1 != test.want1 {
				t.Errorf("RequestIDFromContext() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}

func Test_requestIDHandlerWrapper(t *testing.T) {
	tests := []struct {
		name string
		args server.HandlerFunc
		want server.HandlerFunc
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := requestIDHandlerWrapper(test.args); !reflect.DeepEqual(got, test.want) {
				t.Errorf("requestIDHandlerWrapper() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_requestIDClientWrapper(t *testing.T) {
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
			if got := requestIDClientWrapper(test.args); !reflect.DeepEqual(got, test.want) {
				t.Errorf("requestIDClientWrapper() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_requestIDClientWrapper_Call(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_client.NewMockClient(ctrl)

	w := &_requestIDClientWrapper{m}

	rid := uuid.New()
	req := client.NewRequest("svc", "endpoint", nil)

	type args struct {
		ctx  context.Context
		req  client.Request
		rsp  interface{}
		opts []client.CallOption
	}

	tests := []struct {
		name    string
		w       *_requestIDClientWrapper
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			"With a request ID",
			w,
			args{ContextWithRequestID(context.Background(), rid), req, nil, []client.CallOption{}},
			rid,
			false,
		},
		// {"Without a request ID", w},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			m.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
				func(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) {
					got, ok := metadata.Get(ctx, requestIDMetadataKey)
					if !ok && test.want != nil {
						t.Errorf("Metadata request ID not set")
					}
					if ok && test.want == nil {
						t.Errorf("Metadata request ID unexpectedly set")
					}

					want := test.want.Uuid.String()
					if got != want {
						t.Errorf("Metadata request ID = %v, want %v", got, want)
					}
				},
			)

			if err := test.w.Call(test.args.ctx, test.args.req, test.args.rsp, test.args.opts...); (err != nil) != test.wantErr {
				t.Errorf("_requestIDClientWrapper.Call() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
