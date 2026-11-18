package main

import (
	"fmt"
	"maps"
	"slices"
)

// Bài 3: đặt ghế trong map trạng thái, trả về thông báo.
//   - ghế không có trong map → "❌ Ghế không tồn tại"
//   - ghế đang "đã đặt"      → "❌ Ghế đã có người đặt"
//   - còn lại: đổi thành "đã đặt" → "✅ Đặt thành công"
func datGhe(trangThai map[string]string, ma string) string {
	// "comma ok": ok = false nếu key không có trong map
	tt, ok := trangThai[ma]
	if !ok {
		return "❌ Ghế không tồn tại"
	}
	if tt == "đã đặt" {
		return "❌ Ghế đã có người đặt"
	}
	// Tới đây chắc chắn ghế còn trống (early return đã loại hết trường hợp sai).
	// Map truyền vào hàm → sửa ở đây thì bên ngoài main cũng thấy, không cần &.
	trangThai[ma] = "đã đặt"
	return "✅ Đặt thành công"
}

func main() {
	fmt.Println("🗺️ Trạng thái ghế Seatly")

	// Bài 1: tạo map bằng literal (dòng cuối cũng phải có dấu phẩy)
	trangThai := map[string]string{
		"A1": "trống",
		"A2": "đã đặt",
		"A3": "trống",
		"A4": "trống",
		"A5": "đã đặt",
	}
	fmt.Println("A2:", trangThai["A2"])
	fmt.Println("Tổng số ghế:", len(trangThai))

	// Bài 2: if có câu lệnh khởi tạo trước dấu ;
	// _ vì chỉ cần biết có hay không, không cần giá trị
	if _, ok := trangThai["Z9"]; !ok {
		fmt.Println("Z9 không tồn tại")
	}

	// Bài 3
	fmt.Println("--- Đặt ghế ---")
	fmt.Println("Đặt A1:", datGhe(trangThai, "A1"))
	fmt.Println("Đặt A2:", datGhe(trangThai, "A2"))
	fmt.Println("Đặt Z9:", datGhe(trangThai, "Z9"))

	// Bài 4: thứ tự duyệt map là ngẫu nhiên → lấy key, sắp xếp, rồi mới duyệt
	fmt.Println("--- Sơ đồ (đã sắp xếp) ---")
	var gheTrong []string // slice nil vẫn append được (khác map nil)
	for _, ma := range slices.Sorted(maps.Keys(trangThai)) {
		fmt.Printf("%s: %s\n", ma, trangThai[ma])
		if trangThai[ma] == "trống" {
			gheTrong = append(gheTrong, ma)
		}
	}
	fmt.Println("Ghế còn trống:", gheTrong)

	// Bài 5: key chưa có → trả 0 → ++ thành 1. Không cần kiểm tra trước.
	// Phải tạo bằng map[string]int{}; nếu chỉ "var thongKe map[string]int" thì ghi vào sẽ panic.
	fmt.Println("--- Thống kê ---")
	thongKe := map[string]int{}
	for _, tt := range trangThai {
		thongKe[tt]++
	}
	fmt.Printf("Trống: %d, đã đặt: %d\n", thongKe["trống"], thongKe["đã đặt"])

	// Bài 6
	fmt.Println("--- Gỡ ghế hỏng A3 ---")
	delete(trangThai, "A3")
	fmt.Println("Tổng số ghế:", len(trangThai))
}
