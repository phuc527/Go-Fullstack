package main

import (
	"fmt"
	"slices"
)

// Bài 4: đặt ghế = xóa ghế khỏi danh sách trống.
// Trả về (danh sách trống mới, đặt thành công hay không).
// Ghế không có trong danh sách → trả lại danh sách cũ và false.
func datGhe(gheTrong []string, ma string) ([]string, bool) {
	// TODO: dùng slices.Index và slices.Delete
	i := slices.Index(gheTrong, ma)
	if i == -1 {
		return gheTrong, false
	}
	gheTrong = slices.Delete(gheTrong, i, i+1)
	return gheTrong, true
}

func main() {
	fmt.Println("✂️ Chia ghế Seatly")

	// TODO Bài 1: hangA gồm A1..A6. Cắt ra: VIP = 3 ghế đầu, Thường = phần còn lại, ghế giữa = A3 A4
	hangA := []string{"A1", "A2", "A3", "A4", "A5", "A6"}
	vip := hangA[:3]

	// TODO Bài 2: vip := hangA[:3]; sửa vip[0] = "ĐÃ BÁN" rồi in hangA. Sau đó thử với slices.Clone
	vip[0] = "ĐÃ BÁN"
	fmt.Println("Hang A sau khi sửa:", hangA)

	// TODO Bài 3 + 4: gheTrong gồm A1..A6. Kiểm tra A4 còn trống không, đặt A4, rồi đặt A4 lần nữa
	gheTrong := []string{"A1", "A2", "A3", "A4", "A5", "A6"}
	fmt.Println("Ghế còn trống:", gheTrong)
	gheTrong, success := datGhe(gheTrong, "A4")
	if success {
		fmt.Println("Đặt ghế A4 thành công")
	} else {
		fmt.Println("Đặt ghế A4 thất bại")
	}
	fmt.Println("Ghế còn trống sau khi đặt:", gheTrong)

	// TODO Bài 5: sắp xếp []string{"C2", "A5", "B1", "A1"}
	slice := []string{"C2", "A5", "B1", "A1"}
	slices.Sort(slice)
	fmt.Println("Danh sách ghế sau khi sắp xếp:", slice)

	// TODO Bài 6: chạy đoạn code len & cap trong README và quan sát
	var s []string
	for i := 1; i <= 5; i++ {
		s = append(s, fmt.Sprintf("C%d", i))
		fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
	}

}
