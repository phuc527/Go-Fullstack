package postgres_test

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"bankly/internal/repository/postgres"
	"bankly/internal/service"
)

// Integration test: chạy với Postgres THẬT. Không có TEST_DATABASE_URL thì bỏ qua (skip).
//
//	docker compose up -d postgres
//	TEST_DATABASE_URL=postgres://bankly:bankly@localhost:5433/bankly_test?sslmode=disable go test -race ./...
//
// Dùng DB riêng bankly_test (tạo bởi postgres/init.sql). Sổ cái không xóa được (append-only),
// nên test không dọn dữ liệu: mỗi test tự mở tài khoản mới.
func setup(t *testing.T) (*service.BankService, *pgxpool.Pool) {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL chưa đặt → bỏ qua integration test")
	}
	pool, err := postgres.Open(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		reconcile(t, pool) // sau MỖI test: sổ sách phải cân
		pool.Close()
	})
	return service.NewBankService(postgres.NewRepository(pool)), pool
}

// key tạo Idempotency-Key ngẫu nhiên, giống crypto.randomUUID() ở frontend.
func key() string { return rand.Text() }

// newAccount mở tài khoản và nạp sẵn tiền (0 = không nạp).
func newAccount(t *testing.T, svc *service.BankService, balance int64) int64 {
	t.Helper()
	ctx := context.Background()
	acc, err := svc.CreateAccount(ctx, t.Name())
	if err != nil {
		t.Fatal(err)
	}
	if balance > 0 {
		if _, _, err := svc.Deposit(ctx, key(), acc.ID, balance); err != nil {
			t.Fatal(err)
		}
	}
	return acc.ID
}

func balanceOf(t *testing.T, pool *pgxpool.Pool, id int64) int64 {
	t.Helper()
	var b int64
	if err := pool.QueryRow(context.Background(), `SELECT balance FROM accounts WHERE id = $1`, id).Scan(&b); err != nil {
		t.Fatal(err)
	}
	return b
}

// reconcile là phép đối soát của kế toán kép:
//  1. Tổng số dư mọi tài khoản (kể cả quỹ ngân hàng) = 0 → tiền không tự sinh ra, không biến mất.
//  2. Số dư cache của từng tài khoản = cộng dồn bút toán trong sổ cái.
func reconcile(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	var total int64
	if err := pool.QueryRow(ctx, `SELECT COALESCE(SUM(balance), 0) FROM accounts`).Scan(&total); err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Errorf("tổng số dư = %d, phải = 0", total)
	}
	var mismatched int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM accounts a
		WHERE a.balance <> (SELECT COALESCE(SUM(CASE e.direction WHEN 'CREDIT' THEN e.amount ELSE -e.amount END), 0)
		                    FROM ledger_entries e WHERE e.account_id = a.id)`).Scan(&mismatched); err != nil {
		t.Fatal(err)
	}
	if mismatched != 0 {
		t.Errorf("%d tài khoản có số dư lệch với sổ cái", mismatched)
	}
}
