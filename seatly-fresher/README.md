# Seatly Fresher: app đặt vé (Go + React)

App đặt vé **cấu trúc giống production nhưng vừa sức fresher**. Hiện có 1 sự kiện, 15 ghế, đặt từng ghế, dữ liệu lưu trong RAM. Học tới đâu thì nâng level tới đó (xem [Các level](#các-level)).

```
seatly-fresher/
├── backend/         # Go API: cmd/ + internal/ (handler → service → repository)
│                    #   → đọc backend/README.md để hiểu từng file
├── frontend/        # React + Vite + TypeScript + React Query
│   ├── src/api.ts   # gọi API
│   ├── src/App.tsx      # trang đặt ghế: query, mutation, trạng thái tải/lỗi
│   ├── src/SeatMap.tsx  # sơ đồ ghế + chú thích
│   ├── src/format.ts    # định dạng giá tiền
│   ├── nginx.conf   # chỉ dùng khi chạy Docker
│   └── Dockerfile
├── monitoring/      # Prometheus + Grafana (dashboard có sẵn)
│   ├── prometheus.yml
│   └── grafana/
│       ├── provisioning/   # tự thêm datasource + nạp dashboard khi Grafana khởi động
│       └── dashboards/seatly.json
└── compose.yaml     # chạy tất cả bằng Docker
```

## Sơ đồ

Mở bằng trình duyệt (file HTML tự chạy, không cần server):

- [docs/architecture.html](docs/architecture.html): kiến trúc tổng thể, gồm Docker, 3 layer backend và Prometheus/Grafana
- [docs/booking-sequence.html](docs/booking-sequence.html): đường đi của 1 request đặt ghế, lần 1 được 201, lần 2 bị 409

Sửa sơ đồ: chỉnh file `.json` cùng tên, rồi bảo Claude chạy lại `/archify`.

## Cách 1: chạy khi code hằng ngày (nhanh nhất)

```bash
# Terminal 1: backend
cd seatly-fresher/backend
go run ./cmd/api
```

```bash
# Terminal 2: frontend
cd seatly-fresher/frontend
npm install   # chỉ cần lần đầu
npm run dev   # → mở http://localhost:5173
```

Frontend gọi `/api/...`, Vite tự chuyển các request này sang `localhost:8080` (xem `vite.config.ts`), nên **không cần cấu hình CORS**.

Go không tự reload như nodemon: sửa code xong phải Ctrl+C rồi chạy lại `go run ./cmd/api`.

## Cách 2: chạy bằng Docker (cả monitoring)

Cần mở Docker Desktop trước.

```bash
cd seatly-fresher
docker compose up --build -d
docker compose logs -f backend   # xem log JSON của backend
docker compose down              # tắt
```

| Service      | Mở ở                  | Ghi chú                                                                     |
| ------------ | --------------------- | --------------------------------------------------------------------------- |
| `frontend`   | http://localhost:3000 | nginx phục vụ React, chuyển `/api` sang `backend`                           |
| `backend`    | http://localhost:8080 | Chỉ chạy `frontend` khi `/healthz` của backend đã OK                        |
| `prometheus` | http://localhost:9090 | Cứ 5s kéo `backend:8080/metrics` về 1 lần. Xem ở Status → Targets           |
| `grafana`    | http://localhost:3001 | Mở là thấy dashboard **Seatly Backend**. Đăng nhập admin/admin nếu muốn sửa |

### Thử monitoring

1. Mở http://localhost:3000 và đặt vài ghế, cố tình đặt trùng.
2. Mở Grafana tại http://localhost:3001. Sau khoảng 10 giây sẽ thấy request/giây, độ trễ p95, số lượt đặt theo kết quả (`booked`/`conflict`/`not_found`), goroutine và RAM.
3. Muốn có nhiều số liệu hơn thì bắn tải từ Git Bash:

   ```bash
   for i in $(seq 1 200); do curl -s localhost:8080/api/seats > /dev/null; done
   ```

4. Thử viết truy vấn PromQL ở http://localhost:9090, ví dụ `rate(seatly_http_requests_total[1m])`.

## API

Mọi response đều có cùng dạng. HTTP status vẫn đúng chuẩn (200/201/400/404/409/500):

```jsonc
// thành công
{ "success": true, "data": { "id": "A1", "booked": true }, "message": "Đặt ghế A1 thành công" }
// lỗi
{ "success": false, "data": null, "message": "Ghế A1 đã có người đặt" }
```

| Method | URL             | Body               | Kết quả                                                                      |
| ------ | --------------- | ------------------ | ---------------------------------------------------------------------------- |
| GET    | `/api/seats`    |                    | `200` danh sách ghế                                                          |
| POST   | `/api/bookings` | `{"seatId": "A1"}` | `201` đặt được · `404` ghế không có · `409` đã có người đặt · `400` body sai |
| GET    | `/healthz`      |                    | `200 {"status":"ok"}`                                                        |
| GET    | `/metrics`      |                    | Số liệu cho Prometheus                                                       |

```bash
curl localhost:8080/api/seats
curl -X POST localhost:8080/api/bookings -d '{"seatId":"A1"}'
curl -X POST localhost:8080/api/bookings -d '{"seatId":"A1"}'   # lần 2 → 409
```

## Test backend

```bash
cd seatly-fresher/backend
go test -race ./...
```

`-race` bật **race detector**: Go theo dõi xem có 2 goroutine nào cùng đọc/ghi một dữ liệu mà không khóa hay không. Trên Windows, cờ này cần `gcc` (máy đã cài qua WinLibs).

Bài thử: xóa 2 dòng `r.mu.Lock()` / `defer r.mu.Unlock()` trong `Book` của `internal/repository/memory/seat_repository.go`, rồi chạy lại test. Bạn sẽ thấy `WARNING: DATA RACE` và test 100 người giành 1 ghế bị fail. Đây là bài toán "đặt trùng ghế" sẽ học kỹ ở Day 28.

## Lint (kiểm tra code tự động)

```bash
cd seatly-fresher/frontend && npm run lint       # oxlint (nhanh hơn ESLint ~50 lần)
cd seatly-fresher/backend  && gofmt -l . && go vet ./...
```

Lint được siết dần theo level, để lúc mới học không bị ngợp:

| Level              | Frontend (`.oxlintrc.json`)                                                                                                                                    | Backend                                                                    |
| ------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| **Fresher** ✅     | Nhóm `correctness` báo **error** (code gần như chắc chắn sai). Nhóm `suspicious` chỉ **warn**. Thêm rule hooks, `jsx-key`, `eqeqeq` (`===`), a11y (`jsx-a11y`) | `gofmt` (format chuẩn, Go không tranh cãi style) + `go vet` (lỗi hay gặp)  |
| Junior (L4)        | Bật thêm `perf`, chặn warning: `oxlint --deny-warnings`                                                                                                        | Thêm `golangci-lint` gồm `errcheck` (quên xử lý error), `staticcheck`      |
| Middle (L8, có CI) | CI chạy lint + `tsc` + test trên mỗi PR, fail thì không được merge. Bật lint có dùng thông tin kiểu TypeScript (type-aware)                                    | `golangci-lint` bật thêm `gosec` (bảo mật), `revive`; CI fail khi lint lỗi |

## Đọc code backend: những thứ chưa học

| Trong code                                              | Là gì                                              | Học kỹ ở |
| ------------------------------------------------------- | -------------------------------------------------- | -------- |
| `type Seat struct {...}`                                | Giống TS `type Seat = {...}`                       | Day 11   |
| `func (r *SeatRepository) Book(...)`                    | Method gắn vào struct, giống method của class      | Day 12   |
| `errors.New`, `errors.Is`                               | Lỗi là giá trị, so sánh để biết lỗi gì             | Day 13   |
| `package`, `import "seatly/internal/..."`               | Chia code thành package                            | Day 14   |
| `type SeatRepository interface {...}`                   | Giống TS `interface`, nhưng không cần `implements` | Day 16   |
| `` `json:"id"` ``                                       | Tên field khi đổi sang JSON                        | Day 17   |
| `t.Run`, bảng `tests := []struct{...}`                  | Table-driven test                                  | Day 19   |
| `http.Handler`, `mux.HandleFunc("GET /api/seats", ...)` | Handler và route, giống Express                    | Day 21   |
| `func(next http.Handler) http.Handler`                  | Middleware, giống `app.use`                        | Day 26   |
| `context.Context`, `r.Context()`                        | Hủy/timeout cho request, giống `AbortSignal`       | Day 27   |
| `sync.Mutex`, `defer`, `go func()`, `chan`              | Khóa, chạy khi hàm kết thúc, goroutine, channel    | Day 28   |
| `signal.NotifyContext`, `srv.Shutdown`                  | Graceful shutdown                                  | Day 46   |
| `slog`, Prometheus                                      | Log có cấu trúc, metrics                           | Day 47   |

## Các level

Mỗi level gắn với các ngày trong [ROADMAP.md](../ROADMAP.md). Cấu trúc layer đã có sẵn từ đầu, nên mỗi level chỉ **thêm vào đúng chỗ** chứ không phải đập đi viết lại.

| Level     | Sau khi học | Nâng cấp                                                                                                                                                           |
| --------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **L1** ✅ | (hiện tại)  | Layer handler → service → repository · config bằng env · log JSON (`slog`) · graceful shutdown · `/healthz` · Prometheus + Grafana · table-driven test + race test |
| L2        | Day 17–18   | Thêm `repository/file`: lưu ghế ra file JSON. Đổi 1 dòng trong `main.go` là chạy được, tận mắt thấy lợi ích của interface                                          |
| L3        | Day 21–25   | Nhiều sự kiện: `GET /api/events`, `GET /api/events/{id}/seats`. Frontend có trang chọn sự kiện. Viết spec OpenAPI                                                  |
| L4        | Day 26–30   | Middleware request ID (gắn vào log), timeout bằng `context`, giới hạn số request (rate limit)                                                                      |
| L5        | Day 31–35   | Thêm `repository/postgres` + migration, Postgres trong `compose.yaml`, panel số kết nối DB trên Grafana                                                            |
| L6        | Day 36–40   | Đặt nhiều ghế trong 1 transaction, `UPDATE ... WHERE booked = false`, integration test với DB thật                                                                 |
| L7        | Day 41–45   | Đăng ký/đăng nhập JWT, trang "vé của tôi"; frontend có trang login                                                                                                 |
| L8        | Day 46–55   | CI bằng GitHub Actions, cảnh báo Prometheus (alert rules), tracing OpenTelemetry, deploy lên internet                                                              |
