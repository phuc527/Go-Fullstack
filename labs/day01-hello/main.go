package main

import "fmt"

func main() {
	// Bài 1: chạy thử dòng này
	fmt.Println("🎫 Chào mừng đến Seatly!")

	// TODO Bài 2: dùng fmt.Println in ra: Sự kiện: Đêm nhạc Mùa Thu
	fmt.Println("Sự kiện: Đêm nhạc Mùa Thu")

	// TODO Bài 3: dùng fmt.Printf in ra: Ghế: A12 | Giá: 500000 đ
	fmt.Printf("Ghế: A12 | Giá: 500000 đ\n")

	// TODO Bài 4: tạo 3 biến suKien, ghe, gia rồi dùng chúng in ra: Đêm nhạc Mùa Thu - Ghế A12 - 500000 đ
	suKien := "Đêm nhạc Mùa Thu"
	ghe := "A12"
	gia := 500000
	fmt.Printf("%s - Ghế %s - %d đ\n", suKien, ghe, gia)
}
