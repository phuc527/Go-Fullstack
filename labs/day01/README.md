# Day 1 — Toolchain, string/rune, slice, map

## Đọc (30')

- [Tour of Go: Basics](https://go.dev/tour/basics/1) (lướt nhanh) → [More types](https://go.dev/tour/moretypes/1) (đọc kỹ slice/map)
- [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings) — **quan trọng nhất hôm nay**
- [Go Slices: usage and internals](https://go.dev/blog/slices-intro)

## Khác Node chỗ nào

| Node                        | Go                                                                   |
| --------------------------- | -------------------------------------------------------------------- |
| string = UTF-16, `s.length` | string = **bytes UTF-8**, `len(s)` = số byte (`len("ệ") == 3`)       |
| `[...s].reverse().join("")` | `[]rune(s)` → đảo → `string(r)`                                      |
| `s[i]` = ký tự              | `s[i]` = **byte**; `for _, r := range s` mới lặp theo rune           |
| `{}` / `Map`                | `map[string]int` — key thiếu trả **zero value** (`m["x"]++` chạy ok) |
| Array JS                    | slice `[]T` (view lên array có `len`/`cap`)                          |
| `undefined`                 | zero value: `0`, `""`, `nil`, `false`                                |
| `node index.js`             | `go run .` · `go build` → 1 binary, không cần runtime                |

## Lab

Code trong [text.go](text.go) đang `panic("TODO")`. Làm cho test xanh:

```bash
go test ./labs/day01/ -v
```

1. **`Reverse`** — thử bản sai trước (đảo theo byte `s[i]`) để thấy `"Việt"` vỡ ra ký tự rác, rồi sửa bằng `[]rune`.
2. **`WordFreq`** — `strings.FieldsFunc` + `unicode`. Nhớ trả map rỗng (không phải `nil`) với input rỗng? Test dùng `maps.Equal`, thử xem `nil` có pass không và vì sao.
3. **`TopN`** — map trong Go **lặp ngẫu nhiên thứ tự**, nên phải sort. Dùng `slices.SortFunc` + `cmp.Compare` (gợi ý: `cmp.Or` để sort 2 tiêu chí).
4. **CLI** [cmd/topwords/main.go](cmd/topwords/main.go) — làm theo TODO, rồi:

```bash
go run ./labs/day01/cmd/topwords -n 5 < ROADMAP.md
```

5. **Cross-compile** (PowerShell):

```powershell
$env:GOOS="linux"; go build -o topwords ./labs/day01/cmd/topwords; Remove-Item Env:GOOS
```

## Done khi

- [ ] `go test ./labs/day01/` xanh
- [ ] `go vet ./...` không báo gì
- [ ] CLI in đúng top N từ ROADMAP.md
- [ ] Build ra binary Linux từ Windows (đừng commit binary)

## Tự kiểm tra

1. `len("Việt")` bằng mấy? Vì sao?
2. `for i, r := range s` — `i` tăng 1 mỗi vòng không?
3. Ghi vào `nil` map thì sao? Đọc từ `nil` map thì sao?
4. Vì sao Go cố tình random thứ tự lặp map?
