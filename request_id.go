package micro

import (
	"context"

	"github.com/koverto/uuid"
)

type ridContextKey struct{}

func ContextWithRequestID(ctx context.Context, rid *uuid.UUID) context.Context {
	return context.WithValue(ctx, ridContextKey{}, rid)
}

func RequestIDFromContext(ctx context.Context) *uuid.UUID {
	return ctx.Value(ridContextKey{}).(*uuid.UUID)
}
