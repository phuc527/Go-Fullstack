package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Bài 1: bỏ khoảng trắng 2 đầu và viết hoa. "  a12 " → "A12"
func chuanHoa(ma string) string {
	// TODO: strings.TrimSpace, strings.ToUpper
	ma = strings.TrimSpace(ma)
	ma = strings.ToUpper(ma)
	return ma
}

// Bài 2: tách mã ghế thành (hàng, số, hợp lệ hay không). "A12" → "A", 12, true
// Không hợp lệ khi: ngắn hơn 2 ký tự, ký tự đầu không phải A-Z, phần sau không phải số.
func tachMaGhe(ma string) (string, int, bool) {
	// TODO: ma[:1], ma[1:], strconv.Atoi
	ma = chuanHoa(ma)
	if len(ma) < 2 {
		return "", 0, false
	}
	hang := ma[:1]
	soStr := ma[1:]
	if hang < "A" || hang > "Z" {
		return "", 0, false
	}
	so, err := strconv.Atoi(soStr)
	if err != nil {
		return "", 0, false
	}
	return hang, so, true
}

// Bài 5: giữ tối đa toiDa ký tự, dài hơn thì cắt và thêm "..."
func rutGon(s string, toiDa int) string {
	// TODO: []rune(s)
	runes := []rune(s)
	if len(runes) > toiDa {
		return string(runes[:toiDa]) + "..."
	}
	return s
}

func main() {
	fmt.Println("🔤 Xử lý mã ghế Seatly")

	// TODO Bài 1: chuẩn hóa "  a12 " và in kết quả bằng %q
	ma1 := chuanHoa("  a12 ")
	fmt.Printf("Bài 1: %q\n", ma1)
	// TODO Bài 2: tách "A12", "B5", "A", "1A", "AX" và in kết quả
	hang1, so1, hopLe1 := tachMaGhe("A12")
	fmt.Printf("Bài 2: %s, %d, %t\n", hang1, so1, hopLe1)
	hang2, so2, hopLe2 := tachMaGhe("B5")
	fmt.Printf("Bài 2: %s, %d, %t\n", hang2, so2, hopLe2)
	hang3, so3, hopLe3 := tachMaGhe("A")
	fmt.Printf("Bài 2: %s, %d, %t\n", hang3, so3, hopLe3)
	hang4, so4, hopLe4 := tachMaGhe("1A")
	fmt.Printf("Bài 2: %s, %d, %t\n", hang4, so4, hopLe4)
	hang5, so5, hopLe5 := tachMaGhe("AX")
	fmt.Printf("Bài 2: %s, %d, %t\n", hang5, so5, hopLe5)

	// TODO Bài 3: từ "A12" tạo mã ghế bên cạnh "A13" bằng strconv.Itoa
	ma3 := fmt.Sprintf("%s%d", hang1, so1+1)
	fmt.Printf("Bài 3: %s\n", ma3)

	// TODO Bài 4: tách "a1, A2 ,b10" thành [A1 A2 B10], rồi nối lại bằng " | "
	ma4 := strings.Split("a1, A2 ,b10", ", ")
	for i, v := range ma4 {
		ma4[i] = chuanHoa(v)
	}
	ma4Str := strings.Join(ma4, " | ")
	fmt.Printf("Bài 4: %s\n", ma4Str)

	// TODO Bài 5: với "Đêm nhạc Mùa Thu": in len, số ký tự, viết hoa, cắt sai ten[:3], rút gọn 8 ký tự
	s := "Đêm nhạc Mùa Thu"
	fmt.Printf("Bài 5: len=%d, số ký tự=%d, viết hoa=%q, cắt=%q, rút gọn=%q\n", len(s), utf8.RuneCountInString(s), strings.ToUpper(s), s[:3], rutGon(s, 8))
}
