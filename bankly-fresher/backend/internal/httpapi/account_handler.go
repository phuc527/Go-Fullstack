package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"bankly/internal/domain"
)

// BankService là những gì handler cần từ service layer.
type BankService interface {
	CreateAccount(ctx context.Context, owner string) (domain.Account, error)
	ListAccounts(ctx context.Context) ([]domain.Account, error)
	History(ctx context.Context, accountID int64) ([]domain.Entry, error)
	SetStatus(ctx context.Context, id int64, status string) (domain.Account, error)
	Deposit(ctx context.Context, key string, accountID, amount int64) (domain.Transaction, bool, error)
	Transfer(ctx context.Context, key string, fromID, toID, amount int64) (domain.Transaction, bool, error)
	Reverse(ctx context.Context, key string, txID int64) (domain.Transaction, bool, error)
}

type handler struct {
	svc BankService
}

// GET /api/accounts
func (h *handler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.svc.ListAccounts(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, accounts, "")
}

// POST /api/accounts  body: {"owner": "Nguyễn Văn A"}
func (h *handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Owner string `json:"owner"`
	}
	if err := readJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, `Body phải có dạng {"owner": "Nguyễn Văn A"}`)
		return
	}
	acc, err := h.svc.CreateAccount(r.Context(), body.Owner)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeSuccess(w, http.StatusCreated, acc, "Đã mở tài khoản #"+strconv.FormatInt(acc.ID, 10))
}

// GET /api/accounts/{id}/entries
func (h *handler) history(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	entries, err := h.svc.History(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, entries, "")
}

// PUT /api/accounts/{id}/status  body: {"status": "FROZEN"}
func (h *handler) setStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := readJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, `Body phải có dạng {"status": "FROZEN"}`)
		return
	}
	acc, err := h.svc.SetStatus(r.Context(), id, body.Status)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, acc, "")
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id phải là số nguyên dương")
		return 0, false
	}
	return id, true
}

// fail đổi lỗi ra response. Lỗi lạ: ghi log chi tiết, chỉ trả câu chung cho client.
func (h *handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status, _, message := lookupError(err)
	if status == http.StatusInternalServerError {
		slog.Error("internal error", "err", err, "method", r.Method, "path", r.URL.Path)
	}
	writeError(w, status, message)
}
