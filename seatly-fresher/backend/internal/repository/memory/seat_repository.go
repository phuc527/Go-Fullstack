// Package memory lưu ghế trong RAM. v6 sẽ thêm package postgres có cùng method.
package memory

import (
	"context"
	"slices"
	"sync"

	"seatly/internal/domain"
)

// SeatRepository an toàn khi nhiều goroutine gọi cùng lúc nhờ mutex (Day 28).
type SeatRepository struct {
	mu    sync.Mutex
	seats []domain.Seat
}

func NewSeatRepository(seed []domain.Seat) *SeatRepository {
	return &SeatRepository{seats: slices.Clone(seed)}
}

func (r *SeatRepository) List(ctx context.Context) ([]domain.Seat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Trả bản copy: bên ngoài sửa slice cũng không ảnh hưởng dữ liệu bên trong.
	return slices.Clone(r.seats), nil
}

// Book kiểm tra + đánh dấu đã đặt trong CÙNG 1 lần khóa → không thể đặt trùng.
// (Postgres sau này: UPDATE seats SET booked = true WHERE id = $1 AND booked = false)
func (r *SeatRepository) Book(ctx context.Context, seatID string) (domain.Seat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i := slices.IndexFunc(r.seats, func(s domain.Seat) bool { return s.ID == seatID })
	if i == -1 {
		return domain.Seat{}, domain.ErrSeatNotFound
	}
	if r.seats[i].Booked {
		return domain.Seat{}, domain.ErrSeatAlreadyBooked
	}
	r.seats[i].Booked = true
	return r.seats[i], nil
}
