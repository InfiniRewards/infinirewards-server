package logs

import (
	"context"
	connections "infinirewards/nats"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var Logger *slog.Logger
var LogLevel slog.Level = slog.LevelInfo

// multiWriter implements io.Writer to write to multiple destinations
type multiWriter struct {
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return len(p), nil
}

// natsLogWriter implements io.Writer for NATS logging
type natsLogWriter struct {
	subject string
	nc      *nats.Conn
}

func (n *natsLogWriter) Write(p []byte) (int, error) {
	if n.nc != nil {
		err := n.nc.Publish(n.subject, p)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

// customHandler implements slog.Handler with dual logging capability
type customHandler struct {
	terminalHandler slog.Handler
	natsHandler     slog.Handler
	level           slog.Level
	natsOnly        bool // if true, only logs to NATS
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	// Always log to terminal unless it's NATS-only
	if !h.natsOnly {
		if err := h.terminalHandler.Handle(ctx, r.Clone()); err != nil {
			return err
		}
	}

	// Log to NATS if connection is available
	if h.natsHandler != nil {
		if err := h.natsHandler.Handle(ctx, r.Clone()); err != nil {
			return err
		}
	}

	return nil
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	termHandler := h.terminalHandler.WithAttrs(attrs)
	var natsHandler slog.Handler
	if h.natsHandler != nil {
		natsHandler = h.natsHandler.WithAttrs(attrs)
	}
	return &customHandler{
		terminalHandler: termHandler,
		natsHandler:     natsHandler,
		level:          h.level,
		natsOnly:       h.natsOnly,
	}
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	termHandler := h.terminalHandler.WithGroup(name)
	var natsHandler slog.Handler
	if h.natsHandler != nil {
		natsHandler = h.natsHandler.WithGroup(name)
	}
	return &customHandler{
		terminalHandler: termHandler,
		natsHandler:     natsHandler,
		level:          h.level,
		natsOnly:       h.natsOnly,
	}
}

// InitHandler initializes the logging system
func InitHandler(id string) {
	// Create terminal handler with text format
	terminalHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: LogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Format time for better readability
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	})

	// Create NATS handler if connection is available (keep JSON format)
	var natsHandler slog.Handler
	if connections.NC != nil {
		natsWriter := &natsLogWriter{
			subject: "infinirewards.log." + id,
			nc:      connections.NC,
		}
		natsHandler = slog.NewJSONHandler(natsWriter, &slog.HandlerOptions{
			Level: LogLevel,
		})
	}

	// Create custom handler
	handler := &customHandler{
		terminalHandler: terminalHandler,
		natsHandler:     natsHandler,
		level:          LogLevel,
		natsOnly:       false,
	}

	Logger = slog.New(handler)
}

// UpdateLogLevel updates the logging level
func UpdateLogLevel(ctx context.Context, level int) {
	LogLevel = slog.Level(level)
	Logger.Log(ctx, LogLevel, "log level updated",
		slog.String("handler", "UpdateLogLevel"),
		slog.String("new_level", LogLevel.String()),
	)
}

// CreateNatsOnlyLogger creates a logger that only writes to NATS
func CreateNatsOnlyLogger(id string) *slog.Logger {
	if connections.NC == nil {
		return Logger // Fallback to main logger if NATS is unavailable
	}

	natsWriter := &natsLogWriter{
		subject: "infinirewards.log." + id,
		nc:      connections.NC,
	}

	handler := &customHandler{
		terminalHandler: nil,
		natsHandler:     slog.NewJSONHandler(natsWriter, nil),
		level:          LogLevel,
		natsOnly:       true,
	}

	return slog.New(handler)
}
