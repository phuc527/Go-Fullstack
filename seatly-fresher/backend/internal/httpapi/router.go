// Package httpapi là layer HTTP: route, handler, middleware, đọc/ghi JSON, metrics.
package httpapi

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter gắn route vào handler. Thứ tự middleware: log → metrics → recover → route.
func NewRouter(svc BookingService) http.Handler {
	seats := &seatHandler{svc: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeSuccess(w, http.StatusOK, map[string]string{"status": "ok"}, "")
	})
	// ponytail: /metrics chung cổng với API cho gọn; production thật nên tách cổng nội bộ
	// hoặc chặn ở reverse proxy (nginx của FE chỉ chuyển /api nên user ngoài không thấy).
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /api/seats", seats.list)
	mux.HandleFunc("POST /api/bookings", seats.book)

	return logRequests(measure(recoverPanic(mux)))
}
