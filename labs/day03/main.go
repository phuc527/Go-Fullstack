package main

import "fmt"

func main() {
	fmt.Println("🎟️ Quầy vé Seatly")

	// TODO Bài 1 (switch): hang := "VIP" → tìm giá vé
	//   VIP = 1000000, Thường = 500000, Sinh viên = 300000, khác → in "Hạng không hợp lệ"
	hang := "VIP"
	var giaVe int
	switch hang {
	case "VIP":
		giaVe = 1000000
	case "Thường":
		giaVe = 500000
	case "Sinh viên":
		giaVe = 300000
	default:
		fmt.Println("Hạng không hợp lệ")
		return
	}
	fmt.Printf("Hạng: %v, Giá vé: %v\n", hang, giaVe)

	// TODO Bài 2 (if/else): soVe := 4, tính tổng; mua từ 4 vé trở lên thì giảm 10%
	soVe := 4
	tongTien := giaVe * soVe
	if soVe >= 4 {
		giamGia := tongTien * 10 / 100
		tongTien -= giamGia
		fmt.Printf("Mua %v vé, giảm giá: %v\n", soVe, giamGia)
	}
	fmt.Printf("Tổng tiền thanh toán: %v\n", tongTien)

	// TODO Bài 3 (for lồng nhau): in sơ đồ 3 hàng x 5 ghế, ghế số 3 hàng 2 đã bán → in X
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 5; j++ {
			if i == 2 && j == 3 {
				fmt.Print("X ")
			} else {
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}

	// TODO Bài 4 (for kiểu while): conLai := 3, bán từng vé cho tới khi hết
	conLai := 3
	for conLai > 0 {
		fmt.Printf("Còn lại %v vé\n", conLai)
		conLai--
	}
}
