package main

import "fmt"

// Ôn tập Day 1–4: máy bán vé trên terminal.
// Xem các bước trong README.md. Copy lại hàm cần dùng từ labs/day04/main.go.

func giaVe(hang string) int {
	switch hang {
	case "VIP":
		return 200000
	case "Thường":
		return 100000
	case "Sinh viên":
		return 50000
	default:
		return 0
	}
}

func tinhTien(soVe int, hang string) (int, int) {
	gia := giaVe(hang)
	if gia == 0 {
		return 0, 0
	}

	tongTien := gia * soVe
	var giamGia int
	if soVe >= 4 {
		giamGia = tongTien * 10 / 100
		tongTien -= giamGia
	}
	return tongTien, giamGia
}

func kiemTraSoVe(soVe int) (bool, string) {
	if soVe <= 0 {
		return false, "Phải mua ít nhất 1 vé"
	} else if soVe > 10 {
		return false, "Tối đa 10 vé mỗi lần"
	}

	return true, "OK"
}

func main() {
	fmt.Println("🎫 Máy bán vé Seatly")

	// TODO Bước 1: in menu 3 hạng vé
	fmt.Println("1. VIP")
	fmt.Println("2. Thường")
	fmt.Println("3. Sinh viên")

	// TODO Bước 2: nhập lựa chọn (1-3), sai thì báo lỗi và dừng
	var luaChon int
	fmt.Print("Nhập lựa chọn (1-3): ")
	_, err := fmt.Scan(&luaChon)
	if err != nil {
		fmt.Println("❌ Lựa chọn không hợp lệ")
		return
	}

	// TODO Bước 3: nhập số vé, kiểm tra hợp lệ
	var soVe int
	fmt.Print("Nhập số vé: ")
	_, err = fmt.Scan(&soVe)
	if err != nil {
		fmt.Println("❌ Vui lòng nhập một con số")
		return
	}

	isValid, message := kiemTraSoVe(soVe)
	if !isValid {
		fmt.Println(message)
		return
	}

	// TODO Bước 4: tính tiền và in hóa đơn
	var hang string
	switch luaChon {
	case 1:
		hang = "VIP"
	case 2:
		hang = "Thường"
	case 3:
		hang = "Sinh viên"
	default:
		fmt.Println("Lựa chọn không hợp lệ")
		return
	}

	tongTien, giamGia := tinhTien(soVe, hang)
	fmt.Printf("Hóa đơn:\n")
	fmt.Printf("Hạng: %v\n", hang)
	fmt.Printf("Số vé: %v\n", soVe)
	fmt.Printf("Tổng tiền: %v\n", tongTien)
	if giamGia > 0 {
		fmt.Printf("Giảm giá: %v\n", giamGia)
	}
}
