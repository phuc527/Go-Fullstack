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
		Name: "bankly_http_requests_total",
		Help: "Số request HTTP theo route và status.",
	}, []string{"method", "route", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "bankly_http_request_duration_seconds",
		Help:    "Thời gian xử lý request.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})

	// Metric nghiệp vụ: giao dịch theo loại và kết quả
	// (posted, replayed = request lặp lại bị chặn, insufficient_funds, frozen...).
	postings = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bankly_postings_total",
		Help: "Số yêu cầu ghi sổ theo loại giao dịch và kết quả.",
	}, []string{"kind", "result"})

	// Tổng tiền đã ghi sổ thành công. Tiền là int nhưng Prometheus chỉ lưu float64:
	// chỉ dùng để VẼ BIỂU ĐỒ, không bao giờ dùng số này để đối soát.
	postedAmount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bankly_posted_amount_vnd_total",
		Help: "Tổng số tiền (đồng) đã ghi sổ theo loại giao dịch.",
	}, []string{"kind"})
)

// measure ghi metrics cho mỗi request. Label route dùng pattern ("POST /api/transfers"),
// KHÔNG dùng URL thật có id → tránh "cardinality explosion".
func measure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		route := r.Pattern
		if route == "" {
			route = "unmatched"
		}
		httpRequests.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
		httpDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}
