// Package service chứa luật nghiệp vụ. Không biết gì về HTTP hay DB cụ thể.
package service

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"bankly/internal/domain"
)

// Repository là "hợp đồng" service cần. Interface đặt ở nơi DÙNG nó (kiểu Go).
type Repository interface {
	CreateAccount(ctx context.Context, owner string) (domain.Account, error)
	ListAccounts(ctx context.Context) ([]domain.Account, error)
	SetStatus(ctx context.Context, id int64, status string) (domain.Account, error)
	AccountEntries(ctx context.Context, accountID int64) ([]domain.Entry, error)
	Post(ctx context.Context, p domain.Posting) (domain.Transaction, bool, error)
}

type BankService struct {
	repo Repository
}

func NewBankService(repo Repository) *BankService {
	return &BankService{repo: repo}
}

func (s *BankService) CreateAccount(ctx context.Context, owner string) (domain.Account, error) {
	owner = strings.TrimSpace(owner)
	if owner == "" || utf8.RuneCountInString(owner) > 100 {
		return domain.Account{}, domain.ErrInvalidOwner
	}
	return s.repo.CreateAccount(ctx, owner)
}

func (s *BankService) ListAccounts(ctx context.Context) ([]domain.Account, error) {
	return s.repo.ListAccounts(ctx)
}

func (s *BankService) History(ctx context.Context, accountID int64) ([]domain.Entry, error) {
	return s.repo.AccountEntries(ctx, accountID)
}

func (s *BankService) SetStatus(ctx context.Context, id int64, status string) (domain.Account, error) {
	if status != domain.StatusActive && status != domain.StatusFrozen {
		return domain.Account{}, domain.ErrInvalidStatus
	}
	return s.repo.SetStatus(ctx, id, status)
}

// Deposit nạp tiền: Quỹ tiền mặt (DEBIT) → tài khoản khách (CREDIT).
// bool trả về = true khi đây là request lặp lại (idempotent replay).
func (s *BankService) Deposit(ctx context.Context, key string, accountID, amount int64) (domain.Transaction, bool, error) {
	if err := validate(key, amount, accountID); err != nil {
		return domain.Transaction{}, false, err
	}
	return s.repo.Post(ctx, domain.Posting{
		IdempotencyKey: key,
		Fingerprint:    fmt.Sprintf("DEPOSIT|%d|%d", accountID, amount),
		Kind:           domain.KindDeposit,
		Amount:         amount,
		Legs:           []domain.Leg{{AccountID: domain.SystemAccountID, Direction: domain.Debit}, {AccountID: accountID, Direction: domain.Credit}},
	})
}

// Transfer chuyển tiền: tài khoản nguồn (DEBIT) → tài khoản đích (CREDIT).
func (s *BankService) Transfer(ctx context.Context, key string, fromID, toID, amount int64) (domain.Transaction, bool, error) {
	if err := validate(key, amount, fromID, toID); err != nil {
		return domain.Transaction{}, false, err
	}
	if fromID == toID {
		return domain.Transaction{}, false, domain.ErrSameAccount
	}
	return s.repo.Post(ctx, domain.Posting{
		IdempotencyKey: key,
		Fingerprint:    fmt.Sprintf("TRANSFER|%d|%d|%d", fromID, toID, amount),
		Kind:           domain.KindTransfer,
		Amount:         amount,
		Legs:           []domain.Leg{{AccountID: fromID, Direction: domain.Debit}, {AccountID: toID, Direction: domain.Credit}},
	})
}

// Reverse đảo ngược 1 giao dịch bằng giao dịch mới. Dòng cũ giữ nguyên để kiểm toán.
func (s *BankService) Reverse(ctx context.Context, key string, txID int64) (domain.Transaction, bool, error) {
	if strings.TrimSpace(key) == "" || len(key) > 100 {
		return domain.Transaction{}, false, domain.ErrMissingIdempotencyKey
	}
	return s.repo.Post(ctx, domain.Posting{
		IdempotencyKey: key,
		Fingerprint:    fmt.Sprintf("REVERSAL|%d", txID),
		Kind:           domain.KindReversal,
		ReversesID:     &txID,
	})
}

func validate(key string, amount int64, accountIDs ...int64) error {
	if strings.TrimSpace(key) == "" || len(key) > 100 {
		return domain.ErrMissingIdempotencyKey
	}
	if amount <= 0 || amount > domain.MaxAmount {
		return domain.ErrInvalidAmount
	}
	for _, id := range accountIDs {
		if id == domain.SystemAccountID {
			return domain.ErrSystemAccount
		}
	}
	return nil
}
