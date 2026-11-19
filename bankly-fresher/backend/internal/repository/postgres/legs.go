package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"bankly/internal/domain"
)

type lockedAccount struct {
	status   string
	balance  int64
	isSystem bool
}

// applyLegs khóa tài khoản rồi ghi từng chân của giao dịch vào sổ cái.
func applyLegs(ctx context.Context, tx pgx.Tx, txID int64, p domain.Posting) error {
	ids := make([]int64, 0, len(p.Legs))
	for _, leg := range p.Legs {
		ids = append(ids, leg.AccountID)
	}

	// Khóa bi quan (pessimistic lock), LUÔN theo thứ tự id tăng dần.
	// A→B và B→A cùng lúc đều khóa min(A,B) trước → không thể khóa chéo nhau (deadlock).
	// Nếu khóa lần lượt "nguồn rồi đích", A→B giữ A chờ B, B→A giữ B chờ A → Postgres báo deadlock.
	rows, err := tx.Query(ctx,
		`SELECT id, status, balance, is_system FROM accounts WHERE id = ANY($1) ORDER BY id FOR UPDATE`, ids)
	if err != nil {
		return err
	}
	accounts := map[int64]lockedAccount{}
	var id int64
	var a lockedAccount
	if _, err := pgx.ForEachRow(rows, []any{&id, &a.status, &a.balance, &a.isSystem}, func() error {
		accounts[id] = a
		return nil
	}); err != nil {
		return err
	}
	if len(accounts) != len(ids) {
		return domain.ErrAccountNotFound
	}

	for _, leg := range p.Legs {
		acc := accounts[leg.AccountID]
		if acc.status == domain.StatusFrozen {
			return domain.ErrAccountFrozen
		}
		delta := p.Amount
		if leg.Direction == domain.Debit {
			// Kiểm tra số dư SAU khi đã khóa → không ai chen vào rút tiền giữa chừng (double spending).
			if !acc.isSystem && acc.balance < p.Amount {
				return domain.ErrInsufficientFunds
			}
			delta = -p.Amount
		}

		var balanceAfter int64
		if err := tx.QueryRow(ctx, `UPDATE accounts SET balance = balance + $2 WHERE id = $1 RETURNING balance`,
			leg.AccountID, delta).Scan(&balanceAfter); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO ledger_entries (transaction_id, account_id, direction, amount, balance_after)
			VALUES ($1, $2, $3, $4, $5)`,
			txID, leg.AccountID, leg.Direction, p.Amount, balanceAfter); err != nil {
			return err
		}
	}
	return nil
}
