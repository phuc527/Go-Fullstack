// Package postgres lưu tài khoản và sổ cái trong PostgreSQL.
// Bài toán tiền bạc cần ACID thật (khóa dòng, UNIQUE, transaction), nên không dùng RAM như Seatly.
package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schema string // file SQL được nhúng thẳng vào binary lúc build

// Open kết nối DB (pool nhiều kết nối) và tạo bảng nếu chưa có.
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	if _, err := pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return pool, nil
}

// Repository làm việc với DB. Mọi thao tác ghi tiền nằm trong posting.go.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// querier: thứ chạy được câu SELECT. Cả *pgxpool.Pool và pgx.Tx đều thỏa mãn,
// nên hàm đọc dùng được ở ngoài lẫn bên trong transaction.
type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
