package main

import "fmt"

func main() {
	fmt.Println("🎫 Hóa đơn Seatly")

	// TODO Bài 1: tạo 4 biến rồi in ra kèm kiểu của chúng (%v và %T)
	//   tenSuKien = "Đêm nhạc Mùa Thu", giaVe = 500000, soVe = 3, conVe = true
	tenSuKien := "Đêm nhạc Mùa Thu"
	giaVe := 500000
	soVe := 3
	conVe := true

	fmt.Printf("Tên sự kiện: %v, Kiểu: %T\n", tenSuKien, tenSuKien)
	fmt.Printf("Giá vé: %v, Kiểu: %T\n", giaVe, giaVe)
	fmt.Printf("Số vé: %v, Kiểu: %T\n", soVe, soVe)
	fmt.Printf("Còn vé: %v, Kiểu: %T\n", conVe, conVe)

	// TODO Bài 2: tính tamTinh = giaVe * soVe rồi in ra
	tamTinh := giaVe * soVe
	fmt.Printf("Tạm tính: %v\n", tamTinh)

	// TODO Bài 3: tạo hằng số phiDichVu = 0.05, tính phí và tổng thanh toán
	//   (gợi ý: int và float64 không nhân trực tiếp được với nhau)
	const phiDichVu = 0.05
	phi := float64(tamTinh) * phiDichVu
	tongThanhToan := float64(tamTinh) + phi
	fmt.Printf("Phí dịch vụ: %.0f\n", phi)
	fmt.Printf("Tổng thanh toán: %.0f\n", tongThanhToan)

	// TODO Bài 4: khai báo 3 biến bằng var mà KHÔNG gán giá trị, in ra xem Go cho giá trị gì
	var tenKhachHang string
	var soLuongVe int
	var tongTien float64
	fmt.Printf("Tên khách hàng: %v, Kiểu: %T\n", tenKhachHang, tenKhachHang)
	fmt.Printf("Số lượng vé: %v, Kiểu: %T\n", soLuongVe, soLuongVe)
	fmt.Printf("Tổng tiền: %v, Kiểu: %T\n", tongTien, tongTien)
}
