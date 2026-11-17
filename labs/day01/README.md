# Day 1 — Setup & Toolchain

## Mục tiêu

- Hiểu `module` / `package` / `package main` khác nhau thế nào
- Dùng thành thạo: `go mod init`, `go run`, `go build`, `go fmt`, `go vet`
- Cross-compile ra binary cho OS khác

## Đọc (20–30')

- [Tutorial: Get started with Go](https://go.dev/doc/tutorial/getting-started)
- [Go by Example: Hello World](https://gobyexample.com/hello-world), [Command-Line Arguments](https://gobyexample.com/command-line-arguments), [Exit](https://gobyexample.com/exit)
- VS Code: cài extension **Go** (`golang.go`), rồi `Ctrl+Shift+P` → `Go: Install/Update Tools` → chọn hết

## So với Node

| Node                    | Go                                                |
| ----------------------- | ------------------------------------------------- |
| `npm init`              | `go mod init <module-path>`                       |
| `node index.js`         | `go run .` (compile vào thư mục tạm rồi chạy)     |
| bundler / `pkg`         | `go build` → **1 file binary**, không cần runtime |
| `process.argv.slice(2)` | `os.Args[1:]`                                     |
| `process.exit(1)`       | `os.Exit(1)` (⚠️ không chạy `defer`)              |
| `console.error`         | `fmt.Fprintln(os.Stderr, ...)`                    |
| Prettier                | `go fmt` — 1 style duy nhất, không có config      |

## Lab

### 1. Khởi tạo module

```bash
cd labs/day01
go mod init github.com/<username>/go-fullstack/labs/day01
cd ../..
go work use ./labs/day01   # để VS Code/gopls nhận module này
```

Mở `go.mod` và `go.work` xem chúng chứa gì.

### 2. CLI `hello`

Tạo `main.go` (tự gõ, đừng copy) đáp ứng:

| Lệnh                     | Output                                   | Exit code |
| ------------------------ | ---------------------------------------- | --------- |
| `go run . Phuc`          | `Xin chào, Phuc!`                        | 0         |
| `go run . Phuc Lan Minh` | `Xin chào, Phuc, Lan, Minh!`             | 0         |
| `go run .`               | `usage: hello <name>...` (ra **stderr**) | 1         |

Gợi ý: `os.Args`, `strings.Join`, `fmt.Fprintln`, `os.Exit`.

Kiểm tra exit code: `go run .; echo $?` (bash) hoặc `go run .; $LASTEXITCODE` (PowerShell).

### 3. Toolchain

1. Cố tình viết lệch thụt lề, thừa khoảng trắng → chạy `go fmt ./...` → xem diff.
2. Thêm dòng `fmt.Printf("%d\n", "abc")` → `go run .` vẫn chạy được, nhưng `go vet ./...` báo lỗi. Đọc lỗi rồi xoá dòng đó.
3. Import 1 package mà không dùng (vd `"time"`) → compile lỗi. Go **cấm** import thừa và biến khai báo không dùng.

### 4. Build & cross-compile

```powershell
go build -o hello.exe .
./hello.exe Phuc

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o hello-linux .
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o hello-mac .
Remove-Item Env:GOOS, Env:GOARCH
```

- So sánh kích thước 3 file. Build lại với `-ldflags="-s -w"` → nhỏ hơn bao nhiêu?
- `go env GOOS GOARCH GOPATH GOROOT` — mỗi biến là gì?
- `go tool dist list` — Go build được cho bao nhiêu nền tảng?

## Hoàn thành khi

- [ ] `go.mod` có, `go work use` xong, VS Code không báo lỗi
- [ ] 3 case trong bảng ở bước 2 chạy đúng output **và** exit code
- [ ] `go vet ./...` sạch, `gofmt -l .` không in gì
- [ ] Có 3 binary Windows/Linux/macOS, đã ghi lại kích thước vào `NOTES.md`

## Stretch

- Thêm flag `-upper` (dùng package `flag`) để in hoa: `go run . -upper phuc` → `Xin chào, PHUC!`
- Đọc `go help build` và thử `go build -x` để xem Go thực sự làm gì khi build

## Tự trả lời (ghi vào `NOTES.md`)

1. `module`, `package`, `package main` khác nhau thế nào?
2. `go run` khác `go build` thế nào?
3. Vì sao binary Go chạy được trên server không cài Go, còn app Node thì cần cài Node?
