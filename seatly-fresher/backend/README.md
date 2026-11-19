# Seatly Backend (Go)

API đặt vé viết bằng Go, **chỉ dùng thư viện chuẩn** (`net/http`, `encoding/json`, `log/slog`). Thư viện ngoài duy nhất là Prometheus client để đo metrics.

Cấu trúc tham khảo từ các project Go "clean" trên GitHub, nhưng giữ ở mức fresher đọc được:

- [leeprovoost/go-rest-api-template](https://github.com/leeprovoost/go-rest-api-template): chỉ dùng stdlib, route kiểu `"GET /path"`, graceful shutdown.
- [evrone/go-clean-template](https://github.com/evrone/go-clean-template): tách layer theo Clean Architecture, có `cmd/` + `internal/`.
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout): quy ước đặt tên thư mục `cmd/`, `internal/`.

## Chạy

```bash
cd seatly-fresher/backend
go run ./cmd/api                  # → http://localhost:8080/api/seats
PORT=9000 go run ./cmd/api        # đổi cổng bằng biến môi trường
go test -race ./...               # test tất cả package, bật race detector
go build -o bin/seatly ./cmd/api  # build ra file chạy được
```

`./...` nghĩa là "package này và mọi package con". Lệnh này giống chạy test cho cả monorepo.

## Request đi qua những đâu

```
HTTP request
   │
   ▼
logRequests → measure → recoverPanic        middleware.go, metrics.go
   │
   ▼
ServeMux: "POST /api/bookings"              router.go
   │
   ▼
seatHandler.book                            seat_handler.go     ← biết HTTP: đọc JSON, trả status
   │  gọi
   ▼
BookingService.BookSeat                     booking_service.go  ← luật nghiệp vụ, KHÔNG biết HTTP
   │  gọi qua interface
   ▼
memory.SeatRepository.Book                  seat_repository.go  ← lưu trữ (RAM hôm nay, Postgres sau)
```

Mỗi layer chỉ gọi xuống layer dưới, không gọi ngược lên. So với NestJS: **handler ≈ controller**, **service ≈ service**, **repository ≈ repository/DAO**.

## Từng file

```
backend/
├── cmd/api/main.go
├── internal/
│   ├── config/config.go
│   ├── domain/seat.go
│   ├── service/
│   │   ├── booking_service.go
│   │   └── booking_service_test.go
│   ├── repository/memory/seat_repository.go
│   └── httpapi/
│       ├── router.go
│       ├── seat_handler.go
│       ├── seat_handler_test.go
│       ├── response.go
│       ├── middleware.go
│       └── metrics.go
├── go.mod / go.sum
├── Dockerfile
└── .dockerignore / .gitignore
```

### `cmd/api/main.go`: điểm khởi động

Đây là nơi duy nhất **nối các layer lại** ("dependency injection bằng tay"):

```go
repo := memory.NewSeatRepository(seedSeats())
svc := service.NewBookingService(repo)
router := httpapi.NewRouter(svc)
```

Ngoài ra file này còn làm:

- **Log JSON** bằng `slog`. Mỗi dòng log là 1 JSON object, nên các công cụ như Loki/Datadog đọc và lọc được.
- **Timeout cho server** (`ReadHeaderTimeout`, `ReadTimeout`...). `http.ListenAndServe` mặc định không có timeout, nên client có thể giữ kết nối mãi mãi.
- **Graceful shutdown**. Khi nhận Ctrl+C hoặc `docker stop` (tín hiệu SIGTERM), server ngừng nhận request mới và chờ tối đa 10s cho các request đang chạy xong rồi mới tắt. Nếu không có bước này, người đang đặt ghế dở sẽ bị cắt ngang.
- `seedSeats()` tạo dữ liệu mẫu: 3 hàng × 5 ghế.

Vì sao có thư mục `cmd/`? Một repo có thể có nhiều chương trình, ví dụ `cmd/api`, `cmd/worker`, `cmd/migrate`. Mỗi thư mục là 1 `package main`.

### `internal/`: code riêng của project

Go **cấm** project khác import thứ gì nằm trong `internal/`. Quy tắc này do compiler kiểm tra, không chỉ là quy ước đặt tên.

### `internal/config/config.go`: cấu hình

File này đọc biến môi trường `PORT`; nếu không có thì dùng `8080`. Nó tương đương `process.env.PORT ?? 8080`. Production không hard-code cấu hình, vì cùng 1 image Docker phải chạy được ở cả dev, staging và prod.

### `internal/domain/seat.go`: từ vựng nghiệp vụ

- `Seat`: struct ghế, kèm json tag.
- Các lỗi nghiệp vụ: `ErrSeatNotFound`, `ErrSeatAlreadyBooked`, `ErrInvalidSeatID`.

Package này **không import gì cả**, nên layer nào cũng dùng được mà không bị vòng import. Lỗi được khai báo ở đây là để service và repository trả **cùng một lỗi**, rồi handler dùng `errors.Is` để đổi lỗi đó ra HTTP status.

### `internal/service/booking_service.go`: luật nghiệp vụ

- Chuẩn hóa mã ghế: `" a1 "` → `"A1"`. Mã rỗng thì trả `ErrInvalidSeatID`.
- Khai báo **interface `SeatRepository`** ngay tại service, tức là ở nơi DÙNG nó. Đây là kiểu Go: "accept interfaces, return structs". Sau này đổi sang Postgres thì service không phải sửa dòng nào.

Service hôm nay còn mỏng. Khi có thêm luật như "mỗi người tối đa 4 vé" hay "không đặt khi sự kiện đã bắt đầu", các luật đó sẽ nằm ở đây chứ không nằm trong handler.

### `internal/service/booking_service_test.go`

- **Table-driven test**: mỗi dòng trong bảng là 1 trường hợp, chạy bằng `t.Run`. Đây là cách viết test phổ biến nhất trong Go.
- Test chuẩn hóa mã ghế.
- **100 goroutine cùng giành 1 ghế** thì chỉ đúng 1 người thành công. Chạy với `-race`, Go sẽ báo nếu có data race.

Package của file test là `service_test` (có hậu tố `_test`), nên test chỉ thấy những gì bên ngoài thấy, giống người dùng thật của package.

### `internal/repository/memory/seat_repository.go`: lưu trữ

- Giữ slice ghế trong RAM, có `sync.Mutex` bảo vệ.
- `Book` **kiểm tra và đánh dấu trong cùng 1 lần khóa**. Nếu tách thành 2 bước (đọc xong mới ghi), 2 request có thể cùng thấy "còn trống" và đặt trùng ghế.
- `List` trả về **bản copy** (`slices.Clone`), nên bên ngoài sửa kết quả cũng không làm hỏng dữ liệu bên trong.

Khi lên Postgres sẽ có thêm `repository/postgres/` với **cùng tên method**. Ở `main.go` chỉ cần đổi 1 dòng.

### `internal/httpapi/`: layer HTTP

Tên thư mục không đặt là `http`, để khỏi trùng tên với package `net/http`.

| File                   | Làm gì                                                                                                                                                                                                                                                                  |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `router.go`            | Khai báo route (`GET /api/seats`, `POST /api/bookings`, `GET /healthz`, `GET /metrics`) và bọc middleware.                                                                                                                                                              |
| `seat_handler.go`      | Đọc request, gọi service, **đổi lỗi nghiệp vụ → HTTP status** (`ErrSeatNotFound` → 404, `ErrSeatAlreadyBooked` → 409). Lỗi lạ trả 500 với câu chung chung, còn chi tiết chỉ ghi vào log để không lộ nội bộ cho client. Handler cũng đếm metric `seatly_bookings_total`. |
| `response.go`          | `writeSuccess` / `writeError` bọc mọi response trong envelope `{success, data, message}`. `readJSON` chặn body > 1MB, chặn field lạ (`DisallowUnknownFields`) và chặn body chứa 2 JSON liền nhau.                                                                       |
| `middleware.go`        | `recoverPanic`: 1 request bị panic thì trả 500 và server vẫn sống. `logRequests`: mỗi request ghi 1 dòng log gồm method, path, status, thời gian.                                                                                                                       |
| `metrics.go`           | Khai báo metric Prometheus và middleware `measure` đếm request và đo thời gian xử lý.                                                                                                                                                                                   |
| `seat_handler_test.go` | Test qua HTTP thật bằng `httptest`: 201/409/404/400, danh sách ghế, và `/metrics`.                                                                                                                                                                                      |

Thứ tự middleware: `logRequests(measure(recoverPanic(mux)))`. Request đi từ ngoài vào trong, nên panic bị bắt ở trong cùng. Nhờ vậy log và metrics vẫn ghi nhận status 500.

### `go.mod` / `go.sum`

`go.mod` ≈ `package.json` (tên module `seatly` và danh sách thư viện). `go.sum` ≈ `package-lock.json` (hash để chắc thư viện không bị sửa). Cả 2 file đều phải commit.

### `Dockerfile`

Build 2 bước (multi-stage):

1. Image `golang` build ra 1 file binary tĩnh. `CGO_ENABLED=0` để binary không phụ thuộc thư viện C. `-ldflags="-s -w"` bỏ thông tin debug cho binary nhỏ hơn.
2. Image `alpine` chỉ chứa file binary đó (khoảng 20MB) và chạy bằng user `nobody` (không phải root) cho an toàn.

`go.mod`/`go.sum` được copy **trước** code, nên Docker cache được bước tải thư viện: sửa code xong build lại không phải tải lại.

## Endpoints

Mọi response đều có cùng dạng. HTTP status vẫn đúng chuẩn (200/201/400/404/409/500):

```jsonc
// thành công
{ "success": true, "data": { "id": "A1", "booked": true }, "message": "Đặt ghế A1 thành công" }
// lỗi
{ "success": false, "data": null, "message": "Ghế A1 đã có người đặt" }
```

| Method | Path            | Mô tả                                                                          |
| ------ | --------------- | ------------------------------------------------------------------------------ |
| GET    | `/api/seats`    | Danh sách ghế                                                                  |
| POST   | `/api/bookings` | `{"seatId":"A1"}` → `201` · `400` body sai · `404` không có ghế · `409` đã đặt |
| GET    | `/healthz`      | `{"status":"ok"}`. Docker/Kubernetes gọi để biết app còn sống không            |
| GET    | `/metrics`      | Số liệu cho Prometheus (dạng text)                                             |

## Metrics

| Metric                                 | Loại      | Label                                                            | Dùng để                                                         |
| -------------------------------------- | --------- | ---------------------------------------------------------------- | --------------------------------------------------------------- |
| `seatly_http_requests_total`           | Counter   | `method`, `route`, `status`                                      | Tính request/giây, tỉ lệ lỗi                                    |
| `seatly_http_request_duration_seconds` | Histogram | `method`, `route`                                                | Tính độ trễ p95/p99                                             |
| `seatly_bookings_total`                | Counter   | `result` (`booked`, `conflict`, `not_found`, `invalid`, `error`) | Metric nghiệp vụ: bao nhiêu vé bán được, bao nhiêu lần bị trùng |
| `go_*`, `process_*`                    | có sẵn    |                                                                  | Goroutine, RAM, CPU của chương trình Go                         |

Label `route` dùng **pattern** (`POST /api/bookings`) chứ không dùng URL thật. Lý do: khi có URL dạng `/api/events/123`, mỗi id sẽ sinh ra 1 chuỗi số liệu riêng, và Prometheus sẽ hết RAM (lỗi này gọi là "cardinality explosion").

`/metrics` đang dùng chung cổng với API. Như vậy không lộ ra ngoài, vì nginx của frontend chỉ chuyển `/api`. Production thật thường tách `/metrics` sang cổng nội bộ riêng.
