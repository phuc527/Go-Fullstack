package postgres_test

import (
	"context"
	"errors"
	"testing"

	"bankly/internal/domain"
)

func TestReversal(t *testing.T) {
	svc, pool := setup(t)
	ctx := context.Background()
	a, b := newAccount(t, svc, 1_000), newAccount(t, svc, 0)

	tx, _, err := svc.Transfer(ctx, key(), a, b, 400)
	if err != nil {
		t.Fatal(err)
	}
	rev, _, err := svc.Reverse(ctx, key(), tx.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rev.Kind != domain.KindReversal || *rev.ReversesID != tx.ID || len(rev.Entries) != 2 {
		t.Errorf("reversal sai: %+v", rev)
	}
	if balanceOf(t, pool, a) != 1_000 || balanceOf(t, pool, b) != 0 {
		t.Errorf("sau đảo ngược: a = %d, b = %d", balanceOf(t, pool, a), balanceOf(t, pool, b))
	}

	tests := []struct {
		name    string
		txID    int64
		wantErr error
	}{
		{"đảo ngược lần 2 (key khác)", tx.ID, domain.ErrAlreadyReversed},
		{"đảo ngược 1 bút toán đảo ngược", rev.ID, domain.ErrCannotReverse},
		{"giao dịch không tồn tại", 1 << 60, domain.ErrTransactionNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := svc.Reverse(ctx, key(), tt.txID); !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestFrozenAccount(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	a, b := newAccount(t, svc, 1_000), newAccount(t, svc, 0)

	if _, err := svc.SetStatus(ctx, b, domain.StatusFrozen); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.Transfer(ctx, key(), a, b, 100); !errors.Is(err, domain.ErrAccountFrozen) {
		t.Errorf("err = %v, want ErrAccountFrozen", err)
	}
	if _, err := svc.SetStatus(ctx, domain.SystemAccountID, domain.StatusFrozen); !errors.Is(err, domain.ErrAccountNotFound) {
		t.Errorf("đóng băng quỹ ngân hàng: err = %v, want ErrAccountNotFound", err)
	}
}

// Sổ cái bất biến: chính DB (trigger) từ chối sửa/xóa, kể cả khi ai đó chạy SQL tay.
func TestLedgerIsAppendOnly(t *testing.T) {
	svc, pool := setup(t)
	newAccount(t, svc, 1_000)

	for _, sql := range []string{
		`UPDATE ledger_entries SET amount = 1 WHERE id = (SELECT max(id) FROM ledger_entries)`,
		`DELETE FROM ledger_entries WHERE id = (SELECT max(id) FROM ledger_entries)`,
		`UPDATE transactions SET amount = 1 WHERE id = (SELECT max(id) FROM transactions)`,
		`TRUNCATE ledger_entries`,
	} {
		if _, err := pool.Exec(context.Background(), sql); err == nil {
			t.Errorf("%q phải bị từ chối", sql)
		}
	}
}
