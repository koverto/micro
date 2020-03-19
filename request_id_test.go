package micro

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	mock_client "github.com/koverto/micro/v2/mocks/client"
	mock_server "github.com/koverto/micro/v2/mocks/server"

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		req server.Request
		rsp interface{}
	}

	rid := uuid.New()
	req := mock_server.NewMockRequest(ctrl)

	tests := []struct {
		name    string
		wantErr bool
		fn      server.HandlerFunc
		args    args
	}{
		{
			"With a request ID in the metadata", false,
			func(ctx context.Context, req server.Request, rsp interface{}) error {
				if id, ok := RequestIDFromContext(ctx); !ok {
					return fmt.Errorf("no request ID in context")
				} else if id.Uuid.String() != rid.Uuid.String() {
					return fmt.Errorf("request ID = %v, want %v", id, rid)
				}
				return nil
			},
			args{metadata.Set(context.Background(), requestIDMetadataKey, rid.Uuid.String()), req, nil},
		},
		{
			"Without a request ID in the metadata", false,
			func(ctx context.Context, req server.Request, rsp interface{}) error {
				if id, _ := RequestIDFromContext(ctx); id != nil {
					return fmt.Errorf("expected no request ID in context, got %v", id)
				}
				return nil
			},
			args{context.Background(), req, nil},
		},
		{
			"With an invalid request ID", true,
			func(_ context.Context, _ server.Request, _ interface{}) error {
				return nil
			},
			args{metadata.Set(context.Background(), requestIDMetadataKey, "invalid UUID"), req, nil},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			err := requestIDHandlerWrapper(test.fn)(test.args.ctx, test.args.req, test.args.rsp)

			if (!test.wantErr && err != nil) || (test.wantErr && err == nil) {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
		})
	}
}

func Test_requestIDClientWrapper_Call(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx  context.Context
		req  client.Request
		rsp  interface{}
		opts []client.CallOption
	}

	rid := uuid.New()
	req := client.NewRequest("svc", "endpoint", nil)

	tests := []struct {
		name    string
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			"With a request ID",
			args{ContextWithRequestID(context.Background(), rid), req, nil, []client.CallOption{}},
			rid,
			false,
		},
		{
			"Without a request ID",
			args{context.Background(), req, nil, []client.CallOption{}},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			m := mock_client.NewMockClient(ctrl)
			w := requestIDClientWrapper(m)

			m.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
				func(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) {
					got, ok := metadata.Get(ctx, requestIDMetadataKey)
					if (!ok && test.want != nil) || (ok && test.want == nil) {
						t.Errorf("Metadata request ID ok %v, want %v", ok, test.want)
					}

					if test.want != nil {
						want := test.want.Uuid.String()
						if got != want {
							t.Errorf("Metadata request ID = %v, want %v", got, want)
						}
					}
				},
			)

			if err := w.Call(test.args.ctx, test.args.req, test.args.rsp, test.args.opts...); (err != nil) != test.wantErr {
				t.Errorf("_requestIDClientWrapper.Call() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
