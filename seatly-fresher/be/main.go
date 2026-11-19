// Seatly BE — phiên bản đơn giản nhất: dữ liệu lưu trong RAM, tắt server là mất.
// Chạy: go run .   →  http://localhost:8080/api/seats
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Seat là 1 ghế. struct giống TS type (Day 11).
// Phần `json:"id"` gọi là tag: đặt tên field khi đổi sang JSON (Day 17).
type Seat struct {
	ID     string `json:"id"`
	Row    string `json:"row"`
	Number int    `json:"number"`
	Price  int    `json:"price"`
	Booked bool   `json:"booked"`
}

// Server Go xử lý mỗi request trên 1 goroutine riêng → 2 người cùng đặt 1 ghế
// có thể chạy CÙNG LÚC. mu (mutex) bắt họ xếp hàng từng người một (Day 28).
var (
	mu    sync.Mutex
	seats = newSeats()
)

func newSeats() []Seat {
	var list []Seat
	for _, row := range []string{"A", "B", "C"} {
		price := 500000
		if row == "A" {
			price = 1000000
		}
		for n := 1; n <= 5; n++ {
			list = append(list, Seat{ID: row + strconv.Itoa(n), Row: row, Number: n, Price: price})
		}
	}
	return list
}

// GET /api/seats → danh sách ghế
func listSeats(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock() // defer: chạy khi hàm kết thúc, giống finally trong JS → luôn mở khóa
	writeJSON(w, http.StatusOK, seats)
}

// POST /api/bookings  body: {"seatId": "A1"}
func createBooking(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SeatID string `json:"seatId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.SeatID == "" {
		writeError(w, http.StatusBadRequest, "Body phải có dạng {\"seatId\": \"A1\"}")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range seats {
		if seats[i].ID != body.SeatID {
			continue
		}
		if seats[i].Booked {
			writeError(w, http.StatusConflict, "Ghế "+body.SeatID+" đã có người đặt")
			return
		}
		// seats[i] chứ không phải biến copy từ range → sửa thẳng vào slice (Day 7)
		seats[i].Booked = true
		writeJSON(w, http.StatusCreated, seats[i])
		return
	}
	writeError(w, http.StatusNotFound, "Ghế "+body.SeatID+" không tồn tại")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	// Giống app.get / app.post của Express
	mux.HandleFunc("GET /api/seats", listSeats)
	mux.HandleFunc("POST /api/bookings", createBooking)
	return mux
}

func main() {
	log.Println("🎫 Seatly BE chạy ở http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
