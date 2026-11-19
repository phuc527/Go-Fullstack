package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"seatly/internal/domain"
)

// BookingService là những gì handler cần từ service layer.
type BookingService interface {
	ListSeats(ctx context.Context) ([]domain.Seat, error)
	BookSeat(ctx context.Context, seatID string) (domain.Seat, error)
}

type seatHandler struct {
	svc BookingService
}

// GET /api/seats
func (h *seatHandler) list(w http.ResponseWriter, r *http.Request) {
	seats, err := h.svc.ListSeats(r.Context())
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, seats, "")
}

// POST /api/bookings  body: {"seatId": "A1"}
func (h *seatHandler) book(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SeatID string `json:"seatId"`
	}
	if err := readJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, `Body phải có dạng {"seatId": "A1"}`)
		return
	}

	seat, err := h.svc.BookSeat(r.Context(), body.SeatID)
	switch {
	case err == nil:
		bookings.WithLabelValues("booked").Inc()
		writeSuccess(w, http.StatusCreated, seat, "Đặt ghế "+seat.ID+" thành công")
	// Đổi lỗi nghiệp vụ → HTTP status. Chỉ handler biết về HTTP.
	case errors.Is(err, domain.ErrInvalidSeatID):
		bookings.WithLabelValues("invalid").Inc()
		writeError(w, http.StatusBadRequest, `Body phải có dạng {"seatId": "A1"}`)
	case errors.Is(err, domain.ErrSeatNotFound):
		bookings.WithLabelValues("not_found").Inc()
		writeError(w, http.StatusNotFound, "Ghế "+body.SeatID+" không tồn tại")
	case errors.Is(err, domain.ErrSeatAlreadyBooked):
		bookings.WithLabelValues("conflict").Inc()
		writeError(w, http.StatusConflict, "Ghế "+body.SeatID+" đã có người đặt")
	default:
		bookings.WithLabelValues("error").Inc()
		h.serverError(w, r, err)
	}
}

// Lỗi không lường trước: ghi log chi tiết, nhưng chỉ trả câu chung cho client.
func (h *seatHandler) serverError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "err", err, "method", r.Method, "path", r.URL.Path)
	writeError(w, http.StatusInternalServerError, "Lỗi máy chủ, thử lại sau")
}
