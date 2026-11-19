package httpapi_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"seatly/internal/domain"
	"seatly/internal/httpapi"
	"seatly/internal/repository/memory"
	"seatly/internal/service"
)

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil))) // tắt log khi test
	m.Run()
}

func newRouter() http.Handler {
	seed := []domain.Seat{{ID: "A1", Row: "A", Number: 1, Price: 100}}
	return httpapi.NewRouter(service.NewBookingService(memory.NewSeatRepository(seed)))
}

func TestBookingAPI(t *testing.T) {
	router := newRouter()

	// Chạy theo thứ tự: A1 được đặt ở dòng đầu nên dòng 2 phải 409.
	tests := []struct {
		name string
		body string
		want int
	}{
		{"đặt được", `{"seatId":"A1"}`, http.StatusCreated},
		{"đặt lại", `{"seatId":"A1"}`, http.StatusConflict},
		{"không tồn tại", `{"seatId":"Z9"}`, http.StatusNotFound},
		{"body không phải JSON", `hello`, http.StatusBadRequest},
		{"thiếu seatId", `{}`, http.StatusBadRequest},
		{"field lạ", `{"seatId":"A1","x":1}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d (body %s)", rec.Code, tt.want, rec.Body)
			}
		})
	}
}

func TestMetrics(t *testing.T) {
	router := newRouter()
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/seats", nil))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	// route phải là pattern, không phải URL thật
	want := `seatly_http_requests_total{method="GET",route="GET /api/seats",status="200"}`
	if !strings.Contains(rec.Body.String(), want) {
		t.Errorf("/metrics thiếu %s", want)
	}
}

func TestListSeats(t *testing.T) {
	rec := httptest.NewRecorder()
	newRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/seats", nil))

	var got struct {
		Success bool          `json:"success"`
		Data    []domain.Seat `json:"data"`
		Message string        `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || !got.Success || len(got.Data) != 1 || got.Data[0].ID != "A1" {
		t.Errorf("got %d %+v", rec.Code, got)
	}
}

func TestErrorEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(`{"seatId":"Z9"}`))
	newRouter().ServeHTTP(rec, req)

	want := `{"success":false,"data":null,"message":"Ghế Z9 không tồn tại"}`
	if got := strings.TrimSpace(rec.Body.String()); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
