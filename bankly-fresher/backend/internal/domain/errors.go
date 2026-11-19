package domain

import "errors"

// Lỗi nghiệp vụ (sentinel error). Handler dùng errors.Is để đổi ra HTTP status.
var (
	ErrInvalidOwner          = errors.New("invalid owner")
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrMissingIdempotencyKey = errors.New("missing idempotency key")
	ErrSameAccount           = errors.New("cannot transfer to the same account")
	ErrSystemAccount         = errors.New("system account cannot be used directly")
	ErrAccountNotFound       = errors.New("account not found")
	ErrTransactionNotFound   = errors.New("transaction not found")
	ErrAccountFrozen         = errors.New("account frozen")
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrIdempotencyKeyReused  = errors.New("idempotency key reused with a different request")
	ErrAlreadyReversed       = errors.New("transaction already reversed")
	ErrCannotReverse         = errors.New("a reversal cannot be reversed")
)
