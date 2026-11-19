package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"bankly/internal/domain"
)

// Post ghi 1 giao dịch kế toán kép. TẤT CẢ nằm trong 1 DB transaction (ACID):
// lỗi ở bất kỳ bước nào → ROLLBACK toàn bộ, không có chuyện trừ tiền A mà chưa cộng B.
//
//  1. Ghi dòng transactions với idempotency_key UNIQUE (trùng key → trả kết quả cũ)
//  2. Khóa các tài khoản liên quan theo id tăng dần (SELECT ... FOR UPDATE)
//  3. Kiểm tra số dư/trạng thái, cập nhật balance, thêm ledger_entries
//  4. Ghi outbox_events (Transactional Outbox)
func (r *Repository) Post(ctx context.Context, p domain.Posting) (result domain.Transaction, replayed bool, err error) {
	err = pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if p.Kind == domain.KindReversal {
			if err := resolveReversal(ctx, tx, &p); err != nil {
				return err
			}
		}

		// Request trùng key đang chạy song song sẽ ĐỨNG CHỜ ở đây cho tới khi request đầu COMMIT,
		// rồi rơi vào DO NOTHING → không ghi sổ lần 2.
		var id int64
		err := tx.QueryRow(ctx, `
			INSERT INTO transactions (idempotency_key, request_fingerprint, kind, amount, reverses_id)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (idempotency_key) DO NOTHING RETURNING id`,
			p.IdempotencyKey, p.Fingerprint, p.Kind, p.Amount, p.ReversesID).Scan(&id)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			replayed = true
			result, err = replay(ctx, tx, p)
			return err
		case isUniqueViolation(err, "transactions_reverses_id_key"):
			return domain.ErrAlreadyReversed
		case err != nil:
			return err
		}

		if err := applyLegs(ctx, tx, id, p); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO outbox_events (topic, payload)
			VALUES ('transaction.posted', jsonb_build_object('transactionId', $1::bigint, 'kind', $2::text, 'amount', $3::bigint))`,
			id, p.Kind, p.Amount); err != nil {
			return err
		}
		result, err = loadTransaction(ctx, tx, id)
		return err
	})
	return result, replayed, err
}

// replay: key đã được dùng. Cùng nội dung → trả lại giao dịch cũ. Khác nội dung → lỗi.
func replay(ctx context.Context, tx pgx.Tx, p domain.Posting) (domain.Transaction, error) {
	var id int64
	var fingerprint string
	if err := tx.QueryRow(ctx, `SELECT id, request_fingerprint FROM transactions WHERE idempotency_key = $1`,
		p.IdempotencyKey).Scan(&id, &fingerprint); err != nil {
		return domain.Transaction{}, err
	}
	if fingerprint != p.Fingerprint {
		return domain.Transaction{}, domain.ErrIdempotencyKeyReused
	}
	return loadTransaction(ctx, tx, id)
}

// resolveReversal: bút toán đảo ngược = các chân của giao dịch gốc, đổi DEBIT ↔ CREDIT.
func resolveReversal(ctx context.Context, tx pgx.Tx, p *domain.Posting) error {
	orig, err := loadTransaction(ctx, tx, *p.ReversesID)
	if err != nil {
		return err
	}
	if orig.Kind == domain.KindReversal {
		return domain.ErrCannotReverse
	}
	p.Amount, p.Legs = orig.Amount, nil
	for _, e := range orig.Entries {
		dir := domain.Credit
		if e.Direction == domain.Credit {
			dir = domain.Debit
		}
		p.Legs = append(p.Legs, domain.Leg{AccountID: e.AccountID, Direction: dir})
	}
	return nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
