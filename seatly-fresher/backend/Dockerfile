# Bước 1: build ra 1 file binary bằng image có sẵn Go
FROM golang:1.27-alpine AS build
WORKDIR /src
# Copy go.mod/go.sum trước → Docker cache bước tải thư viện, sửa code không phải tải lại
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -trimpath -ldflags="-s -w": bỏ đường dẫn máy build + thông tin debug → binary nhỏ hơn
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /seatly ./cmd/api

# Bước 2: image chạy chỉ chứa file binary, không cần Go → rất nhẹ
FROM alpine:3.22
COPY --from=build /seatly /seatly
EXPOSE 8080
USER nobody
CMD ["/seatly"]
