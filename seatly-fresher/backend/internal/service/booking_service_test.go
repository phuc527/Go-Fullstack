package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"seatly/internal/domain"
	"seatly/internal/repository/memory"
	"seatly/internal/service"
)

func newService() *service.BookingService {
	seed := []domain.Seat{{ID: "A1", Row: "A", Number: 1, Price: 100}}
	return service.NewBookingService(memory.NewSeatRepository(seed))
}

// Table-driven test (Day 19): mỗi dòng là 1 trường hợp.
func TestBookSeat(t *testing.T) {
	svc := newService()
	ctx := context.Background()

	tests := []struct {
		name    string
		seatID  string
		wantErr error
	}{
		{"đặt được", "A1", nil},
		{"đặt lại ghế đã đặt", "A1", domain.ErrSeatAlreadyBooked},
		{"ghế không tồn tại", "Z9", domain.ErrSeatNotFound},
		{"id rỗng", "   ", domain.ErrInvalidSeatID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.BookSeat(ctx, tt.seatID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookSeatNormalizesID(t *testing.T) {
	seat, err := newService().BookSeat(context.Background(), " a1 ")
	if err != nil || seat.ID != "A1" || !seat.Booked {
		t.Fatalf("got %+v, %v; want A1 booked", seat, err)
	}
}

// 100 người cùng giành A1 → chỉ đúng 1 người thành công.
func TestBookSeatConcurrent(t *testing.T) {
	svc := newService()
	var (
		wg sync.WaitGroup
		mu sync.Mutex
		ok int
	)
	for range 100 {
		wg.Go(func() {
			if _, err := svc.BookSeat(context.Background(), "A1"); err == nil {
				mu.Lock()
				ok++
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	if ok != 1 {
		t.Errorf("%d người đặt được A1, want đúng 1", ok)
	}
}
