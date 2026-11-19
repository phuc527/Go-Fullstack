// Package domain chứa "từ vựng" của nghiệp vụ: Seat và các lỗi nghiệp vụ.
// Package này KHÔNG import gì của HTTP hay DB → layer nào cũng dùng được.
package domain

import "errors"

// Seat là 1 ghế. Tag `json:"id"` đặt tên field khi đổi sang JSON (Day 17).
type Seat struct {
	ID     string `json:"id"`
	Row    string `json:"row"`
	Number int    `json:"number"`
	Price  int    `json:"price"`
	Booked bool   `json:"booked"`
}

// Lỗi nghiệp vụ (sentinel error, Day 13). Handler dùng errors.Is để đổi ra HTTP status.
var (
	ErrInvalidSeatID     = errors.New("invalid seat id")
	ErrSeatNotFound      = errors.New("seat not found")
	ErrSeatAlreadyBooked = errors.New("seat already booked")
)
