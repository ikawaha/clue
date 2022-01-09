package log

import (
	"context"
	"sync"
	"time"
)

type (
	// KeyVals represents a list of key/value pairs.
	KeyVals []interface{}

	// Log entry
	Entry struct {
		Time     time.Time
		Severity Severity
		KeyVals  KeyVals
		Message  string
	}

	// Logger implementation
	logger struct {
		options *options
		lock    sync.Mutex
		keyvals []interface{}
		entries []*Entry
		flushed bool
	}

	// Log severity enum
	Severity int

	// private type for context keys
	ctxKey int
)

const (
	SeverityDebug Severity = iota + 1
	SeverityInfo
	SeverityError
)

const (
	ctxLogger ctxKey = iota + 1
)

// Be kind to tests
var timeNow = time.Now

// Context initializes a context for logging.
func Context(ctx context.Context, opts ...LogOption) context.Context {
	var l *logger
	if v := ctx.Value(ctxLogger); v != nil {
		l = v.(*logger)
	} else {
		l = &logger{options: defaultOptions()}
	}
	for _, opt := range opts {
		opt(l.options)
	}
	return context.WithValue(ctx, ctxLogger, l)
}

// Debug logs a debug message.
func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, SeverityDebug, true, msg, keyvals...)
}

// Print logs an info message and ignores buffering.
func Print(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, SeverityInfo, false, msg, keyvals...)
}

// Info logs an info message.
func Info(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, SeverityInfo, true, msg, keyvals...)
}

// Error logs an error message.
func Error(ctx context.Context, msg string, keyvals ...interface{}) {
	Flush(ctx)
	log(ctx, SeverityError, true, msg, keyvals...)
}

// With adds the given key/value pairs to the log context.
func With(ctx context.Context, keyvals ...interface{}) context.Context {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return ctx
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}
	l.keyvals = append(l.keyvals, keyvals...)
	return ctx
}

// Flush flushes the log entries to the writer.
func Flush(ctx context.Context) {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return // do nothing if context isn't initialized
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	l.flush()
}

// logger lock must be held when calling this function.
func (l *logger) flush() {
	if l.flushed {
		return
	}
	for _, e := range l.entries {
		l.options.w.Write(l.options.format(e))
	}
	l.entries = nil // free up memory
	l.flushed = true
}

func log(ctx context.Context, sev Severity, buffer bool, msg string, keyvals ...interface{}) {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return // do nothing if context isn't initialized
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.options.debug && sev == SeverityDebug {
		return
	}
	if l.options.debug && !l.flushed {
		l.flush()
	}

	keyvals = append(l.keyvals, keyvals...)
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}

	e := &Entry{timeNow().UTC(), sev, keyvals, msg}
	if l.flushed || !buffer {
		l.options.w.Write(l.options.format(e))
		return
	}
	l.entries = append(l.entries, e)
}

// Parse extracts the keys and values from the given key/value pairs. The
// resulting slices are of the same length and ordered in the same way.
func (kv KeyVals) Parse() (keys []string, vals []interface{}) {
	if len(kv) == 0 {
		return
	}
	keys = make([]string, len(kv)/2)
	vals = make([]interface{}, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			key = "<INVALID>"
		}
		keys[i/2] = key
		vals[i/2] = kv[i+1]
	}
	return keys, vals
}

// String returns a string representation of the log severity.
func (l Severity) String() string {
	switch l {
	case SeverityDebug:
		return "DEBUG"
	case SeverityInfo:
		return "INFO"
	case SeverityError:
		return "ERROR"
	default:
		return "<INVALID>"
	}
}

// Code returns a 4-character code for the log severity.
func (l Severity) Code() string {
	switch l {
	case SeverityDebug:
		return "DEBG"
	case SeverityInfo:
		return "INFO"
	case SeverityError:
		return "ERRO"
	default:
		return "<INVALID>"
	}
}

// Color returns an escape sequence that colors the output for the given
// severity.
func (l Severity) Color() string {
	switch l {
	case SeverityDebug:
		return "\033[37m"
	case SeverityInfo:
		return "\033[34m"
	case SeverityError:
		return "\033[1;31m"
	default:
		return ""
	}
}
