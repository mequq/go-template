package biz

import (
	"log/slog"
	"runtime"

	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(
	NewHealthzUseCase,
)

func TraceInfo() slog.Value {
	if _, file, number, ok := runtime.Caller(1); ok {
		return slog.GroupValue(
			slog.String("file", file),
			slog.Int("no", number),
		)
	}
	return slog.StringValue("N/A")
}
