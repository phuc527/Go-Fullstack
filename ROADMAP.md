# Lộ trình Golang cho Fullstack Dev (React/Node) — 72 ngày

> Dựa trên [roadmap.sh/golang](https://roadmap.sh/golang), cắt gọn cho người đã có 5+ năm JS/TS.
> Bỏ qua "lập trình là gì", tập trung vào **khác biệt so với Node** và cách viết Go "chuẩn" (idiomatic).

## Cách học

- **1.5–2h/ngày, 5 ngày/tuần** (~14–15 tuần). Mỗi ngày: 20–30' đọc + phần còn lại làm lab.
- **Mỗi lab = 1 commit**. Từ Phase 2 trở đi, lab nào cũng phải có `go test`.
- Cuối mỗi phase có **project** (3–6 ngày) + **câu hỏi tự kiểm tra** (cũng là câu hỏi phỏng vấn Go hay gặp).
- Đang dùng Go 1.27 → dùng luôn feature mới: `for i := range 10`, `ServeMux` pattern `GET /users/{id}`, `log/slog`, iterators, `b.Loop()`, `testing/synctest`, `wg.Go()`.

```
Go-Fullstack/
├── labs/day01/ ... day66/
├── projects/p1-tasker/ ... p6-taskhub/
└── ROADMAP.md
```

## Bảng ánh xạ Node → Go

| Node / TS                 | Go                                               |
| ------------------------- | ------------------------------------------------ |
| `package.json`, `npm i`   | `go.mod`, `go get`                               |
| `class` + `extends`       | `struct` + method + embedding (không có kế thừa) |
| TS interface (structural) | interface (implicit, không cần `implements`)     |
| `try/catch/throw`         | trả về `error`, `if err != nil`                  |
| `async/await`, Promise    | goroutine + channel                              |
| `Promise.all`             | `sync.WaitGroup` / `errgroup`                    |
| `AbortController`         | `context.Context`                                |
| Express middleware        | `func(http.Handler) http.Handler`                |
| Jest                      | `testing` + table-driven tests                   |
| ESLint / Prettier         | `go vet`, `golangci-lint` / `gofmt`              |
| Prisma / TypeORM          | `sqlc` / GORM                                    |
| zod                       | `go-playground/validator`                        |
| pino / winston            | `log/slog`                                       |
| nodemon                   | `air`                                            |
| pm2 cluster               | Không cần — 1 binary dùng hết CPU core           |

---

## Phase 1 — Nền tảng (Day 1–12)

_roadmap.sh: Learn the Basics_

| Day | Chủ đề                                                                                                             | Lab                                                                                                |
| --- | ------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------- |
| 1   | Cài đặt, `go mod init`, `go run/build/fmt/vet`, `package main`                                                     | CLI `hello` nhận tên qua `os.Args`. Build binary Linux từ Windows (`GOOS=linux go build`)          |
| 2   | Biến, kiểu, zero value, `const`/`iota`, ép kiểu tường minh, `string` vs `[]byte` vs `rune`                         | Hàm đảo chuỗi Unicode đúng: `"Việt Nam"` → `"maN têiV"`. Viết thêm bản sai (đảo bytes) để thấy lỗi |
| 3   | `if` có init statement, `for` (vòng lặp duy nhất), `range` over int, `switch` không cần `break`, labeled break     | Sàng Eratosthenes in số nguyên tố ≤ N; FizzBuzz bằng `switch` không điều kiện                      |
| 4   | Array vs slice (`len`/`cap`, `append`, chung backing array), `copy`, package `slices`                              | Tự viết `Filter`, `Map`, `Chunk` cho `[]int`. Tái hiện bug slice aliasing rồi fix                  |
| 5   | `map`, `struct`, struct tags, embedding                                                                            | Word frequency: đọc file text, in top 10 từ (`maps` + `slices.SortFunc`)                           |
| 6   | Hàm: multi-return, variadic, closure, `defer` (LIFO, tham số evaluate ngay)                                        | `timed(fn)` đo thời gian chạy; counter generator bằng closure; demo thứ tự `defer`                 |
| 7   | Pointer, method, value vs pointer receiver                                                                         | `BankAccount` với `Deposit/Withdraw`; chứng minh value receiver không đổi state                    |
| 8   | `error` interface, `fmt.Errorf("%w")`, `errors.Is/As`, sentinel error, `panic/recover` (và khi nào **không** dùng) | Parser file `key=value` trả lỗi kèm số dòng; wrap & unwrap lỗi                                     |
| 9   | Package, module, visibility (chữ hoa = export), `internal/`, `go get`, `go.sum`                                    | Tách Day 5 thành package `wordcount`, gọi từ `cmd/wc`                                              |

### 🚀 Project 1: `tasker` — CLI quản lý todo (Day 10–12)

- Subcommand: `add`, `list`, `done <id>`, `delete <id>`, `list --status=pending` (dùng `flag.NewFlagSet`)
- Lưu file JSON (`encoding/json`), tự tạo file nếu chưa có
- Lỗi rõ ràng ra `stderr`, exit code ≠ 0 khi lỗi
- Build ra 3 binary: Windows/Linux/macOS
- **Stretch**: viết lại CLI bằng `cobra`

**Tự kiểm tra:** Slice khác array thế nào, khi nào `append` tạo mảng mới? Value vs pointer receiver chọn sao? Vì sao Go không có exception?

---

## Phase 2 — Đi sâu: Interface, Generics, Testing, I/O (Day 13–23)

_roadmap.sh: Going Deeper, Testing_

| Day | Chủ đề                                                                                               | Lab                                                                          |
| --- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| 13  | Interface implicit, interface nhỏ (`io.Reader`, `fmt.Stringer`), "accept interfaces, return structs" | `Shape` interface (`Area/Perimeter`) cho Circle/Rect + implement `Stringer`  |
| 14  | Type assertion, type switch, `any`, bẫy **nil interface ≠ nil**                                      | Hàm `describe(v any)`; tái hiện bug trả `*MyErr(nil)` thành `error` khác nil |
| 15  | Generics: type params, constraints, `comparable`, `cmp.Ordered`, generic types                       | `Stack[T]`, `Set[T comparable]`, `Map/Filter/Reduce` generic                 |
| 16  | Table-driven tests, `t.Run`, `t.Helper`, `testdata/`, `go test -cover`                               | Test Stack/Set đạt ≥ 90% coverage                                            |
| 17  | Benchmark (`b.Loop()`, `-benchmem`), fuzzing                                                         | Benchmark nối chuỗi `+` vs `strings.Builder`; fuzz hàm reverse Day 2         |
| 18  | I/O: compose `io.Reader/Writer`, `bufio.Scanner`, `io.Copy`, `embed`                                 | `mygrep -n pattern file` đọc stream, không load hết file vào RAM             |
| 19  | Custom error types, `errors.Join`; iterators `iter.Seq` (range-over-func)                            | `Lines(r io.Reader) iter.Seq2[int, string]` và dùng với `for range`          |
| 20  | Project layout (`cmd/`, `internal/`), `golangci-lint`, `go doc`, Example tests                       | Setup `golangci-lint` + `Makefile`/`Taskfile` cho cả repo                    |

### 🚀 Project 2: `logscope` — Log analyzer CLI (Day 21–23)

- Đọc nginx access log dạng **stream** (file vài trăm MB) — tự viết script sinh log giả
- Thống kê: top IP, top endpoint, phân bố status code, p50/p95 latency
- Output qua interface `Formatter`: `--format=table|json|csv`
- Table-driven test + benchmark + fuzz cho parser
- **Done khi**: coverage ≥ 80%, xử lý file 500MB mà RAM < ~50MB

**Tự kiểm tra:** Khi nào nên định nghĩa interface (bên dùng hay bên cung cấp)? Generics vs interface chọn cái nào? Vì sao `io.Reader` mạnh?

---

## Phase 3 — Concurrency (Day 24–34)

_roadmap.sh: Goroutines, Channels, Select, Mutex, Context, Concurrency Patterns_

| Day | Chủ đề                                                                          | Lab                                                                              |
| --- | ------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 24  | Goroutine vs event loop Node, scheduler G-M-P, `sync.WaitGroup` (`wg.Go`)       | Fetch 20 URL tuần tự vs song song, đo thời gian                                  |
| 25  | Channel unbuffered/buffered, `close`, `range` channel, channel một chiều        | Pipeline: generator → square → sum                                               |
| 26  | `select`, timeout `time.After`, `default` non-blocking, nil channel trick       | `firstResponse()` gọi 3 mirror, lấy cái nhanh nhất, timeout 2s                   |
| 27  | `Mutex`, `RWMutex`, `Once`, `atomic`, race detector `-race`                     | Counter có race → phát hiện bằng `-race` → fix 3 cách (mutex / atomic / channel) |
| 28  | `context`: cancel, timeout, ctx là tham số đầu tiên, không lạm dụng `WithValue` | Worker dừng sạch khi ctx cancel (so với `AbortController`)                       |
| 29  | Worker pool, fan-out/fan-in, semaphore (buffered chan), `errgroup`              | Hash 1000 file bằng pool N worker; `errgroup.WithContext` dừng ở lỗi đầu tiên    |
| 30  | Rate limit (`time.Ticker`, `x/time/rate`), backpressure, goroutine leak         | Tái hiện goroutine leak, đo `runtime.NumGoroutine`, fix                          |
| 31  | Test code concurrent với `testing/synctest`, `signal.NotifyContext`             | Test rate limiter bằng `synctest` (không `sleep` thật)                           |

### 🚀 Project 3: `crawlr` — Concurrent crawler / dead link checker (Day 32–34)

- Input: URL gốc + `--depth`, `--workers`, `--rps`
- Worker pool, dedupe URL, rate limit theo host, timeout mỗi request bằng `context`
- `Ctrl+C` → graceful shutdown, vẫn in report phần đã crawl
- Report: broken links (4xx/5xx), redirect chain
- **Done khi**: `go test -race` sạch, không leak goroutine (`go.uber.org/goleak`)

**Tự kiểm tra:** Unbuffered vs buffered channel? Ai nên `close` channel? Mutex hay channel — khi nào chọn cái nào? Goroutine leak xảy ra khi nào?

---

## Phase 4 — Web API với `net/http` (Day 35–46)

_roadmap.sh: Web Development, Logging, Auth_

| Day | Chủ đề                                                                                     | Lab                                                           |
| --- | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------- |
| 35  | `Handler`, `HandlerFunc`, `ServeMux` pattern (`GET /notes/{id}`), `r.PathValue`            | CRUD in-memory `/notes`, **không framework**                  |
| 36  | JSON encode/decode, `DisallowUnknownFields`, `MaxBytesReader`, helper `writeJSON/readJSON` | Chuẩn hoá error response `{ "error": { "code", "message" } }` |
| 37  | Middleware chain: logging, recover panic, request ID, CORS                                 | Tự viết 4 middleware trên, không dùng lib                     |
| 38  | Validation (`validator`), config từ env (12-factor)                                        | Validate payload Notes + struct config                        |
| 39  | `log/slog`: JSON handler, level theo env, attrs                                            | Log mỗi request: `request_id`, latency, status                |
| 40  | `httptest.NewRecorder`, `httptest.NewServer`, test middleware                              | Test toàn bộ endpoint Notes ≥ 80%                             |
| 41  | `http.Server` timeouts, graceful shutdown, health check. Tour chi / gin / echo / fiber     | Port Notes sang `chi` trong 1h, so sánh với stdlib            |
| 42  | bcrypt, JWT (`golang-jwt`), auth middleware, đưa user vào ctx                              | Thêm register/login, bảo vệ route                             |

### 🚀 Project 4: `shorty` — URL Shortener API (Day 43–46)

- `POST /links`, `GET /{code}` (redirect), `GET /links/{code}/stats`
- Storage qua interface `LinkStore` (bản in-memory; Phase 5 thêm Postgres)
- JWT auth, rate limit theo IP, `slog`, graceful shutdown
- Đếm click **bất đồng bộ** qua channel + worker (áp dụng Phase 3)
- Viết OpenAPI spec, test bằng Postman collection + `httptest`
- **Done khi**: e2e test pass, `-race` sạch, load test bằng `hey`/k6 và ghi lại số rps

**Tự kiểm tra:** Vì sao Go không cần framework như Express? Server thiếu timeout gặp rủi ro gì? Truyền user qua `context` đúng cách thế nào?

---

## Phase 5 — Database & Persistence (Day 47–58)

_roadmap.sh: ORMs & DB Access, Caching_

| Day | Chủ đề                                                                                | Lab                                                                      |
| --- | ------------------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| 47  | `database/sql` + `pgx`, connection pool, `sql.ErrNoRows`, Postgres qua docker compose | CRUD users bằng raw SQL                                                  |
| 48  | Migration với `goose` hoặc `golang-migrate`                                           | Schema users/links/clicks, up/down                                       |
| 49  | `sqlc`: viết SQL → sinh code type-safe                                                | Chuyển Day 47 sang sqlc                                                  |
| 50  | GORM: ưu/nhược, N+1, khi nào chọn                                                     | Cùng CRUD bằng GORM, so sánh query log với sqlc                          |
| 51  | Transaction, isolation level, `SELECT ... FOR UPDATE`, pagination cursor vs offset    | Chuyển tiền giữa 2 account trong tx + test 100 goroutine chuyển cùng lúc |
| 52  | Redis (`go-redis`): cache-aside, TTL, rate limiter phân tán                           | Cache `GET /{code}` của shorty                                           |
| 53  | Layering handler → service → store, DI thủ công trong `main` (không framework DI)     | Refactor shorty theo layer, thêm `PostgresLinkStore`                     |
| 54  | Integration test với `testcontainers-go`                                              | Test store với Postgres thật                                             |

### 🚀 Project 5: `taskhub` — Team task management API (Day 55–58)

Mini Trello, nền tảng cho Capstone:

- Entity: users, workspaces, boards, tasks, comments
- Postgres + sqlc + migrations, Redis cache
- JWT access + refresh token, RBAC (owner / member)
- Filter, sort, cursor pagination
- Integration test bằng testcontainers
- Dockerfile multi-stage + `docker compose up` chạy cả stack
- **Done khi**: 1 lệnh chạy được toàn bộ, test xanh, README có API docs

**Tự kiểm tra:** sqlc vs GORM? Cấu hình pool (`MaxOpenConns`...) ảnh hưởng gì? Khi nào cần `FOR UPDATE`? Offset pagination chậm vì sao?

---

## Phase 6 — Production, Microservices & Fullstack (Day 59–72)

_roadmap.sh: Microservices, gRPC, Real-time, Tools_

| Day | Chủ đề                                                                       | Lab                                                       |
| --- | ---------------------------------------------------------------------------- | --------------------------------------------------------- |
| 59  | Docker multi-stage, distroless/scratch, `air` hot reload                     | Image taskhub < 30MB                                      |
| 60  | `pprof` (CPU/heap/goroutine), `go tool trace`, escape analysis `-gcflags=-m` | Cài 1 bottleneck cố ý vào taskhub, tìm bằng pprof rồi fix |
| 61  | Prometheus (`client_golang`), OpenTelemetry tracing                          | `/metrics` + dashboard Grafana                            |
| 62  | gRPC + protobuf (`buf`), unary/streaming, interceptor                        | Tách `notifier` service bằng gRPC                         |
| 63  | NATS (hoặc Kafka): pub/sub, at-least-once, idempotency                       | Event `task.updated` → notifier                           |
| 64  | WebSocket (`coder/websocket`) / SSE, hub broadcast                           | Realtime cập nhật board                                   |
| 65  | GitHub Actions: lint, `test -race`, build, push image; `govulncheck`         | CI cho cả repo                                            |
| 66  | Retry + backoff, circuit breaker, idempotency key, timeout budget            | Retry client cho call gRPC                                |

### 🏆 Capstone: `taskhub` Fullstack Realtime (Day 67–72)

- **Frontend**: React + Vite + TS (thế mạnh của bạn) — REST + WebSocket, cập nhật realtime khi người khác sửa task
- **Backend**:
  - `api` — REST gateway (taskhub Phase 5)
  - `notifier` — gRPC server + NATS consumer (gửi email giả / in-app notification)
  - `realtime` — WebSocket hub nhận event từ NATS
- Observability: logs (slog) + metrics (Prometheus) + traces (OTel) xuyên service
- `docker compose up` chạy toàn bộ; CI xanh
- Deploy lên Fly.io / Railway / VPS
- **Done khi**: demo được 2 tab browser thấy update realtime của nhau, xem được trace 1 request đi qua 3 service

---

## Tài liệu

| Khi nào    | Tài liệu                                                                                                                            |
| ---------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| Phase 1    | [A Tour of Go](https://go.dev/tour), [Go by Example](https://gobyexample.com)                                                       |
| Phase 1–2  | [Effective Go](https://go.dev/doc/effective_go), [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)                 |
| Phase 3    | _Concurrency in Go_ — Katherine Cox-Buday                                                                                           |
| Phase 4–5  | _Let's Go_ & _Let's Go Further_ — Alex Edwards                                                                                      |
| Xuyên suốt | _100 Go Mistakes and How to Avoid Them_ — Teiva Harsanyi                                                                            |
| Style      | [Google Go Style Guide](https://google.github.io/styleguide/go/), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) |
