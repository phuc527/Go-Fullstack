// Package domain chứa "từ vựng" của nghiệp vụ ngân hàng: tài khoản, giao dịch, bút toán, lỗi.
// Package này KHÔNG import gì của HTTP hay DB → layer nào cũng dùng được.
package domain

import "time"

// Tiền luôn là int64, đơn vị nhỏ nhất (đồng). KHÔNG BAO GIỜ dùng float64 cho tiền:
// 0.1 + 0.2 = 0.30000000000000004.
const MaxAmount int64 = 1_000_000_000_000 // 1 nghìn tỷ đồng / 1 giao dịch

// SystemAccountID là "Quỹ tiền mặt" của ngân hàng: tiền nạp vào đi ra từ đây.
// Nó được phép âm, nên tổng số dư MỌI tài khoản luôn = 0 (đó là kế toán kép).
const SystemAccountID int64 = 1

const (
	StatusActive = "ACTIVE"
	StatusFrozen = "FROZEN"

	KindDeposit  = "DEPOSIT"
	KindTransfer = "TRANSFER"
	KindReversal = "REVERSAL"

	Debit  = "DEBIT"  // trừ tiền tài khoản
	Credit = "CREDIT" // cộng tiền tài khoản
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	Balance   int64     `json:"balance"` // cached_balance: luôn bằng tổng bút toán của tài khoản
	IsSystem  bool      `json:"isSystem"`
	CreatedAt time.Time `json:"createdAt"`
}

// Entry là 1 dòng trong sổ cái (ledger_entries). Chỉ được thêm, không sửa, không xóa.
type Entry struct {
	ID            int64     `json:"id"`
	TransactionID int64     `json:"transactionId"`
	AccountID     int64     `json:"accountId"`
	Direction     string    `json:"direction"`
	Amount        int64     `json:"amount"`
	BalanceAfter  int64     `json:"balanceAfter"`
	CreatedAt     time.Time `json:"createdAt"`
	Kind          string    `json:"kind"`     // loại giao dịch chứa bút toán này
	Reversed      bool      `json:"reversed"` // giao dịch này đã bị đảo ngược chưa
}

type Transaction struct {
	ID         int64     `json:"id"`
	Kind       string    `json:"kind"`
	Amount     int64     `json:"amount"`
	ReversesID *int64    `json:"reversesId"`
	CreatedAt  time.Time `json:"createdAt"`
	Entries    []Entry   `json:"entries"`
}

// Leg là 1 "chân" của giao dịch: tài khoản nào, bị trừ hay được cộng.
type Leg struct {
	AccountID int64
	Direction string
}

// Posting là yêu cầu ghi sổ mà service gửi xuống repository.
type Posting struct {
	IdempotencyKey string
	Fingerprint    string // tóm tắt nội dung request, để phát hiện dùng lại key cho request KHÁC
	Kind           string
	Amount         int64
	Legs           []Leg  // rỗng khi Kind = REVERSAL: repository tự đảo chân của giao dịch gốc
	ReversesID     *int64 // chỉ có khi Kind = REVERSAL
}
