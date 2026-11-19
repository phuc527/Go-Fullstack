package httpapi

import (
	"errors"
	"net/http"

	"bankly/internal/domain"
)

// Bảng đổi lỗi nghiệp vụ → HTTP status + nhãn metric + câu cho người dùng. Chỉ layer HTTP biết bảng này.
var errorTable = []struct {
	err     error
	status  int
	label   string
	message string
}{
	{domain.ErrInvalidOwner, http.StatusBadRequest, "invalid", "Tên chủ tài khoản phải có 1–100 ký tự"},
	{domain.ErrInvalidAmount, http.StatusBadRequest, "invalid", "Số tiền phải từ 1 đến 1.000.000.000.000 đồng"},
	{domain.ErrInvalidStatus, http.StatusBadRequest, "invalid", `status phải là "ACTIVE" hoặc "FROZEN"`},
	{domain.ErrMissingIdempotencyKey, http.StatusBadRequest, "invalid", "Thiếu header Idempotency-Key (tối đa 100 ký tự)"},
	{domain.ErrSameAccount, http.StatusBadRequest, "invalid", "Tài khoản nguồn và đích phải khác nhau"},
	{domain.ErrSystemAccount, http.StatusBadRequest, "invalid", "Không được dùng trực tiếp tài khoản quỹ của ngân hàng"},
	{domain.ErrAccountNotFound, http.StatusNotFound, "not_found", "Tài khoản không tồn tại"},
	{domain.ErrTransactionNotFound, http.StatusNotFound, "not_found", "Giao dịch không tồn tại"},
	{domain.ErrAccountFrozen, http.StatusConflict, "frozen", "Tài khoản đang bị đóng băng"},
	{domain.ErrAlreadyReversed, http.StatusConflict, "already_reversed", "Giao dịch này đã được đảo ngược rồi"},
	{domain.ErrCannotReverse, http.StatusConflict, "invalid", "Không thể đảo ngược một giao dịch đảo ngược"},
	{domain.ErrInsufficientFunds, http.StatusUnprocessableEntity, "insufficient_funds", "Số dư không đủ"},
	{domain.ErrIdempotencyKeyReused, http.StatusUnprocessableEntity, "key_reused", "Idempotency-Key này đã được dùng cho một yêu cầu khác"},
}

// lookupError tìm lỗi trong bảng. Không có trong bảng = lỗi không lường trước → 500.
func lookupError(err error) (status int, label, message string) {
	for _, e := range errorTable {
		if errors.Is(err, e.err) {
			return e.status, e.label, e.message
		}
	}
	return http.StatusInternalServerError, "error", "Lỗi máy chủ, thử lại sau"
}
