// Điểm khởi động: đọc config, nối các layer lại với nhau, chạy server.
// Chạy (trong thư mục backend): go run ./cmd/api
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"seatly/internal/config"
	"seatly/internal/domain"
	"seatly/internal/httpapi"
	"seatly/internal/repository/memory"
	"seatly/internal/service"
)

func main() {
	// Log dạng JSON: 1 dòng 1 sự kiện, máy đọc được (Datadog, Loki...).
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	// Nối dây (dependency injection bằng tay): repo → service → handler.
	repo := memory.NewSeatRepository(seedSeats())
	svc := service.NewBookingService(repo)
	router := httpapi.NewRouter(svc)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second, // chặn client gửi header chậm cố tình giữ kết nối
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Ctrl+C hoặc `docker stop` (SIGTERM) → ctx bị hủy.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("🎫 Seatly BE started", "addr", "http://localhost:"+cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err // không mở được cổng, v.v.
	case <-ctx.Done():
	}

	// Graceful shutdown: ngừng nhận request mới, chờ request đang chạy xong (tối đa 10s).
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// seedSeats tạo 3 hàng A/B/C × 5 ghế. Khi có DB (v6) dữ liệu này sẽ nằm trong migration.
func seedSeats() []domain.Seat {
	var seats []domain.Seat
	for _, row := range []string{"A", "B", "C"} {
		price := 500_000
		if row == "A" {
			price = 1_000_000
		}
		for n := 1; n <= 5; n++ {
			seats = append(seats, domain.Seat{ID: row + strconv.Itoa(n), Row: row, Number: n, Price: price})
		}
	}
	return seats
}
