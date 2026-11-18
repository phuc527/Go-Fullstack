# Lab 1 — DevBoard CLI

Mục tiêu: làm quen với Go qua một CLI quản lý task, chỉ dùng thư viện chuẩn.
Thời lượng gợi ý: 3–5 buổi tối.

## Chạy

```bash
go test ./...                          # ban đầu sẽ FAIL — nhiệm vụ của bạn là làm cho PASS
go test -v -run TestNew ./internal/task # chạy 1 test cụ thể
go run ./cmd/devboard add "Học goroutine"
go run ./cmd/devboard list -status todo
go run ./cmd/devboard done 1
```

## Cấu trúc

```
devboard/
├── go.mod
├── cmd/devboard/main.go        # CLI (đã viết sẵn — đọc để học pattern run(args, out))
└── internal/task/
    ├── task.go                 # TODO: New, MarkDone
    ├── store.go                # TODO: Open, Add, List, Done, Delete, Save
    ├── task_test.go            # test định nghĩa "thế nào là xong"
    └── store_test.go
```

`internal/` là quy ước có hiệu lực: package bên ngoài module không import được.

## Thứ tự làm

1. `task.go`: `New` → `MarkDone`. Chạy `go test -run 'TestNew|TestMarkDone' ./...`
2. `store.go`: `Open` → `Add` → `List` → `Done` → `Delete` → `Save`
3. Chạy CLI thật, mở `tasks.json` xem kết quả

## Định nghĩa "xong"

- [ ] `go test -race ./...` pass toàn bộ
- [ ] `go vet ./...` và `gofmt -l .` không báo gì
- [ ] Không có `panic` trong code của bạn
- [ ] Mọi lỗi trả ra ngoài đều có ngữ cảnh (`fmt.Errorf("...: %w", err)`)
- [ ] Viết `NOTES.md`: 5 điều khác biệt so với cách bạn làm ở Node

## Khái niệm cần nắm sau lab

| Go | Tương đương Node/TS | Ghi chú |
|---|---|---|
| `struct` + tag `json:""` | interface/class + JSON | Không có class, không kế thừa |
| `error` trả về | `throw` / `try-catch` | Lỗi là giá trị, xử lý ngay tại chỗ |
| `errors.Is` + `%w` | `instanceof` / `cause` | Sentinel error + wrapping |
| Pointer receiver `(t *Task)` | `this` mutable | Value receiver nhận bản sao |
| `for _, v := range` | `for...of` | `v` là **bản sao** — bẫy kinh điển |
| `rune` vs `byte` | string length | `len("Việt") == 6` |
| `t.TempDir()`, table test | Jest `beforeEach`, `test.each` | Thư viện chuẩn có sẵn |
| Inject `now func() time.Time` | `jest.useFakeTimers()` | Dependency injection đơn giản |

## Bẫy hay gặp

- Sửa `t` trong `for _, t := range s.tasks` → không có tác dụng. Dùng index.
- Trả `s.tasks` trực tiếp từ `List` → caller sửa được dữ liệu nội bộ (slice chia sẻ mảng nền).
- Quên trường hợp file chưa tồn tại → lần chạy đầu tiên báo lỗi.
- Add thất bại vẫn tăng `nextID` → test `TestStore_AddAssignsIncreasingIDs` bắt được.

## Bonus (tùy chọn)

1. Viết test cho `run()` trong `main.go` bằng `bytes.Buffer` (gợi ý: `t.Setenv("DEVBOARD_FILE", ...)`).
2. Thêm lệnh `edit <id> <title>`.
3. Thêm `-sort created|title` cho `list` (gợi ý: `slices.SortFunc`).
4. Chuyển `Store` thành interface `Repository` — chuẩn bị cho Lab 3 (Postgres).
