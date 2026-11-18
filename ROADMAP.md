# Go Backend Roadmap — cho dev React/Node 5+ năm

> Tham chiếu [roadmap.sh/golang](https://roadmap.sh/golang) + [roadmap.sh/backend](https://roadmap.sh/backend).
> Bỏ phần "lập trình là gì". Tập trung: **Go khác Node chỗ nào** + **bài toán BE thật** (race, transaction, cache, queue, observability).

## Cách học

- **~2h/ngày, 5 ngày/tuần, 11 tuần (55 ngày).** 30' đọc → 90' lab.
- **Mỗi lab có cột "Done khi"** — chưa đạt thì chưa qua ngày. Ưu tiên chứng minh bằng `go test`.
- **Tuần 1–3:** drill ngôn ngữ trong `labs/dayXX/` (có test sẵn, bạn code cho xanh).
- **Tuần 4–11:** build **Seatly** — API đặt vé sự kiện, lớn dần mỗi ngày trong `seatly/`.
- Cuối tuần trả lời phần **Tự kiểm tra** (cũng là câu phỏng vấn Go hay gặp) không nhìn tài liệu.
- Go 1.27 → dùng feature mới: `for i := range n`, ServeMux `GET /events/{id}`, `log/slog`, iterators, `wg.Go()`, `b.Loop()`, `testing/synctest`.

```
Go-Fullstack/
├── go.mod          # 1 module cho cả repo
├── labs/day01..15/ # drill ngôn ngữ
└── seatly/         # project chính (từ Day 16)
```

## Node → Go

| Node / TS                     | Go                                                 |
| ----------------------------- | -------------------------------------------------- |
| `package.json`, `npm i`       | `go.mod`, `go get`                                 |
| class + extends               | struct + method + embedding (không kế thừa)        |
| TS interface + `implements`   | interface ngầm định — chỉ cần đủ method            |
| `try/catch/throw`             | trả `error`, `if err != nil`, `errors.Is/As`       |
| `async/await`, event loop     | goroutine + channel, scheduler dùng hết CPU core   |
| `Promise.all` / `AbortSignal` | `errgroup` / `context.Context`                     |
| Express middleware            | `func(http.Handler) http.Handler`                  |
| Jest                          | `testing` + table-driven, `-race`, fuzz, benchmark |
| ESLint / Prettier             | `go vet`, `golangci-lint` / `gofmt` (không config) |
| Prisma / TypeORM              | `sqlc` + `pgx` (hoặc GORM)                         |
| zod                           | validate tay hoặc `go-playground/validator`        |
| pino                          | `log/slog`                                         |
| nodemon / pm2                 | `air` / không cần — 1 binary                       |

---

## Tuần 1 — Go cho người viết JS

_roadmap.sh: Learn the Basics_

| Day | Chủ đề                                                                                          | Lab → **Done khi**                                                                                         |
| --- | ----------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| 1   | Toolchain, module/package, kiểu cơ bản, `string`/`byte`/`rune`, slice, map                      | `Reverse` Unicode, `WordFreq`, `TopN` + CLI `topwords` → **test xanh, build được binary Linux từ Windows** |
| 2   | Slice internals (`len`/`cap`, shared backing array), package `slices`/`maps`, struct, embedding | Tái hiện bug `append` aliasing rồi fix → **có test fail ở bản lỗi, pass ở bản fix**                        |
| 3   | Function, closure, `defer`, pointer, value vs pointer receiver                                  | `Account.Deposit/Withdraw` → **test chứng minh value receiver không đổi state**                            |
| 4   | Error: `%w`, `errors.Is/As`, sentinel, custom error type; `panic` chỉ cho bug                   | Parser file `.env` trả `*ParseError{Line}` → **test dùng `errors.As` lấy đúng số dòng**                    |
| 5   | Package, visibility, `internal/`, layout `cmd/` + `internal/`, `golangci-lint`                  | CLI `envcheck --file .env` dùng parser Day 4 → **`go vet` + `golangci-lint run` sạch**                     |

**Tự kiểm tra:** Khi nào `append` cấp mảng mới? Chọn value hay pointer receiver thế nào? Vì sao Go không có exception?

## Tuần 2 — Interface, Generics, Testing

_roadmap.sh: Going Deeper, Testing_

| Day | Chủ đề                                                                                      | Lab → **Done khi**                                                                   |
| --- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| 6   | Interface ngầm định, interface nhỏ, `io.Reader/Writer`, "accept interfaces, return structs" | `CountingWriter` bọc `io.Writer` → **`io.Copy` file qua nó, đếm đúng số byte**       |
| 7   | Type switch, `any`, bẫy **typed nil ≠ nil interface**                                       | Tái hiện hàm trả `*MyErr(nil)` thành `error != nil` → **test bắt được bug, rồi fix** |
| 8   | Generics: type params, constraints, `cmp.Ordered`, generic types                            | `Set[T]`, `Map/Filter`, `LRU[K,V]` → **coverage ≥ 90%**                              |
| 9   | Testing sâu: subtests, `t.Helper`, `t.Cleanup`, golden file, fuzz, benchmark `b.Loop()`     | Fuzz `Reverse` Day 1; bench `+` vs `strings.Builder` → **có số `-benchmem` để so**   |
| 10  | `encoding/json` (tags, `omitempty`, custom marshal), `time`, iterators `iter.Seq`           | Stream file JSON-lines 100MB → thống kê → **RAM không tăng theo kích thước file**    |

**Tự kiểm tra:** Interface nhỏ có lợi gì? Vì sao `err != nil` dù trả `nil`? Generics vs interface khi nào dùng cái nào?

## Tuần 3 — Concurrency (điểm khác Node lớn nhất)

_roadmap.sh: Goroutines, Channels, Select, Mutex, Context_

| Day | Chủ đề                                                                           | Lab → **Done khi**                                                              |
| --- | -------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| 11  | Goroutine, channel buffered/unbuffered, `select`, close/range; so với event loop | Pipeline generate → square → sum → **kết quả đúng, không goroutine nào treo**   |
| 12  | `sync.WaitGroup` (`wg.Go`), `Mutex/RWMutex`, `Once`, `atomic`, flag `-race`      | Counter bị race → **`go test -race` báo lỗi, fix bằng 2 cách đều pass**         |
| 13  | `context` (cancel/timeout/value), `errgroup`                                     | Fetch N URL song song, timeout 2s → **1 URL lỗi thì hủy hết, trả lỗi đầu tiên** |
| 14  | Worker pool, fan-in/out, semaphore, rate limit `x/time/rate`                     | Hash 1000 file với concurrency=4 → **benchmark nhanh hơn chạy tuần tự**         |
| 15  | Goroutine leak, deadlock, `testing/synctest`                                     | Fix 3 bug leak cho sẵn → **test dùng `synctest` chứng minh hết leak**           |

**Tự kiểm tra:** Goroutine rẻ hơn thread vì sao? Khi nào dùng channel, khi nào dùng mutex? Ai nên cancel context?

---

## Tuần 4 — HTTP & REST API → bắt đầu Seatly

_roadmap.sh: Web Frameworks, APIs (REST)_

Seatly: organizer tạo event có sơ đồ ghế, user giữ ghế 10 phút rồi thanh toán.

| Day | Chủ đề                                                                   | Lab → **Done khi**                                                        |
| --- | ------------------------------------------------------------------------ | ------------------------------------------------------------------------- |
| 16  | `net/http`, `ServeMux` pattern, JSON request/response                    | `GET/POST /events`, `GET /events/{id}` in-memory → **gọi được bằng curl** |
| 17  | Middleware: logging, recover, request ID, CORS — tự viết, so với Express | 4 middleware chain → **panic trong handler trả 500, server không chết**   |
| 18  | Validation, lỗi chuẩn RFC 9457 (`problem+json`), pagination cursor       | Mọi lỗi cùng 1 format → **test validation trả 422 đúng field**            |
| 19  | Config từ env, server timeouts, graceful shutdown; khi nào cần chi/gin   | **Kill server khi đang có request chậm → request đó vẫn trả về đủ**       |
| 20  | `httptest`, test handler; Postman collection                             | **Handler coverage ≥ 80%; `postman collection run` pass**                 |

## Tuần 5 — PostgreSQL

_roadmap.sh: Relational DB, ORMs (pgx, sqlc, GORM)_

| Day | Chủ đề                                                               | Lab → **Done khi**                                                       |
| --- | -------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| 21  | Docker Compose Postgres, schema design, migration (`goose`)          | Bảng users/events/seats/bookings → **`goose up/down` chạy 2 chiều sạch** |
| 22  | `database/sql` vs `pgx`, connection pool, query timeout qua context  | Repo pgx thay in-memory → **toàn bộ test Tuần 4 vẫn pass**               |
| 23  | `sqlc` vs GORM — trade-off                                           | Chuyển query sang sqlc → **không còn SQL string trong Go code**          |
| 24  | Transaction, isolation level, `SELECT … FOR UPDATE`, optimistic lock | Đặt ghế trong tx → **100 goroutine giành 1 ghế: đúng 1 thành công**      |
| 25  | Index, `EXPLAIN ANALYZE`, N+1, keyset pagination                     | Seed 1M seats → **ghi lại query time trước/sau index**                   |

**Tự kiểm tra:** `FOR UPDATE` vs optimistic lock chọn khi nào? Pool size đặt bao nhiêu? sqlc hơn GORM ở đâu?

## Tuần 6 — Auth & Security

_roadmap.sh backend: Authentication, Web Security (OWASP)_

| Day | Chủ đề                                                          | Lab → **Done khi**                                               |
| --- | --------------------------------------------------------------- | ---------------------------------------------------------------- |
| 26  | Hash password (argon2id/bcrypt), signup/login                   | **Password không bao giờ xuất hiện trong log/DB dạng thô**       |
| 27  | JWT access + refresh token rotation; so với session cookie      | **Refresh token cũ dùng lại → bị từ chối và revoke cả họ token** |
| 28  | RBAC: user / organizer / admin; authz ở service layer           | **User thường tạo event → 403**                                  |
| 29  | OWASP API Top 10: IDOR, mass assignment, SQLi, security headers | **Test: user A xem booking của user B → 404**                    |
| 30  | OAuth2 Google (`x/oauth2`), quản lý secret                      | **Login Google tạo user mới hoặc link user cũ theo email**       |

## Tuần 7 — Cache, Background Job, Idempotency

_roadmap.sh backend: Caching, Message Brokers_

| Day | Chủ đề                                                          | Lab → **Done khi**                                                          |
| --- | --------------------------------------------------------------- | --------------------------------------------------------------------------- |
| 31  | Redis cache-aside, TTL, invalidation, stampede (`singleflight`) | Cache `GET /events/{id}` → **1000 request đồng thời khi miss = 1 query DB** |
| 32  | Giữ ghế 10': Redis `SET NX EX` hoặc cột `held_until`            | **Hết hạn tự nhả ghế; người khác giữ được ngay**                            |
| 33  | Idempotency-Key, webhook ký HMAC                                | Mock payment webhook gửi 2 lần → **chỉ 1 booking confirmed**                |
| 34  | Job queue (`river` trên Postgres), retry, backoff, dead letter  | Email xác nhận qua Mailpit → **SMTP chết rồi sống lại, email vẫn tới**      |
| 35  | Rate limit phân tán (Redis token bucket)                        | **Vượt quota → 429 + `Retry-After`**                                        |

**Tự kiểm tra:** Cache invalidation khi update? Idempotency khác retry chỗ nào? At-least-once nghĩa là gì với handler?

## Tuần 8 — Kiến trúc & Testing thật

_roadmap.sh: Design Patterns, Testing_

| Day | Chủ đề                                                                         | Lab → **Done khi**                                                        |
| --- | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------- |
| 36  | handler → service → repo; DI qua constructor; interface khai báo phía consumer | Refactor Seatly → **không còn interface nào chỉ có 1 impl mà không test** |
| 37  | Integration test với `testcontainers-go`                                       | **`go test ./...` tự dựng Postgres + Redis thật**                         |
| 38  | Fake viết tay vs mock generate; test service layer                             | **Service test chạy < 1s, không cần Docker**                              |
| 39  | OpenAPI spec (viết spec trước hoặc `oapi-codegen`)                             | **Spec khớp API; Postman import spec chạy pass**                          |
| 40  | `golangci-lint` strict, code review checklist                                  | **0 lint warning; tự review 1 PR theo checklist**                         |

## Tuần 9 — Observability & Performance

_roadmap.sh backend: Observability, Scaling_

| Day | Chủ đề                                             | Lab → **Done khi**                                   |
| --- | -------------------------------------------------- | ---------------------------------------------------- |
| 41  | `slog` JSON, correlation ID xuyên request          | **grep 1 request ID ra toàn bộ log của request đó**  |
| 42  | Prometheus metrics (RED), Grafana                  | **Dashboard: rate, error %, p95 latency theo route** |
| 43  | OpenTelemetry tracing HTTP → DB → Redis, Jaeger    | **Thấy 1 trace đặt vé với span DB + Redis**          |
| 44  | `pprof` CPU/heap, escape analysis, giảm allocation | **Giảm alloc/op ở 1 hot path, có số trước/sau**      |
| 45  | Load test k6: 1000 VU giành ghế                    | **p95 < 200ms, 0 oversell**                          |

## Tuần 10 — Messaging, Microservices, Realtime

_roadmap.sh: Microservices (gRPC, protobuf), Message Brokers, Real-time_

| Day | Chủ đề                                             | Lab → **Done khi**                                                      |
| --- | -------------------------------------------------- | ----------------------------------------------------------------------- |
| 46  | NATS JetStream / RabbitMQ / Kafka — khác nhau gì   | Publish `booking.confirmed` → **consumer chết rồi lên lại vẫn nhận đủ** |
| 47  | Transactional outbox                               | **Kill process giữa commit và publish → event không mất**               |
| 48  | Tách notification service; gRPC + protobuf (`buf`) | **Seatly gọi notification qua gRPC, có deadline**                       |
| 49  | Timeout, retry + jitter, circuit breaker           | **Notification chết → Seatly vẫn đặt vé được**                          |
| 50  | WebSocket / SSE sơ đồ ghế realtime                 | **2 tab: tab A giữ ghế, tab B thấy ghế khóa < 1s**                      |

## Tuần 11 — Ship

_roadmap.sh backend: Containerization, CI/CD_

| Day | Chủ đề                                                        | Lab → **Done khi**                               |
| --- | ------------------------------------------------------------- | ------------------------------------------------ |
| 51  | Dockerfile multi-stage, distroless                            | **Image < 25MB, chạy non-root**                  |
| 52  | GitHub Actions: lint, `test -race`, build image               | **PR đỏ khi test fail**                          |
| 53  | Deploy (Fly.io / VPS), migration khi deploy, health/readiness | **Có URL public; deploy mới không downtime**     |
| 54  | `govulncheck`, secret scanning, dependency update             | **0 vuln known**                                 |
| 55  | README kiến trúc + diagram + ADR; mock interview              | **Giải thích được mọi quyết định trong 15 phút** |

---

## Tài liệu chính

- [Tour of Go](https://go.dev/tour) · [Go by Example](https://gobyexample.com) · [Effective Go](https://go.dev/doc/effective_go)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests) — tuần 1–4
- Sách: _Let's Go_ / _Let's Go Further_ (Alex Edwards) — tuần 4–9; _Concurrency in Go_ (Katherine Cox-Buday) — tuần 3
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) · [100 Go Mistakes](https://100go.co)
