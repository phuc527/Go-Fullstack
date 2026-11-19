package postgres_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"bankly/internal/domain"
)

// Mạng chập chờn: client gửi CÙNG 1 key 20 lần cùng lúc → chỉ ghi sổ đúng 1 lần.
func TestIdempotencyUnderConcurrency(t *testing.T) {
	svc, pool := setup(t)
	acc := newAccount(t, svc, 0)
	k := key()

	var wg sync.WaitGroup
	var fresh atomic.Int32
	ids := make([]int64, 20)
	for i := range ids {
		wg.Go(func() {
			tx, replayed, err := svc.Deposit(context.Background(), k, acc, 500)
			if err != nil {
				t.Error(err)
				return
			}
			ids[i] = tx.ID
			if !replayed {
				fresh.Add(1)
			}
		})
	}
	wg.Wait()

	if fresh.Load() != 1 {
		t.Errorf("%d request được ghi sổ, want đúng 1", fresh.Load())
	}
	for _, id := range ids {
		if id != ids[0] {
			t.Fatalf("các lần gửi lại phải trả cùng giao dịch, got %v", ids)
		}
	}
	if got := balanceOf(t, pool, acc); got != 500 {
		t.Errorf("số dư = %d, want 500", got)
	}
}

func TestKeyReusedForDifferentRequest(t *testing.T) {
	svc, _ := setup(t)
	acc := newAccount(t, svc, 0)
	k := key()
	if _, _, err := svc.Deposit(context.Background(), k, acc, 500); err != nil {
		t.Fatal(err)
	}
	_, _, err := svc.Deposit(context.Background(), k, acc, 999)
	if !errors.Is(err, domain.ErrIdempotencyKeyReused) {
		t.Errorf("err = %v, want ErrIdempotencyKeyReused", err)
	}
}

// A→B và B→A chạy cùng lúc 100 lần: không deadlock (nhờ khóa theo id tăng dần), không mất tiền.
func TestConcurrentTransfersNoDeadlock(t *testing.T) {
	svc, pool := setup(t)
	a, b := newAccount(t, svc, 1_000_000), newAccount(t, svc, 1_000_000)

	var wg sync.WaitGroup
	for i := range 100 {
		from, to := a, b
		if i%2 == 1 {
			from, to = b, a
		}
		wg.Go(func() {
			if _, _, err := svc.Transfer(context.Background(), key(), from, to, 1_000); err != nil {
				t.Error(err) // deadlock → Postgres trả lỗi 40P01 ở đây
			}
		})
	}
	wg.Wait()

	if ga, gb := balanceOf(t, pool, a), balanceOf(t, pool, b); ga != 1_000_000 || gb != 1_000_000 {
		t.Errorf("a = %d, b = %d, want 1.000.000 mỗi bên", ga, gb)
	}
}

// Double spending: có 1.000đ, 10 người cùng rút 300đ → chỉ 3 người thành công.
func TestNoOverdraftUnderConcurrency(t *testing.T) {
	svc, pool := setup(t)
	a, b := newAccount(t, svc, 1_000), newAccount(t, svc, 0)

	var wg sync.WaitGroup
	var ok atomic.Int32
	for range 10 {
		wg.Go(func() {
			_, _, err := svc.Transfer(context.Background(), key(), a, b, 300)
			switch {
			case err == nil:
				ok.Add(1)
			case !errors.Is(err, domain.ErrInsufficientFunds):
				t.Error(err)
			}
		})
	}
	wg.Wait()

	if ok.Load() != 3 || balanceOf(t, pool, a) != 100 {
		t.Errorf("thành công %d lần, số dư %d; want 3 lần, còn 100", ok.Load(), balanceOf(t, pool, a))
	}
}
