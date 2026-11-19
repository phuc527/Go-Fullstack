-- Postgres chỉ chạy file này 1 lần, lúc volume còn trống (lần đầu `docker compose up`).
-- DB riêng cho integration test: sổ cái không xóa được, nên tách hẳn khỏi dữ liệu của app.
CREATE DATABASE bankly_test OWNER bankly;
