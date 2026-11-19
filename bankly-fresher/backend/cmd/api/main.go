// Điểm khởi động: đọc config, kết nối DB, nối các layer lại với nhau, chạy server.
// Chạy (trong thư mục backend, cần Postgres đang chạy): go run ./cmd/api
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bankly/internal/config"
	"bankly/internal/httpapi"
	"bankly/internal/repository/postgres"
	"bankly/internal/service"
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

	// Ctrl+C hoặc `docker stop` (SIGTERM) → ctx bị hủy.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Nối dây (dependency injection bằng tay): repo → service → handler.
	repo := postgres.NewRepository(pool)
	svc := service.NewBankService(repo)
	router := httpapi.NewRouter(svc, pool)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("🏦 Bankly BE started", "addr", "http://localhost:"+cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	// Graceful shutdown: chờ giao dịch đang chạy xong (tối đa 10s) rồi mới đóng DB.
	// Cắt ngang giữa chừng cũng không mất tiền (DB tự ROLLBACK), nhưng client sẽ không biết kết quả.
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
