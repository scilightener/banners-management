// Package slogdiscard contains a slog.Logger implementation that discards all the messages being logged.
package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger returns a new slog.Logger implementation that discards anything being logged.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler is a struct, implementing slog.Handler so that any log being handled is discarded.
type DiscardHandler struct{}

// NewDiscardHandler return a new instance of DiscardHandler.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
