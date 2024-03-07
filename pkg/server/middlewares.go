package server

import (
	"fmt"
	"log/slog"

	huma "github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type middlewareFunc func(ctx huma.Context, next func(huma.Context))

func newLoggerMiddleware(logger *slog.Logger) middlewareFunc {
	return func(ctx huma.Context, next func(huma.Context)) {
		l := logger
		l = l.WithGroup("req")
		l = l.With("method", ctx.Method())
		l = l.With("url", ctx.URL().Path)

		ctx = setContextValue(ctx, loggerKey, l)

		next(ctx)
	}
}

func newRequestIDMiddleware() middlewareFunc {
	const requestIDHeaderName = "X-Request-ID"

	return func(ctx huma.Context, next func(huma.Context)) {
		l := loggerFromRequest(ctx.Context())

		rid := ctx.Header(requestIDHeaderName)
		hasValidID := rid != ""

		if !hasValidID {
			if uid, err := uuid.NewRandom(); err != nil {
				l.Warn(fmt.Errorf("failed to create UUID: %w", err).Error())
			} else {
				rid = uid.String()
				hasValidID = true
			}
		}

		if hasValidID {
			ctx.SetHeader(requestIDHeaderName, rid)
			ctx = setContextValue(ctx, ridKey, rid)

			l = l.With("id", rid)
			ctx = setContextValue(ctx, loggerKey, l)
		}

		next(ctx)
	}
}

func registerMiddlewares(api huma.API, logger *slog.Logger) {
	api.UseMiddleware(newLoggerMiddleware(logger))

	api.UseMiddleware(newRequestIDMiddleware())
}
