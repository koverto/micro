package micro

import (
	"context"

	"github.com/koverto/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
)

const REQUEST_ID_METADATA_KEY = "request_id"

type ridContextKey struct{}

func ContextWithRequestID(ctx context.Context, rid *uuid.UUID) context.Context {
	return context.WithValue(ctx, ridContextKey{}, rid)
}

func RequestIDFromContext(ctx context.Context) *uuid.UUID {
	return ctx.Value(ridContextKey{}).(*uuid.UUID)
}

func requestIDHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if rid, ok := metadata.Get(ctx, REQUEST_ID_METADATA_KEY); ok {
			rid, err := uuid.Parse(rid)
			if err != nil {
				return err
			}

			ctx = ContextWithRequestID(ctx, rid)
		}

		return fn(ctx, req, rsp)
	}
}

type _requestIDClientWrapper struct {
	client.Client
}

func requestIDClientWrapper(c client.Client) client.Client {
	return &_requestIDClientWrapper{c}
}

func (w *_requestIDClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if rid := RequestIDFromContext(ctx); rid != nil {
		ctx = metadata.Set(ctx, REQUEST_ID_METADATA_KEY, rid.Uuid.String())
	}

	return w.Client.Call(ctx, req, rsp, opts...)
}
