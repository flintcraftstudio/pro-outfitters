package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// RecoveryConfig controls panic recovery behavior.
type RecoveryConfig struct {
	LogFunc      func(val any, stack []byte)
	ErrorHandler func(w http.ResponseWriter, r *http.Request, val any)
}

// Recovery returns middleware that catches panics, logs them, and
// returns a 500 response instead of crashing the server process.
func Recovery(config RecoveryConfig) func(http.Handler) http.Handler {
	logFunc := config.LogFunc
	if logFunc == nil {
		logFunc = func(val any, stack []byte) {
			slog.Error("panic recovered",
				"error", val,
				"stack", string(stack),
			)
		}
	}

	errorHandler := config.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(w http.ResponseWriter, r *http.Request, val any) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &recoverWriter{ResponseWriter: w}

			defer func() {
				val := recover()
				if val == nil {
					return
				}

				// Re-panic for ErrAbortHandler — it signals intentional
				// connection termination. Swallowing it leaks connections.
				if val == http.ErrAbortHandler {
					panic(val)
				}

				stack := debug.Stack()
				logFunc(val, stack)

				// Don't write a response body on upgraded connections
				// (WebSocket) — the handshake already happened.
				if r.Header.Get("Connection") == "Upgrade" {
					return
				}

				if rw.written {
					return
				}

				errorHandler(w, r, val)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

type recoverWriter struct {
	http.ResponseWriter
	written bool
}

func (w *recoverWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *recoverWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}
