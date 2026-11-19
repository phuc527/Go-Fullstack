package httpapi_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bankly/internal/domain"
	"bankly/internal/httpapi"
)

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil))) // tắt log khi test
	m.Run()
}

// fakeService nhúng interface: chỉ cần viết method mà test dùng (Transfer).
type fakeService struct {
	httpapi.BankService
	replayed bool
	err      error
}

func (f fakeService) Transfer(_ context.Context, _ string, _, _, amount int64) (domain.Transaction, bool, error) {
	return domain.Transaction{ID: 7, Kind: domain.KindTransfer, Amount: amount}, f.replayed, f.err
}

type okDB struct{}

func (okDB) Ping(context.Context) error { return nil }

// Test layer HTTP: lỗi nghiệp vụ → đúng HTTP status, request lặp lại → 200 + header.
func TestTransferStatus(t *testing.T) {
	tests := []struct {
		name         string
		svc          fakeService
		body         string
		want         int
		wantReplayed bool
	}{
		{"giao dịch mới", fakeService{}, `{"fromAccountId":2,"toAccountId":3,"amount":100}`, 201, false},
		{"request lặp lại", fakeService{replayed: true}, `{"fromAccountId":2,"toAccountId":3,"amount":100}`, 200, true},
		{"body sai", fakeService{}, `hello`, 400, false},
		{"field lạ", fakeService{}, `{"amount":1,"x":1}`, 400, false},
		{"thiếu key", fakeService{err: domain.ErrMissingIdempotencyKey}, `{}`, 400, false},
		{"không có tài khoản", fakeService{err: domain.ErrAccountNotFound}, `{}`, 404, false},
		{"đóng băng", fakeService{err: domain.ErrAccountFrozen}, `{}`, 409, false},
		{"không đủ tiền", fakeService{err: domain.ErrInsufficientFunds}, `{}`, 422, false},
		{"key dùng cho request khác", fakeService{err: domain.ErrIdempotencyKeyReused}, `{}`, 422, false},
		{"lỗi lạ", fakeService{err: io.ErrUnexpectedEOF}, `{}`, 500, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/transfers", strings.NewReader(tt.body))
			req.Header.Set("Idempotency-Key", "k1")
			rec := httptest.NewRecorder()
			httpapi.NewRouter(tt.svc, okDB{}).ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d (body %s)", rec.Code, tt.want, rec.Body)
			}
			if got := rec.Header().Get("Idempotent-Replayed") == "true"; got != tt.wantReplayed {
				t.Errorf("Idempotent-Replayed = %v, want %v", got, tt.wantReplayed)
			}
		})
	}
}

func TestErrorDoesNotLeakInternals(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/transfers", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	httpapi.NewRouter(fakeService{err: io.ErrUnexpectedEOF}, okDB{}).ServeHTTP(rec, req)

	want := `{"success":false,"data":null,"message":"Lỗi máy chủ, thử lại sau"}`
	if got := strings.TrimSpace(rec.Body.String()); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

func TestMetrics(t *testing.T) {
	router := httpapi.NewRouter(fakeService{}, okDB{})
	req := httptest.NewRequest(http.MethodPost, "/api/transfers", strings.NewReader(`{"amount":100}`))
	router.ServeHTTP(httptest.NewRecorder(), req)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	want := `bankly_postings_total{kind="TRANSFER",result="posted"}`
	if !strings.Contains(rec.Body.String(), want) {
		t.Errorf("/metrics thiếu %s", want)
	}
}
