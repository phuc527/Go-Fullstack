package main

import "fmt"

// Bài 1: trả về giá 1 vé theo hạng. Hạng không hợp lệ → 0.
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

// Bài 2: trả về 2 giá trị (tổng phải trả, số tiền được giảm).
// Mua từ 4 vé trở lên được giảm 10%.
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

// Bài 4: kiểm tra số vé. Trả về (hợp lệ hay không, thông báo).
// <= 0 → "Phải mua ít nhất 1 vé", > 10 → "Tối đa 10 vé mỗi lần", còn lại → "OK"
func kiemTraSoVe(soVe int) (bool, string) {
	if soVe <= 0 {
		return false, "Phải mua ít nhất 1 vé"
	} else if soVe > 10 {
		return false, "Tối đa 10 vé mỗi lần"
	}

	return true, "OK"
}

func main() {
	fmt.Println("🧾 Tính tiền Seatly")

	// TODO Bài 3: gọi tinhTien cho VIP x2, Thường x4, Sinh viên x5 rồi in kết quả
	//             sau đó gọi lại Thường x4 nhưng chỉ lấy tổng (bỏ giá trị thứ 2 bằng _)
	tongTien1, giamGia1 := tinhTien(2, "VIP")
	fmt.Printf("VIP x2: Tổng = %v, Giảm = %v\n", tongTien1, giamGia1)
	tongTien2, giamGia2 := tinhTien(4, "Thường")
	fmt.Printf("Thường x4: Tổng = %v, Giảm = %v\n", tongTien2, giamGia2)
	tongTien3, giamGia3 := tinhTien(5, "Sinh viên")
	fmt.Printf("Sinh viên x5: Tổng = %v, Giảm = %v\n", tongTien3, giamGia3)
	tongTien4, _ := tinhTien(4, "Thường")
	fmt.Printf("Thường x4 (chỉ lấy tổng): %v\n", tongTien4)

	// TODO Bài 4: gọi kiemTraSoVe với 0, 12, 3 rồi in kết quả
	isValid1, message1 := kiemTraSoVe(0)
	fmt.Printf("Số vé 0: Hợp lệ = %v, Thông báo = %v\n", isValid1, message1)
	isValid2, message2 := kiemTraSoVe(12)
	fmt.Printf("Số vé 12: Hợp lệ = %v, Thông báo = %v\n", isValid2, message2)
	isValid3, message3 := kiemTraSoVe(3)
	fmt.Printf("Số vé 3: Hợp lệ = %v, Thông báo = %v\n", isValid3, message3)
}
