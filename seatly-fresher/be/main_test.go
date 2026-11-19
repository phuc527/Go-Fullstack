package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Chạy: go test -race -v .   (Test sẽ học kỹ ở Day 19)
func TestBooking(t *testing.T) {
	seats = newSeats()
	mux := newMux()

	book := func(body string) int {
		req := httptest.NewRequest("POST", "/api/bookings", strings.NewReader(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec.Code
	}

	if got := book(`{"seatId":"A1"}`); got != http.StatusCreated {
		t.Errorf("đặt A1: got %d, want 201", got)
	}
	if got := book(`{"seatId":"A1"}`); got != http.StatusConflict {
		t.Errorf("đặt A1 lần 2: got %d, want 409", got)
	}
	if got := book(`{"seatId":"Z9"}`); got != http.StatusNotFound {
		t.Errorf("đặt Z9: got %d, want 404", got)
	}
	if got := book(`hello`); got != http.StatusBadRequest {
		t.Errorf("body sai: got %d, want 400", got)
	}
}

// 100 người cùng giành ghế B1 → chỉ đúng 1 người thành công.
func TestBookingConcurrent(t *testing.T) {
	seats = newSeats()
	mux := newMux()

	results := make(chan int, 100)
	for range 100 {
		go func() {
			req := httptest.NewRequest("POST", "/api/bookings", strings.NewReader(`{"seatId":"B1"}`))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			results <- rec.Code
		}()
	}
	ok := 0
	for range 100 {
		if <-results == http.StatusCreated {
			ok++
		}
	}
	if ok != 1 {
		t.Errorf("%d người đặt được B1, want đúng 1", ok)
	}
}
