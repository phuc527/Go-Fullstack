# Bankly Backend (Go)

API ngân hàng viết bằng Go. Thư viện ngoài: [pgx](https://github.com/jackc/pgx) (driver Postgres) và Prometheus client. Còn lại là thư viện chuẩn (`net/http`, `encoding/json`, `log/slog`, `embed`).

Cấu trúc **giống hệt Seatly** (`cmd/` + `internal/`, handler → service → repository). Nếu chưa đọc thì đọc `seatly-fresher/backend/README.md` trước.

## Chạy

```bash
docker compose up -d postgres     # ở thư mục bankly-fresher
go run ./cmd/api                  # → http://localhost:8090/api/accounts
DATABASE_URL=postgres://... go run ./cmd/api
go test -race ./...               # integration test cần TEST_DATABASE_URL, xem README ngoài
```

## Một lệnh chuyển tiền đi qua những đâu

```
POST /api/transfers  (header Idempotency-Key)
   │
   ▼
logRequests → measure → recoverPanic                middleware.go, metrics.go
   │
   ▼
handler.transfer                                    transaction_handler.go  ← đọc JSON + header, đổi lỗi → status
   │
   ▼
BankService.Transfer                                bank_service.go         ← validate, dựng 2 "chân" DEBIT/CREDIT
   │
   ▼
Repository.Post         BEGIN                       posting.go
   ├─ INSERT transactions ... ON CONFLICT (idempotency_key) DO NOTHING
   │     └─ trùng key → trả giao dịch cũ (replay)
   ├─ SELECT ... FROM accounts WHERE id = ANY(...) ORDER BY id FOR UPDATE   legs.go
   ├─ kiểm tra đóng băng / số dư
   ├─ UPDATE accounts SET balance = balance ± amount  (×2)
   ├─ INSERT ledger_entries                           (×2)
   └─ INSERT outbox_events
                        COMMIT  (lỗi ở bất kỳ bước nào → ROLLBACK hết)
```

## Từng file

```
backend/
├── cmd/api/main.go                      # nối dây: DB → repo → service → router, graceful shutdown
├── internal/
│   ├── config/config.go                 # PORT, DATABASE_URL
│   ├── domain/
│   │   ├── bank.go                      # Account, Entry, Transaction, Leg, Posting, hằng số
│   │   └── errors.go                    # lỗi nghiệp vụ
│   ├── service/
│   │   ├── bank_service.go              # luật: validate, dựng chân giao dịch, fingerprint
│   │   └── bank_service_test.go
│   ├── repository/postgres/
│   │   ├── schema.sql                   # bảng, ràng buộc, trigger append-only, quỹ ngân hàng
│   │   ├── db.go                        # mở pool, chạy schema
│   │   ├── account_repository.go        # tài khoản, sao kê, đọc giao dịch
│   │   ├── posting.go                   # ghi sổ: idempotency, đảo ngược, outbox
│   │   ├── legs.go                      # khóa tài khoản theo thứ tự, cập nhật số dư, ghi bút toán
│   │   └── *_test.go                    # integration test với Postgres thật
│   └── httpapi/
│       ├── router.go                    # route + /healthz (ping DB) + /metrics
│       ├── account_handler.go           # tài khoản, sao kê, đóng băng
│       ├── transaction_handler.go       # nạp, chuyển, đảo ngược; 201 vs 200 replay
│       ├── errors.go                    # BẢNG lỗi → HTTP status + nhãn metric + câu tiếng Việt
│       ├── response.go, middleware.go   # copy từ Seatly
│       ├── metrics.go
│       └── handler_test.go              # test HTTP với service giả
└── Dockerfile
```

### Vì sao số dư lưu ở 2 nơi?

`accounts.balance` là **cache** để đọc nhanh. Nguồn sự thật là `ledger_entries`. Hai nơi luôn được ghi trong cùng 1 transaction, và hàm `reconcile` trong test kiểm tra chúng khớp nhau. Ngân hàng thật chạy phép đối soát này hằng đêm.

### Vì sao `balance_after` nằm trong từng bút toán?

Để in sao kê mà không phải cộng dồn lại từ đầu, và để kiểm toán viên thấy ngay số dư tại từng thời điểm.

### Vì sao lưu `request_fingerprint`?

Client lỡ dùng lại 1 key cho request **khác** (vd đổi số tiền) thì không được trả kết quả cũ, vì đó là ý định khác. Service tóm tắt nội dung request thành chuỗi (`TRANSFER|2|3|300000`). Trùng key mà khác chuỗi thì trả `422`.

### Vì sao quỹ ngân hàng được âm?

Tiền nạp vào phải đến từ đâu đó. Kế toán kép bắt mọi giao dịch có 2 chân, nên nạp tiền = DEBIT quỹ, CREDIT khách. Quỹ âm 2 triệu nghĩa là khách đang gửi ngân hàng 2 triệu. Nhờ đó **tổng mọi số dư luôn = 0**, là phép kiểm tra đơn giản nhất rằng tiền không tự sinh ra.

## Metrics

| Metric                                 | Label                 | Dùng để                                                                |
| -------------------------------------- | --------------------- | ---------------------------------------------------------------------- |
| `bankly_http_requests_total`           | method, route, status | request/giây, tỉ lệ lỗi                                                |
| `bankly_http_request_duration_seconds` | method, route         | p95/p99                                                                |
| `bankly_postings_total`                | kind, result          | `posted`, `replayed`, `insufficient_funds`, `frozen`, `key_reused`...  |
| `bankly_posted_amount_vnd_total`       | kind                  | Tổng tiền ghi sổ. Là float, **chỉ để vẽ biểu đồ**, không dùng đối soát |
