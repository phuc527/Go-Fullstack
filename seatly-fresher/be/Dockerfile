# Bước 1: build ra 1 file binary bằng image có sẵn Go
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /seatly .

# Bước 2: image chạy chỉ chứa file binary, không cần Go → rất nhẹ
FROM alpine:3.22
COPY --from=build /seatly /seatly
EXPOSE 8080
USER nobody
CMD ["/seatly"]
