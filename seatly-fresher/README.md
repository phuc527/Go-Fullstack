# Seatly Fresher — App đặt vé (Go + React)

Phiên bản **đơn giản nhất có thể chạy được**: 1 sự kiện, 15 ghế, đặt từng ghế. Dữ liệu lưu trong RAM nên tắt BE là mất. Mỗi tuần học xong thì nâng cấp thêm 1 bước (xem [Lộ trình nâng cấp](#lộ-trình-nâng-cấp)).

```
seatly-fresher/
├── be/              # Go API (net/http, không dùng framework)
│   ├── main.go
│   ├── main_test.go
│   └── Dockerfile
├── fe/              # React + Vite + TypeScript + React Query
│   ├── src/api.ts   # gọi API
│   ├── src/App.tsx  # sơ đồ ghế
│   ├── nginx.conf   # chỉ dùng khi chạy Docker
│   └── Dockerfile
└── compose.yaml     # chạy cả 2 bằng Docker
```

## Cách 1 — Chạy khi code hằng ngày (nhanh nhất)

Mở 2 terminal:

```bash
# Terminal 1 — BE
cd seatly-fresher/be
go run .
# → 🎫 Seatly BE chạy ở http://localhost:8080
```

```bash
# Terminal 2 — FE
cd seatly-fresher/fe
npm install   # chỉ lần đầu
npm run dev
# → mở http://localhost:5173
```

FE gọi `/api/...`, Vite tự chuyển sang `localhost:8080` (xem `vite.config.ts`), nên **không cần cấu hình CORS**.

Sửa code Go xong thì **tắt BE (Ctrl+C) rồi `go run .` lại**. Go không tự reload như nodemon. Sau này sẽ dùng `air` để tự reload.

## Cách 2 — Chạy bằng Docker

Cần mở Docker Desktop trước.

```bash
cd seatly-fresher
docker compose up --build     # build + chạy cả BE và FE
# → mở http://localhost:3000   (API vẫn gọi thẳng được ở http://localhost:8080)

docker compose down           # tắt
```

Sửa code xong thì chạy lại `docker compose up --build`.

| Service | Image                                           | Port                        |
| ------- | ----------------------------------------------- | --------------------------- |
| `be`    | Go build ra 1 file binary, chạy trên Alpine     | `8080`                      |
| `fe`    | Build React ra file tĩnh, nginx phục vụ         | `3000` → `/api` chuyển sang `be` |

## API

| Method | URL             | Body                 | Kết quả                                                  |
| ------ | --------------- | -------------------- | -------------------------------------------------------- |
| GET    | `/api/seats`    |                      | `200` danh sách ghế                                      |
| POST   | `/api/bookings` | `{"seatId": "A1"}`   | `201` đặt được · `404` ghế không có · `409` đã có người đặt · `400` body sai |

Thử bằng curl (Git Bash):

```bash
curl localhost:8080/api/seats
curl -X POST localhost:8080/api/bookings -d '{"seatId":"A1"}'
curl -X POST localhost:8080/api/bookings -d '{"seatId":"A1"}'   # lần 2 → 409
```

## Test BE

```bash
cd seatly-fresher/be
go test -v .
```

Có 2 test:

- Đặt A1 → 201, đặt lại → 409, Z9 → 404, body sai → 400.
- **100 người cùng lúc giành ghế B1 → chỉ đúng 1 người được.** Thử xóa 2 dòng `mu.Lock()` / `defer mu.Unlock()` trong `createBooking` rồi chạy lại test. Thường bạn sẽ thấy test fail. Đây chính là bài toán "đặt trùng ghế" sẽ học ở Day 28.

## Đọc code BE: những thứ chưa học

`be/main.go` có dùng vài thứ chưa học tới. Đọc chú thích trong code là đủ hiểu để chạy, còn học kỹ thì theo ngày dưới đây:

| Trong code                              | Là gì                                        | Học kỹ ở  |
| --------------------------------------- | -------------------------------------------- | --------- |
| `type Seat struct {...}`                | Giống TS `type Seat = {...}`                 | Day 11    |
| `` `json:"id"` ``                       | Tên field khi đổi sang JSON                  | Day 17    |
| `if err := ...; err != nil`             | Xử lý lỗi                                    | Day 13    |
| `func(w http.ResponseWriter, r *http.Request)` | Handler, giống `(req, res) => {}` của Express | Day 21 |
| `mux.HandleFunc("GET /api/seats", ...)` | Giống `app.get("/api/seats", ...)`           | Day 21    |
| `sync.Mutex`, `mu.Lock()`               | Khóa để 2 request không sửa cùng lúc         | Day 28    |
| `defer`                                 | Chạy khi hàm kết thúc, giống `finally`       | Day 28    |

## Lộ trình nâng cấp

Mỗi bước gắn với các ngày trong [ROADMAP.md](../ROADMAP.md). Học tới đâu thì nâng cấp tới đó.

| Phiên bản  | Sau khi học | Nâng cấp                                                                     |
| ---------- | ----------- | ---------------------------------------------------------------------------- |
| **v0** ✅  | Day 10      | Bản hiện tại: 1 sự kiện, đặt ghế, lưu trong RAM                              |
| v1         | Day 11–14   | Tách code vào package `booking`: `Event` có method `Book(seatID) error`      |
| v2         | Day 17–18   | Lưu dữ liệu ra file JSON → tắt BE không mất ghế đã đặt                       |
| v3         | Day 19–20   | Test viết theo kiểu table-driven                                             |
| v4         | Day 21–25   | Nhiều sự kiện: `GET /api/events`, `GET /api/events/{id}/seats`; FE có trang chọn sự kiện |
| v5         | Day 26–30   | Middleware ghi log và bắt panic; hiểu thật kỹ mutex                          |
| v6         | Day 31–35   | Thêm PostgreSQL vào `compose.yaml`, BE đọc/ghi DB                            |
| v7         | Day 36–40   | Chia layer handler → service → repo; đặt ghế bằng transaction                |
| v8         | Day 41–45   | Đăng ký/đăng nhập JWT, "vé của tôi"; FE có trang login                       |
| v9         | Day 46–55   | Log có cấu trúc, CI chạy test, deploy lên internet                           |
