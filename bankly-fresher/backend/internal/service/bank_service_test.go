package service_test

import (
	"context"
	"errors"
	"testing"

	"bankly/internal/domain"
	"bankly/internal/service"
)

// Các lỗi validate xảy ra TRƯỚC khi chạm DB, nên repo = nil cũng test được.
// Phần ghi sổ thật (khóa, idempotency, đảo ngược) test với Postgres ở repository/postgres.
func TestTransferValidation(t *testing.T) {
	svc := service.NewBankService(nil)
	tests := []struct {
		name     string
		key      string
		from, to int64
		amount   int64
		wantErr  error
	}{
		{"thiếu Idempotency-Key", "", 2, 3, 100, domain.ErrMissingIdempotencyKey},
		{"key toàn khoảng trắng", "   ", 2, 3, 100, domain.ErrMissingIdempotencyKey},
		{"số tiền 0", "k", 2, 3, 0, domain.ErrInvalidAmount},
		{"số tiền âm (lén rút tiền ngược)", "k", 2, 3, -100, domain.ErrInvalidAmount},
		{"vượt hạn mức", "k", 2, 3, domain.MaxAmount + 1, domain.ErrInvalidAmount},
		{"tự chuyển cho mình", "k", 2, 2, 100, domain.ErrSameAccount},
		{"rút từ quỹ ngân hàng", "k", domain.SystemAccountID, 3, 100, domain.ErrSystemAccount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.Transfer(context.Background(), tt.key, tt.from, tt.to, tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestOtherValidation(t *testing.T) {
	svc := service.NewBankService(nil)
	ctx := context.Background()

	if _, err := svc.CreateAccount(ctx, "  "); !errors.Is(err, domain.ErrInvalidOwner) {
		t.Errorf("CreateAccount rỗng: err = %v", err)
	}
	if _, err := svc.SetStatus(ctx, 2, "CLOSED"); !errors.Is(err, domain.ErrInvalidStatus) {
		t.Errorf("SetStatus lạ: err = %v", err)
	}
	if _, _, err := svc.Deposit(ctx, "k", domain.SystemAccountID, 100); !errors.Is(err, domain.ErrSystemAccount) {
		t.Errorf("Deposit vào quỹ: err = %v", err)
	}
	if _, _, err := svc.Reverse(ctx, "", 1); !errors.Is(err, domain.ErrMissingIdempotencyKey) {
		t.Errorf("Reverse thiếu key: err = %v", err)
	}
}
