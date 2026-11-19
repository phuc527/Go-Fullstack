package httpapi

import (
	"log/slog"
	"net/http"
	"time"
)

// Middleware trong Go: hàm nhận handler, trả handler mới bọc ngoài (giống app.use của Express).

// recoverPanic: 1 request bị panic thì trả 500, server vẫn sống.
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "err", err, "method", r.Method, "path", r.URL.Path)
				w.Header().Set("Connection", "close")
				writeError(w, http.StatusInternalServerError, "Lỗi máy chủ, thử lại sau")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// statusRecorder ghi lại status code để middleware log được.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// logRequests: mỗi request 1 dòng log (method, path, status, thời gian).
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}
