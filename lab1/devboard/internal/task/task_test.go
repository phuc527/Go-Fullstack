package task

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var fixedNow = time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)

// Table-driven test: pattern chuẩn của Go. Thêm case = thêm 1 dòng.
func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle string
		wantErr   error
	}{
		{"valid", "Học goroutine", "Học goroutine", nil},
		{"trims spaces", "   viết test   ", "viết test", nil},
		{"empty", "", "", ErrEmptyTitle},
		{"only spaces", "    ", "", ErrEmptyTitle},
		{"100 unicode chars ok", strings.Repeat("ệ", 100), strings.Repeat("ệ", 100), nil},
		{"101 chars too long", strings.Repeat("a", 101), "", ErrTitleTooLong},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(1, tc.title, fixedNow)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if got.Title != tc.wantTitle {
				t.Errorf("Title = %q, want %q", got.Title, tc.wantTitle)
			}
			if got.Status != StatusTodo {
				t.Errorf("Status = %q, want %q", got.Status, StatusTodo)
			}
			if !got.CreatedAt.Equal(fixedNow) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, fixedNow)
			}
		})
	}
}

func TestMarkDone(t *testing.T) {
	tk, err := New(1, "deploy", fixedNow)
	if err != nil {
		t.Fatal(err)
	}

	later := fixedNow.Add(time.Hour)
	if err := tk.MarkDone(later); err != nil {
		t.Fatalf("first MarkDone: %v", err)
	}
	if tk.Status != StatusDone {
		t.Errorf("Status = %q, want done", tk.Status)
	}
	if tk.DoneAt == nil || !tk.DoneAt.Equal(later) {
		t.Errorf("DoneAt = %v, want %v", tk.DoneAt, later)
	}

	if err := tk.MarkDone(later); !errors.Is(err, ErrAlreadyDone) {
		t.Errorf("second MarkDone err = %v, want ErrAlreadyDone", err)
	}
}
