package httpapi

import (
	"net/http"
	"strconv"

	"bankly/internal/domain"
)

// Client tự tạo key (UUID) cho MỖI ý định chuyển tiền và gửi lại ĐÚNG key đó khi retry.
const idempotencyHeader = "Idempotency-Key"

// POST /api/deposits  body: {"accountId": 2, "amount": 500000}
func (h *handler) deposit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AccountID int64 `json:"accountId"`
		Amount    int64 `json:"amount"`
	}
	if err := readJSON(w, r, &body); err != nil {
		badBody(w, domain.KindDeposit, `{"accountId": 2, "amount": 500000}`)
		return
	}
	tx, replayed, err := h.svc.Deposit(r.Context(), r.Header.Get(idempotencyHeader), body.AccountID, body.Amount)
	h.posted(w, r, domain.KindDeposit, tx, replayed, err)
}

// POST /api/transfers  body: {"fromAccountId": 2, "toAccountId": 3, "amount": 100000}
func (h *handler) transfer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FromAccountID int64 `json:"fromAccountId"`
		ToAccountID   int64 `json:"toAccountId"`
		Amount        int64 `json:"amount"`
	}
	if err := readJSON(w, r, &body); err != nil {
		badBody(w, domain.KindTransfer, `{"fromAccountId": 2, "toAccountId": 3, "amount": 100000}`)
		return
	}
	tx, replayed, err := h.svc.Transfer(r.Context(), r.Header.Get(idempotencyHeader),
		body.FromAccountID, body.ToAccountID, body.Amount)
	h.posted(w, r, domain.KindTransfer, tx, replayed, err)
}

// POST /api/transactions/{id}/reversal  (không cần body)
func (h *handler) reverse(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	tx, replayed, err := h.svc.Reverse(r.Context(), r.Header.Get(idempotencyHeader), id)
	h.posted(w, r, domain.KindReversal, tx, replayed, err)
}

// posted trả kết quả ghi sổ và đếm metric.
//   - giao dịch mới:     201
//   - request lặp lại:   200 + header Idempotent-Replayed: true, body y hệt lần đầu
func (h *handler) posted(w http.ResponseWriter, r *http.Request, kind string, tx domain.Transaction, replayed bool, err error) {
	if err != nil {
		_, label, _ := lookupError(err)
		postings.WithLabelValues(kind, label).Inc()
		h.fail(w, r, err)
		return
	}
	msg := "Giao dịch #" + strconv.FormatInt(tx.ID, 10) + " thành công"
	if replayed {
		postings.WithLabelValues(kind, "replayed").Inc()
		w.Header().Set("Idempotent-Replayed", "true")
		writeSuccess(w, http.StatusOK, tx, msg)
		return
	}
	postings.WithLabelValues(kind, "posted").Inc()
	postedAmount.WithLabelValues(kind).Add(float64(tx.Amount))
	writeSuccess(w, http.StatusCreated, tx, msg)
}

func badBody(w http.ResponseWriter, kind, example string) {
	postings.WithLabelValues(kind, "invalid").Inc()
	writeError(w, http.StatusBadRequest, "Body phải có dạng "+example)
}
