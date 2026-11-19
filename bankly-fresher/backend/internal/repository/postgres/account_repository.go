package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"bankly/internal/domain"
)

const accountColumns = `id, owner, currency, status, balance, is_system, created_at`

func scanAccount(row pgx.Row) (domain.Account, error) {
	var a domain.Account
	err := row.Scan(&a.ID, &a.Owner, &a.Currency, &a.Status, &a.Balance, &a.IsSystem, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, domain.ErrAccountNotFound
	}
	return a, err
}

func (r *Repository) CreateAccount(ctx context.Context, owner string) (domain.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx,
		`INSERT INTO accounts (owner) VALUES ($1) RETURNING `+accountColumns, owner))
}

func (r *Repository) ListAccounts(ctx context.Context) ([]domain.Account, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+accountColumns+` FROM accounts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Account, error) {
		return scanAccount(row)
	})
}

// SetStatus đóng băng / mở băng. Đây là UPDATE trên accounts (thông tin tài khoản),
// KHÔNG phải trên sổ cái, nên được phép. Quỹ hệ thống không bị đóng băng.
func (r *Repository) SetStatus(ctx context.Context, id int64, status string) (domain.Account, error) {
	return scanAccount(r.pool.QueryRow(ctx,
		`UPDATE accounts SET status = $2 WHERE id = $1 AND NOT is_system RETURNING `+accountColumns,
		id, status))
}

// AccountEntries trả lịch sử biến động số dư, mới nhất trước.
func (r *Repository) AccountEntries(ctx context.Context, accountID int64) ([]domain.Entry, error) {
	var exists bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1)`, accountID).
		Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAccountNotFound
	}
	// ponytail: lấy 100 dòng mới nhất; thêm phân trang (?before=<id>) khi lịch sử dài hơn.
	return queryEntries(ctx, r.pool, `WHERE e.account_id = $1 ORDER BY e.id DESC LIMIT 100`, accountID)
}

func queryEntries(ctx context.Context, q querier, where string, args ...any) ([]domain.Entry, error) {
	rows, err := q.Query(ctx, `
		SELECT e.id, e.transaction_id, e.account_id, e.direction, e.amount, e.balance_after, e.created_at,
		       t.kind, EXISTS (SELECT 1 FROM transactions r WHERE r.reverses_id = t.id)
		FROM ledger_entries e JOIN transactions t ON t.id = e.transaction_id `+where, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Entry, error) {
		var e domain.Entry
		err := row.Scan(&e.ID, &e.TransactionID, &e.AccountID, &e.Direction, &e.Amount,
			&e.BalanceAfter, &e.CreatedAt, &e.Kind, &e.Reversed)
		return e, err
	})
}

func loadTransaction(ctx context.Context, q querier, id int64) (domain.Transaction, error) {
	var t domain.Transaction
	err := q.QueryRow(ctx, `SELECT id, kind, amount, reverses_id, created_at FROM transactions WHERE id = $1`, id).
		Scan(&t.ID, &t.Kind, &t.Amount, &t.ReversesID, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, domain.ErrTransactionNotFound
	}
	if err != nil {
		return t, err
	}
	t.Entries, err = queryEntries(ctx, q, `WHERE e.transaction_id = $1 ORDER BY e.id`, id)
	return t, err
}
