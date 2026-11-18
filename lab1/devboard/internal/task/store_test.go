package task

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// newTestStore mở store trong thư mục tạm (tự xóa sau test) với đồng hồ cố định.
func newTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "tasks.json")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	s.now = func() time.Time { return fixedNow }
	return s, path
}

func TestOpen_MissingFileIsEmpty(t *testing.T) {
	s, _ := newTestStore(t)
	if got := s.List(""); len(got) != 0 {
		t.Fatalf("expected empty store, got %d tasks", len(got))
	}
}

func TestOpen_CorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(path); err == nil {
		t.Fatal("expected error for corrupt file, got nil")
	}
}

func TestStore_AddAssignsIncreasingIDs(t *testing.T) {
	s, _ := newTestStore(t)

	a, err := s.Add("task A")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add(""); !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("Add empty err = %v, want ErrEmptyTitle", err)
	}
	b, err := s.Add("task B")
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != 1 || b.ID != 2 {
		t.Errorf("IDs = %d, %d; want 1, 2 (failed Add must not consume an ID)", a.ID, b.ID)
	}
}

func TestStore_ListFilter(t *testing.T) {
	s, _ := newTestStore(t)
	s.Add("one")
	s.Add("two")
	s.Add("three")
	if err := s.Done(2); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		filter Status
		want   int
	}{
		{"", 3},
		{StatusTodo, 2},
		{StatusDone, 1},
	}
	for _, tc := range tests {
		if got := len(s.List(tc.filter)); got != tc.want {
			t.Errorf("List(%q) len = %d, want %d", tc.filter, got, tc.want)
		}
	}
}

func TestStore_ListReturnsCopy(t *testing.T) {
	s, _ := newTestStore(t)
	s.Add("original")

	list := s.List("")
	list[0].Title = "hacked"

	if s.List("")[0].Title != "original" {
		t.Error("List must return a copy; modifying it changed the store")
	}
}

func TestStore_DoneAndDeleteNotFound(t *testing.T) {
	s, _ := newTestStore(t)

	if err := s.Done(42); !errors.Is(err, ErrNotFound) {
		t.Errorf("Done(42) err = %v, want ErrNotFound", err)
	}
	if err := s.Delete(42); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete(42) err = %v, want ErrNotFound", err)
	}
}

func TestStore_Delete(t *testing.T) {
	s, _ := newTestStore(t)
	s.Add("keep")
	s.Add("remove")

	if err := s.Delete(2); err != nil {
		t.Fatal(err)
	}
	list := s.List("")
	if len(list) != 1 || list[0].Title != "keep" {
		t.Errorf("after Delete got %+v", list)
	}

	// ID không được tái sử dụng sau khi xóa.
	c, _ := s.Add("new")
	if c.ID != 3 {
		t.Errorf("new ID = %d, want 3", c.ID)
	}
}

func TestStore_SaveAndReopen(t *testing.T) {
	s, path := newTestStore(t)
	s.Add("persist me")
	s.Add("and me")
	s.Done(1)

	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !errors.Is(err, os.ErrNotExist) {
		t.Error("temp file should be renamed away after Save")
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	list := reopened.List("")
	if len(list) != 2 {
		t.Fatalf("reopened has %d tasks, want 2", len(list))
	}
	if list[0].Status != StatusDone || list[0].DoneAt == nil {
		t.Errorf("task 1 not persisted as done: %+v", list[0])
	}

	reopened.now = func() time.Time { return fixedNow }
	next, _ := reopened.Add("third")
	if next.ID != 3 {
		t.Errorf("nextID not persisted: got %d, want 3", next.ID)
	}
}
