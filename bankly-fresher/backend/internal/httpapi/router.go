// Package httpapi là layer HTTP: route, handler, middleware, đọc/ghi JSON, metrics.
package httpapi

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Pinger: thứ kiểm tra được DB còn sống (*pgxpool.Pool thỏa mãn).
type Pinger interface {
	Ping(ctx context.Context) error
}

// NewRouter gắn route vào handler. Thứ tự middleware: log → metrics → recover → route.
func NewRouter(svc BankService, db Pinger) http.Handler {
	h := &handler{svc: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "DB không phản hồi")
			return
		}
		writeSuccess(w, http.StatusOK, map[string]string{"status": "ok"}, "")
	})
	// ponytail: /metrics chung cổng với API; nginx của FE chỉ chuyển /api nên user ngoài không thấy.
	mux.Handle("GET /metrics", promhttp.Handler())

	mux.HandleFunc("GET /api/accounts", h.listAccounts)
	mux.HandleFunc("POST /api/accounts", h.createAccount)
	mux.HandleFunc("GET /api/accounts/{id}/entries", h.history)
	mux.HandleFunc("PUT /api/accounts/{id}/status", h.setStatus)
	mux.HandleFunc("POST /api/deposits", h.deposit)
	mux.HandleFunc("POST /api/transfers", h.transfer)
	mux.HandleFunc("POST /api/transactions/{id}/reversal", h.reverse)

	return logRequests(measure(recoverPanic(mux)))
}
