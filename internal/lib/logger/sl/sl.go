package sl

import (
	"log/slog"
)

// Err returns slog.Attr with key "error" and value err.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
