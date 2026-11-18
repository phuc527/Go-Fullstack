# Học Go Backend từ số 0 — cho dev React/Node

> Tham chiếu [roadmap.sh/golang](https://roadmap.sh/golang) + [roadmap.sh/backend](https://roadmap.sh/backend).
> Bạn đã biết lập trình (JS/Node) → chỉ cần học **cách viết lại những thứ bạn đã biết bằng Go**. Đi chậm, chắc.

## Project xuyên suốt: Seatly — đặt vé sự kiện

Mọi bài tập đều xoay quanh 1 domain: **sự kiện (event) → ghế (seat) → đặt vé (booking)**.

- Tuần 1–4: bài tập nhỏ dùng dữ liệu vé/ghế để học cú pháp.
- Tuần 5 trở đi: ghép lại thành **Seatly API** thật, lớn dần mỗi ngày.

Bài toán cốt lõi của đặt vé — **2 người không được đặt trùng 1 ghế** — sẽ xuất hiện nhiều lần, mỗi lần ở mức khó hơn (bộ nhớ → mutex → database transaction).

## Cách học

- **1–1.5h/ngày, 5 ngày/tuần.** Mỗi ngày **1 khái niệm**, không hơn.
- Mỗi ngày có thư mục `labs/dayXX/` gồm:
  - `README.md`: giải thích ngắn + **so sánh với JS** + bài tập + **đáp án giấu sẵn** (bấm mở khi bí)
  - `main.go`: file khởi đầu, bạn điền code vào chỗ `TODO`
- Chạy: `go run ./labs/day01` → so output với **"Kết quả mong đợi"** trong README. Khớp là xong.
- **Ngày thứ 5 mỗi tuần = ôn tập**: 1 bài nhỏ gom kiến thức cả tuần. Không có kiến thức mới.
- Bí quá 20 phút → mở đáp án, đọc hiểu, **xóa đi rồi tự gõ lại**. Không sao cả.
- Viết test (`go test`) bắt đầu từ Tuần 4, khi đã quen cú pháp.

---

## Giai đoạn 1 — Làm quen Go (Tuần 1–4)

_roadmap.sh: Learn the Basics_

### Tuần 1 — Hello Go

| Day | Học gì                                                   | Bài tập                                          |
| --- | -------------------------------------------------------- | ------------------------------------------------ |
| 1   | Chạy chương trình đầu tiên, `fmt.Println`, `fmt.Printf`  | In thông tin 1 vé: tên sự kiện, ghế, giá         |
| 2   | Biến (`:=`, `var`), kiểu `int`/`float64`/`string`/`bool` | Tính tổng tiền 3 vé + phí dịch vụ 5%             |
| 3   | `if/else`, `for` (Go chỉ có `for`), `switch`             | Giá vé theo hạng (VIP/thường); in sơ đồ hàng ghế |
| 4   | Hàm, trả **nhiều giá trị** cùng lúc                      | `tinhTien(soVe, hang)` trả `(tong, giamGia)`     |
| 5   | **Ôn tập**                                               | Máy bán vé terminal: nhập số vé → in hóa đơn     |

### Tuần 2 — Danh sách & từ điển

| Day | Học gì                                                    | Bài tập                                          |
| --- | --------------------------------------------------------- | ------------------------------------------------ |
| 6   | Slice (giống array JS): tạo, `append`, `len`, `for range` | Danh sách ghế `["A1","A2",...]`, thêm/đếm ghế    |
| 7   | Cắt slice `s[1:3]`, package `slices` (sort, contains)     | Kiểm tra ghế có tồn tại, sắp xếp ghế             |
| 8   | Map (giống object/`Map` JS): thêm, đọc, xóa, kiểm tra key | `map[ghế]trạngThái`: đặt ghế, xem ghế còn trống  |
| 9   | String: package `strings`                                 | Tách `"A12"` → hàng `A`, số `12`                 |
| 10  | **Ôn tập**                                                | Sơ đồ phòng chiếu trên terminal, chọn ghế để đặt |

### Tuần 3 — Struct, method, error

| Day | Học gì                                                        | Bài tập                                             |
| --- | ------------------------------------------------------------- | --------------------------------------------------- |
| 11  | Struct (giống TS `type Event = {...}`)                        | Struct `Event`, `Seat`, in danh sách sự kiện        |
| 12  | Method + pointer (quy tắc: muốn **sửa** dữ liệu thì dùng `*`) | `event.Book("A1")` đánh dấu ghế đã đặt              |
| 13  | Error: `if err != nil` thay cho `try/catch`                   | `Book` báo lỗi: ghế không tồn tại / đã có người đặt |
| 14  | Package & module: tách file, **chữ Hoa = export**             | Tách code ra package `booking`                      |
| 15  | **Ôn tập**                                                    | CLI đặt vé: xem sự kiện → chọn ghế → đặt → hủy      |

### Tuần 4 — Interface, JSON, file, test

| Day | Học gì                                                 | Bài tập                                         |
| --- | ------------------------------------------------------ | ----------------------------------------------- |
| 16  | Interface (giống TS interface, không cần `implements`) | `PriceRule` cho giá thường / giảm giá sinh viên |
| 17  | JSON: struct ↔ JSON (như `JSON.stringify/parse`)       | In `Event` ra JSON                              |
| 18  | Đọc/ghi file                                           | CLI đặt vé lưu vào `bookings.json`              |
| 19  | Unit test đầu tiên với `go test` (như Jest)            | Test: đặt ghế trống ok, đặt ghế đã đặt → lỗi    |
| 20  | **Ôn tập**                                             | CLI đặt vé hoàn chỉnh + test                    |

## Giai đoạn 2 — Seatly API (Tuần 5–6)

_roadmap.sh: Web Frameworks, APIs_ — Bê CLI lên thành REST API, giống làm với Express.

### Tuần 5 — HTTP server

| Day | Học gì                                         | Bài tập                                     |
| --- | ---------------------------------------------- | ------------------------------------------- |
| 21  | Server đầu tiên với `net/http` (như `app.get`) | `GET /health`                               |
| 22  | Trả JSON                                       | `GET /events`                               |
| 23  | Đọc body JSON                                  | `POST /events` tạo sự kiện                  |
| 24  | Path param `/events/{id}`                      | `GET /events/{id}/seats` xem sơ đồ ghế      |
| 25  | **Ôn tập**: status code, lỗi; Postman          | `POST /bookings` đặt ghế (lưu trong bộ nhớ) |

### Tuần 6 — Middleware & goroutine nhập môn

| Day | Học gì                                        | Bài tập                                               |
| --- | --------------------------------------------- | ----------------------------------------------------- |
| 26  | Middleware (như Express middleware)           | Middleware log mỗi request                            |
| 27  | Goroutine: chạy song song bằng chữ `go`       | Giả lập 100 người cùng bấm đặt vé                     |
| 28  | `sync.Mutex` — **vấn đề đặt trùng ghế lần 1** | Thấy bug 2 người cùng đặt được 1 ghế → sửa bằng mutex |
| 29  | Channel cơ bản                                | Hàng đợi gửi "email xác nhận" sau khi đặt vé          |
| 30  | **Ôn tập**                                    | Seatly API + test bằng `httptest`                     |

## Giai đoạn 3 — Database & Auth (Tuần 7–9)

_roadmap.sh: Relational DB, Authentication, Web Security_

### Tuần 7 — PostgreSQL

| Day | Học gì                                     | Bài tập                                     |
| --- | ------------------------------------------ | ------------------------------------------- |
| 31  | Chạy Postgres bằng Docker, ôn SQL          | Thiết kế bảng `events`, `seats`, `bookings` |
| 32  | Kết nối DB bằng `pgx`, query đơn giản      | `GET /events` đọc từ DB                     |
| 33  | Insert / update                            | Tạo sự kiện + sinh ghế vào DB               |
| 34  | Migration với `goose` (như Prisma migrate) | Thêm cột `price` cho ghế bằng migration     |
| 35  | **Ôn tập**                                 | Seatly chạy hoàn toàn trên Postgres         |

### Tuần 8 — Cấu trúc project & transaction

| Day | Học gì                                                                              | Bài tập                                                      |
| --- | ----------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| 36  | Config từ `.env`                                                                    | Đọc DB URL, port từ env                                      |
| 37  | Chia layer: handler → service → repository                                          | Refactor Seatly                                              |
| 38  | Validate input, format lỗi thống nhất                                               | Đặt ghế không tồn tại → 404, thiếu field → 400               |
| 39  | Transaction — **vấn đề đặt trùng ghế lần 2** (mutex không đủ khi chạy nhiều server) | Đặt ghế trong transaction, 2 request cùng ghế → 1 thành công |
| 40  | **Ôn tập**                                                                          | Phân trang `GET /events?page=2`                              |

### Tuần 9 — Đăng nhập

| Day | Học gì                       | Bài tập                                             |
| --- | ---------------------------- | --------------------------------------------------- |
| 41  | Hash password với `bcrypt`   | `POST /signup`                                      |
| 42  | JWT                          | `POST /login` trả token                             |
| 43  | Middleware xác thực          | Phải đăng nhập mới đặt vé được                      |
| 44  | Phân quyền: user / organizer | Chỉ organizer tạo sự kiện; user chỉ xem vé của mình |
| 45  | **Ôn tập**                   | Postman: signup → login → đặt vé → xem vé của tôi   |

## Giai đoạn 4 — Đưa lên production (Tuần 10–11)

_roadmap.sh: Testing, Logging, Containerization, CI/CD_

| Day | Học gì                                          | Bài tập                                 |
| --- | ----------------------------------------------- | --------------------------------------- |
| 46  | Table-driven test (cách viết test chuẩn của Go) | Test service đặt vé                     |
| 47  | Test handler với `httptest`                     | Test API không cần chạy server          |
| 48  | Log có cấu trúc với `log/slog`                  | Log JSON mỗi lượt đặt vé                |
| 49  | Graceful shutdown                               | Tắt server không mất lượt đặt đang chạy |
| 50  | **Ôn tập** + `golangci-lint`                    | Sửa hết cảnh báo lint                   |
| 51  | Dockerfile                                      | Đóng gói Seatly thành image             |
| 52  | Docker Compose: API + Postgres                  | `docker compose up` là chạy             |
| 53  | GitHub Actions                                  | Tự chạy test mỗi lần push               |
| 54  | Deploy (Render / Fly.io)                        | Có link API public                      |
| 55  | **Ôn tập**                                      | README giới thiệu Seatly                |

## Giai đoạn 5 — Nâng cao (khi đã vững, sẽ chi tiết sau)

Giữ ghế 10 phút rồi tự nhả (Redis) · thanh toán + webhook không bị trừ tiền 2 lần (idempotency) · gửi email bằng background job · rate limit · sơ đồ ghế realtime (WebSocket) · load test 1000 người cùng giành vé · `pprof` · gRPC.

---

## Tài liệu (đọc kèm, không bắt buộc)

- [Tour of Go](https://go.dev/tour) — tương tác ngay trên web, rất hợp người mới
- [Go by Example](https://gobyexample.com) — mỗi trang 1 ví dụ ngắn, dùng để tra cứu
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests) — từ Tuần 4
