# Bankly Fresher: hệ thống giao dịch ngân hàng (Go + React)

App ngân hàng mini **cấu trúc giống Seatly** (handler → service → repository), nhưng bài toán khó hơn: **tiền không được phép biến mất hay tự sinh ra**. Có mở tài khoản, nạp tiền, chuyển khoản, sao kê, đóng băng và đảo ngược giao dịch.

Khác Seatly ở 1 điểm lớn: dữ liệu nằm trong **PostgreSQL** ngay từ đầu, vì các kỹ thuật cốt lõi (khóa dòng, `UNIQUE`, transaction ACID, trigger chặn sửa sổ cái) chỉ có ý nghĩa với DB thật.

```
bankly-fresher/
├── backend/          # Go API: cmd/ + internal/ (handler → service → repository)
│                     #   → đọc backend/README.md để hiểu từng file
├── frontend/         # React + Vite + TypeScript + React Query
│   ├── src/api.ts             # gọi API, gắn header Idempotency-Key
│   ├── src/App.tsx            # bố cục: danh sách tài khoản | chi tiết
│   ├── src/AccountList.tsx    # danh sách tài khoản + quỹ ngân hàng
│   ├── src/AccountPanel.tsx   # số dư, đóng băng
│   ├── src/MoneyForm.tsx      # nạp tiền / chuyển khoản
│   ├── src/History.tsx        # sao kê + nút đảo ngược
│   ├── src/useIdempotencyKey.ts
│   └── nginx.conf / Dockerfile
├── monitoring/       # Prometheus + Grafana (dashboard có sẵn)
├── postgres/init.sql # tạo thêm DB bankly_test cho integration test
└── compose.yaml      # chạy tất cả bằng Docker
```

## Những ý tưởng chính (đọc code theo thứ tự này)

| Ý tưởng                                                                            | Ở đâu                                               |
| ---------------------------------------------------------------------------------- | --------------------------------------------------- |
| **Tiền là `BIGINT` (đồng)**, không bao giờ `float`                                 | `schema.sql`, `domain/bank.go`                      |
| **Kế toán kép**: mỗi giao dịch có 1 dòng DEBIT + 1 dòng CREDIT, tổng mọi số dư = 0 | `service/bank_service.go` (các "chân"), `legs.go`   |
| **Quỹ tiền mặt ngân hàng** (tài khoản #1, được âm) là nơi tiền nạp đi ra           | `schema.sql`, `AccountList.tsx`                     |
| **ACID**: trừ A + cộng B + ghi sổ + outbox nằm trong 1 transaction                 | `repository/postgres/posting.go`                    |
| **Khóa bi quan + thứ tự khóa** (`FOR UPDATE ... ORDER BY id`) → không deadlock     | `repository/postgres/legs.go`                       |
| **Idempotency-Key**: gửi lại cùng key → trả kết quả cũ, không trừ tiền lần 2       | `posting.go`, `useIdempotencyKey.ts`                |
| **Sổ cái bất biến**: trigger trong DB chặn `UPDATE`/`DELETE`/`TRUNCATE`            | `schema.sql`                                        |
| **Bút toán đảo ngược** thay cho sửa/xóa; `reverses_id UNIQUE` → chỉ đảo 1 lần      | `posting.go` (`resolveReversal`)                    |
| **Transactional Outbox**: ghi sự kiện vào `outbox_events` cùng transaction         | `posting.go`                                        |
| **Đối soát**: tổng số dư = 0, số dư cache = cộng dồn sổ cái (chạy sau mỗi test)    | `repository/postgres/helpers_test.go` (`reconcile`) |

## Cách 1: chạy khi code hằng ngày

```bash
# Terminal 1: chỉ bật Postgres bằng Docker
cd bankly-fresher
docker compose up -d postgres
```

```bash
# Terminal 2: backend (tự tạo bảng khi khởi động)
cd bankly-fresher/backend
go run ./cmd/api   # → http://localhost:8090/api/accounts
```

```bash
# Terminal 3: frontend
cd bankly-fresher/frontend
npm install   # chỉ cần lần đầu
npm run dev   # → mở http://localhost:5173
```

## Cách 2: chạy tất cả bằng Docker (cả monitoring)

```bash
cd bankly-fresher
docker compose up --build -d
docker compose logs -f backend   # xem log JSON
docker compose down              # tắt (giữ dữ liệu)
docker compose down -v           # tắt và XÓA dữ liệu Postgres
```

Cổng khác Seatly nên chạy song song được:

| Service      | Mở ở                  | Ghi chú                                                          |
| ------------ | --------------------- | ---------------------------------------------------------------- |
| `frontend`   | http://localhost:3100 | nginx phục vụ React, chuyển `/api` sang `backend`                |
| `backend`    | http://localhost:8090 | Chỉ chạy khi Postgres sẵn sàng; FE chỉ chạy khi `/healthz` OK    |
| `postgres`   | `localhost:5433`      | user/pass/db: `bankly`. Có thêm DB `bankly_test` cho test        |
| `prometheus` | http://localhost:9091 | Cứ 5s kéo `backend:8090/metrics`                                 |
| `grafana`    | http://localhost:3101 | Dashboard **Bankly Backend**. Đăng nhập admin/admin nếu muốn sửa |

Xem thẳng dữ liệu trong DB:

```bash
docker compose exec postgres psql -U bankly -c "SELECT * FROM ledger_entries ORDER BY id DESC LIMIT 10"
docker compose exec postgres psql -U bankly -c "SELECT SUM(balance) FROM accounts"   # luôn = 0
docker compose exec postgres psql -U bankly -c "DELETE FROM ledger_entries"          # bị trigger từ chối
```

## API

Mọi response có dạng `{ "success": bool, "data": ..., "message": "..." }`, HTTP status vẫn đúng chuẩn.

| Method | URL                                  | Body                                             | Kết quả                                            |
| ------ | ------------------------------------ | ------------------------------------------------ | -------------------------------------------------- |
| GET    | `/api/accounts`                      |                                                  | `200` danh sách tài khoản (kể cả quỹ ngân hàng #1) |
| POST   | `/api/accounts`                      | `{"owner":"An"}`                                 | `201`                                              |
| GET    | `/api/accounts/{id}/entries`         |                                                  | `200` sao kê 100 dòng mới nhất · `404`             |
| PUT    | `/api/accounts/{id}/status`          | `{"status":"FROZEN"}`                            | `200` · `404`                                      |
| POST   | `/api/deposits` 🔑                   | `{"accountId":2,"amount":500000}`                | xem bảng dưới                                      |
| POST   | `/api/transfers` 🔑                  | `{"fromAccountId":2,"toAccountId":3,"amount":1}` | xem bảng dưới                                      |
| POST   | `/api/transactions/{id}/reversal` 🔑 |                                                  | xem bảng dưới · `409` đã đảo ngược rồi             |
| GET    | `/healthz`                           |                                                  | `200` khi DB còn sống, `503` khi không             |

🔑 = bắt buộc header `Idempotency-Key`. Kết quả:

| Status | Khi nào                                                                                |
| ------ | -------------------------------------------------------------------------------------- |
| `201`  | Giao dịch mới được ghi sổ                                                              |
| `200`  | Key đã dùng với **cùng** nội dung → trả lại giao dịch cũ, header `Idempotent-Replayed` |
| `400`  | Thiếu key, số tiền ≤ 0, tự chuyển cho mình, dùng quỹ ngân hàng                         |
| `404`  | Không có tài khoản / giao dịch                                                         |
| `409`  | Tài khoản đóng băng, giao dịch đã đảo ngược                                            |
| `422`  | Không đủ số dư · key đã dùng cho request **khác** nội dung                             |

### Thử Idempotency bằng tay (Git Bash)

```bash
curl -X POST localhost:8090/api/accounts -d '{"owner":"An"}'     # → id 2
curl -X POST localhost:8090/api/accounts -d '{"owner":"Bình"}'   # → id 3
curl -X POST localhost:8090/api/deposits -H "Idempotency-Key: nap-1" -d '{"accountId":2,"amount":1000000}'

# Gửi 3 lần cùng key: lần 1 → 201, lần 2-3 → 200 + Idempotent-Replayed, chỉ trừ tiền 1 lần
for i in 1 2 3; do
  curl -si -X POST localhost:8090/api/transfers -H "Idempotency-Key: ck-1" \
    -d '{"fromAccountId":2,"toAccountId":3,"amount":300000}' | grep -E "^HTTP|Replayed"
done
```

Trên giao diện, dòng `Idempotency-Key: ...` dưới nút "Chuyển tiền" cho thấy key hiện tại: bấm lại sau lỗi thì key giữ nguyên, giao dịch thành công hoặc sửa số tiền thì key đổi.

### Thử monitoring

1. Mở http://localhost:3100, nạp/chuyển vài lần, cố tình chuyển quá số dư, đảo ngược 1 giao dịch.
2. Mở Grafana http://localhost:3101: request/giây, p95, giao dịch theo kết quả (`posted`, `replayed`, `insufficient_funds`...), số request lặp lại bị chặn.
3. PromQL ở http://localhost:9091: `sum by (result) (bankly_postings_total)`.

## Test

```bash
cd bankly-fresher/backend
go test -race ./...    # unit test (không cần DB); integration test tự skip

# Có Postgres (docker compose up -d postgres) → chạy cả integration test:
TEST_DATABASE_URL="postgres://bankly:bankly@localhost:5433/bankly_test?sslmode=disable" go test -race ./...
```

Integration test chứng minh các yêu cầu phi chức năng với DB thật:

| Test                                | Chứng minh                                                       |
| ----------------------------------- | ---------------------------------------------------------------- |
| `TestIdempotencyUnderConcurrency`   | 20 request cùng key cùng lúc → ghi sổ đúng 1 lần                 |
| `TestConcurrentTransfersNoDeadlock` | 100 lệnh A→B và B→A song song → không deadlock, không mất tiền   |
| `TestNoOverdraftUnderConcurrency`   | 1.000đ, 10 người cùng rút 300đ → đúng 3 người thành công         |
| `TestReversal`                      | Đảo ngược trả lại số dư; không đảo 2 lần; không đảo bút toán đảo |
| `TestLedgerIsAppendOnly`            | `UPDATE`/`DELETE`/`TRUNCATE` sổ cái bị DB từ chối                |
| `reconcile` (sau mọi test)          | Tổng số dư = 0, số dư cache khớp sổ cái                          |

Bài thử: trong `legs.go`, đổi sang khóa từng tài khoản theo thứ tự "nguồn rồi đích" (`SELECT ... WHERE id = $1 FOR UPDATE` trong vòng lặp), chạy lại test. `TestConcurrentTransfersNoDeadlock` sẽ fail với `deadlock detected (SQLSTATE 40P01)`.

## Lint

```bash
cd bankly-fresher/frontend && npm run lint
cd bankly-fresher/backend  && gofmt -l . && go vet ./...
```

## Chưa làm (so với bài thiết kế) và các level

| Level     | Nâng cấp                                                                                                                    |
| --------- | --------------------------------------------------------------------------------------------------------------------------- |
| **L1** ✅ | Kế toán kép · ACID · khóa bi quan có thứ tự · Idempotency-Key · sổ cái bất biến · đảo ngược · outbox (chỉ ghi) · Prometheus |
| L2        | Tiến trình `cmd/relay` đọc `outbox_events` chưa gửi (`FOR UPDATE SKIP LOCKED`), log/gửi đi rồi đánh dấu `published_at`      |
| L3        | Lưu cả giao dịch **thất bại** (status `FAILED`) để kiểm toán; hiện tại lỗi = ROLLBACK, key chưa bị tiêu                     |
| L4        | Migration thật (goose) thay cho `schema.sql` chạy lúc khởi động; phân trang sao kê                                          |
| L5        | Đăng nhập JWT: chỉ chủ tài khoản được chuyển tiền; vai trò "kiểm soát viên" mới được đảo ngược/đóng băng                    |
| L6        | Khóa lạc quan (`version`) so sánh với khóa bi quan bằng load test (k6); tài khoản "nóng" chia sub-account                   |
| L7        | Kafka + Saga cho chuyển liên ngân hàng, CQRS read-replica cho sao kê, alert Prometheus khi `reconcile` lệch                 |

Tham khảo: [techschool/simplebank](https://github.com/techschool/simplebank) (Go + Postgres, xử lý deadlock khi chuyển tiền), [formancehq/ledger](https://github.com/formancehq/ledger) (ledger kế toán kép production).
