// Package service chứa luật nghiệp vụ. Không biết gì về HTTP hay DB cụ thể.
package service

import (
	"context"
	"strings"

	"seatly/internal/domain"
)

// SeatRepository là "hợp đồng" service cần. Interface đặt ở nơi DÙNG nó (kiểu Go, Day 16):
// memory.SeatRepository hôm nay, postgres.SeatRepository ở v6 — service không phải sửa.
type SeatRepository interface {
	List(ctx context.Context) ([]domain.Seat, error)
	Book(ctx context.Context, seatID string) (domain.Seat, error)
}

type BookingService struct {
	repo SeatRepository
}

func NewBookingService(repo SeatRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) ListSeats(ctx context.Context) ([]domain.Seat, error) {
	return s.repo.List(ctx)
}

func (s *BookingService) BookSeat(ctx context.Context, seatID string) (domain.Seat, error) {
	seatID = strings.ToUpper(strings.TrimSpace(seatID))
	if seatID == "" {
		return domain.Seat{}, domain.ErrInvalidSeatID
	}
	return s.repo.Book(ctx, seatID)
}
