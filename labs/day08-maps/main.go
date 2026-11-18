package main

import "fmt"

// Bài 3: đặt ghế trong map trạng thái, trả về thông báo.
//   - ghế không có trong map → "❌ Ghế không tồn tại"
//   - ghế đang "đã đặt"      → "❌ Ghế đã có người đặt"
//   - còn lại: đổi thành "đã đặt" → "✅ Đặt thành công"
func datGhe(trangThai map[string]string, ma string) string {
	// TODO
	return ""
}

func main() {
	fmt.Println("🗺️ Trạng thái ghế Seatly")

	// TODO Bài 1: tạo map trangThai: A1 trống, A2 đã đặt, A3 trống, A4 trống, A5 đã đặt
	//             in trạng thái A2 và tổng số ghế

	// TODO Bài 2: kiểm tra ghế "Z9" bằng cách v, ok := ...

	// TODO Bài 3: đặt A1, A2, Z9 và in thông báo

	// TODO Bài 4: in sơ đồ theo thứ tự A1 → A5, rồi in danh sách ghế còn trống

	// TODO Bài 5: đếm số ghế theo từng trạng thái bằng map[string]int

	// TODO Bài 6: xóa ghế hỏng A3 khỏi map, in lại tổng số ghế
}
