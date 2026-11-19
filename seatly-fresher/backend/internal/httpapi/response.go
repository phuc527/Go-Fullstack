package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxBodyBytes = 1 << 20 // 1 MB

// envelope là khuôn chung cho MỌI response, FE chỉ cần 1 cách đọc:
//
//	thành công: {"success": true,  "data": ..., "message": "..."}
//	lỗi:        {"success": false, "data": null, "message": "Ghế A1 đã có người đặt"}
//
// HTTP status code vẫn đúng chuẩn (201, 404, 409...), envelope chỉ bổ sung chứ không thay thế.
type envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("write json", "err", err)
	}
}

func writeSuccess(w http.ResponseWriter, status int, data any, message string) {
	writeJSON(w, status, envelope{Success: true, Data: data, Message: message})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, envelope{Success: false, Data: nil, Message: message})
}

// readJSON đọc body vào dst. Chặn body quá lớn, field lạ và nhiều object liền nhau.
func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("body chỉ được chứa 1 JSON object")
	}
	return nil
}
