// Package task chứa domain model của DevBoard.
// Lab 1: chỉ dùng thư viện chuẩn, không framework.
package task

import (
	"errors"
	"time"
)

// Status là "enum" kiểu Go: một named type + các hằng số.
type Status string

const (
	StatusTodo Status = "todo"
	StatusDone Status = "done"
)

// Sentinel errors: caller dùng errors.Is(err, task.ErrNotFound) để kiểm tra.
var (
	ErrEmptyTitle   = errors.New("title must not be empty")
	ErrTitleTooLong = errors.New("title must be at most 100 characters")
	ErrNotFound     = errors.New("task not found")
	ErrAlreadyDone  = errors.New("task already done")
)

const MaxTitleLen = 100

// Task — để ý struct tag `json:"..."` quyết định tên field khi encode.
type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	DoneAt    *time.Time `json:"done_at,omitempty"` // pointer = có thể nil
}

// New tạo task mới sau khi validate title.
//
// TODO(lab1):
//   - Cắt khoảng trắng 2 đầu title (gợi ý: strings.TrimSpace)
//   - Title rỗng       -> trả ErrEmptyTitle
//   - Title > 100 ký tự -> trả ErrTitleTooLong
//     (đếm theo ký tự, không theo byte! "Việt" có 4 ký tự nhưng 6 byte.
//     Gợi ý: utf8.RuneCountInString)
//   - Status mặc định là StatusTodo, CreatedAt = now
func New(id int, title string, now time.Time) (Task, error) {
	return Task{}, errors.New("TODO: implement New")
}

// MarkDone chuyển task sang done.
// Receiver là pointer (*Task) vì method này sửa dữ liệu của task.
//
// TODO(lab1):
//   - Nếu đã done -> trả ErrAlreadyDone
//   - Ngược lại set Status = StatusDone, DoneAt = &now
func (t *Task) MarkDone(now time.Time) error {
	return errors.New("TODO: implement MarkDone")
}
