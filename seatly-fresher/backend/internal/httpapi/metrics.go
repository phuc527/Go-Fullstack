package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics. Prometheus vào GET /metrics mỗi 5s để "scrape" (kéo) số liệu về.
var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "seatly_http_requests_total",
		Help: "Số request HTTP theo route và status.",
	}, []string{"method", "route", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "seatly_http_request_duration_seconds",
		Help:    "Thời gian xử lý request.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})

	// Metric nghiệp vụ: bao nhiêu lượt đặt thành công / bị trùng / sai ghế.
	bookings = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "seatly_bookings_total",
		Help: "Số lượt đặt ghế theo kết quả.",
	}, []string{"result"})
)

// measure ghi metrics cho mỗi request.
// Label route dùng pattern ("POST /api/bookings"), KHÔNG dùng URL thật:
// URL chứa id (vd /api/events/123) sẽ tạo vô số chuỗi số liệu → Prometheus nổ RAM.
func measure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		route := r.Pattern // ServeMux điền sau khi khớp route
		if route == "" {
			route = "unmatched"
		}
		httpRequests.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
		httpDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}
