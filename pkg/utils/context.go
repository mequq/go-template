package utils

import (
	"context"
	"log/slog"
)

// import (
// 	"context"
// 	"net/http"

// 	"log/slog"
// )

// type appContextKey int

// const (
// 	requestID appContextKey = iota
// 	reqClientIP
// )

// type AppContext struct {
// 	// Context
// 	ctx context.Context
// 	// Logger
// }

// // New AppContext
// func LogContext(ctx context.Context) *AppContext {
// 	return &AppContext{
// 		ctx: ctx,
// 	}
// }

// // logValue returns a slog.Value with all the context values
// func (a *AppContext) LogValue() slog.Value {

// 	var attrs []slog.Attr

// 	if a.ctx.Value(reqClientIP) != nil {
// 		attrs = append(attrs, slog.String("clientIP", a.ctx.Value(reqClientIP).(string)))
// 	}

// 	if a.ctx.Value(requestID) != nil {
// 		attrs = append(attrs, slog.String("requestID", a.ctx.Value(requestID).(string)))
// 	}

// 	attrs = append(attrs, slog.Any("err", a.ctx.Err()))
// 	return slog.GroupValue(attrs...)

// }

// // SetContextFromHttpReq set context from http request
// func SetContextFromHttpReq(ctx context.Context, r *http.Request) context.Context {
// 	nCtx := context.WithValue(ctx, requestID, r.Header.Get("x-request-id"))
// 	var requestIP string
// 	if r.Header.Get("x-forwarded-for") != "" {
// 		requestIP = r.Header.Get("x-forwarded-for")
// 	} else {
// 		requestIP = r.RemoteAddr
// 	}
// 	nCtx = context.WithValue(nCtx, reqClientIP, requestIP)
// 	return nCtx
// }

type keyType int

const LoggerContext keyType = 1

func SetLoggerContext(ctx context.Context, attr slog.Attr) context.Context {
	attrs := []slog.Attr{attr}
	if ctxattrs, ok := ctx.Value(LoggerContext).([]slog.Attr); ok {
		attrs = append(attrs, ctxattrs...)
	}
	return context.WithValue(ctx, LoggerContext, attrs)
}

func GetLoggerContext(ctx context.Context) slog.Value {
	if ctx == nil {
		return slog.GroupValue()
	}

	if ctx.Value(LoggerContext) != nil {
		attrs := ctx.Value(LoggerContext).([]slog.Attr)
		return slog.GroupValue(attrs...)
	}
	return slog.GroupValue()
}

func GetLoggerContextAsAttrs(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}

	if ctx.Value(LoggerContext) != nil {
		attrs := ctx.Value(LoggerContext).([]slog.Attr)
		return attrs
	}
	return nil
}

type ContextLoggerHandler struct {
	slog.Handler
}

func NewContextLoggerHandler(handler slog.Handler) slog.Handler {
	return &ContextLoggerHandler{
		Handler: handler,
	}
}

func (c *ContextLoggerHandler) Handle(ctx context.Context, r slog.Record) error {

	attr := slog.GroupValue(GetLoggerContextAsAttrs(ctx)...)
	r.AddAttrs(slog.Attr{
		Key:   "context",
		Value: attr,
	})

	return c.Handler.Handle(ctx, r)
}

func (c *ContextLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextLoggerHandler{
		Handler: c.Handler.WithAttrs(attrs),
	}
}

func (c *ContextLoggerHandler) WithGroup(name string) slog.Handler {
	return &ContextLoggerHandler{
		Handler: c.Handler.WithGroup(name),
	}
}
func (c *ContextLoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.Handler.Enabled(ctx, level)
}
