package main

import "fmt"

// b1
// func main() {
// 	// Dùng slice để lưu trữ các số thỏa mãn điều kiện
// 	var numbers []string

// 	// Các số cần loại bỏ
// 	skipMap := map[int]bool{6: true, 48: true, 75: true, 89: true}

// 	for i := 1; i <= 100; i++ {
// 		if !skipMap[i] {
// 			numbers = append(numbers, fmt.Sprintf("%d", i))
// 		}
// 	}

// 	// Ghép các phần tử lại bằng dấu phẩy
// 	result := strings.Join(numbers, ",")
// 	fmt.Println(result)
// }

// b2
// func main() {
// 	count := 0
// 	var line []string

// 	for i := 1; i <= 100; i += 2 { // Duyệt các số lẻ
// 		line = append(line, fmt.Sprintf("%d", i))
// 		count++

// 		// Khi đủ 3 số trên một dòng
// 		if count == 3 {
// 			fmt.Println(line[0] + "," + line[1] + "," + line[2])
// 			line = nil // Reset lại mảng cho dòng tiếp theo
// 			count = 0
// 		}
// 	}

// 	// In nốt các số còn thừa ở dòng cuối cùng (nếu có)
// 	if len(line) > 0 {
// 		// Nối các phần tử còn lại bằng dấu phẩy thủ công hoặc dùng strings.Join
// 		output := ""
// 		for idx, val := range line {
// 			if idx > 0 {
// 				output += ","
// 			}
// 			output += val
// 		}
// 		fmt.Println(output)
// 	}
// }

func main() {
	var start, end int

	// Nhập số bắt đầu và kết thúc
	fmt.Print("Nhập số bắt đầu (1-10): ")
	_, err1 := fmt.Scan(&start)
	fmt.Print("Nhập số kết thúc (1-10): ")
	_, err2 := fmt.Scan(&end)

	// Kiểm tra điều kiện input
	if err1 != nil || err2 != nil {
		fmt.Println("Lỗi: Vui lòng chỉ nhập số nguyên!")
		return
	}

	if start < 1 || end > 10 || start > end {
		fmt.Println("Lỗi: Giá trị không hợp lệ! Yêu cầu: 1 <= số bắt đầu <= số kết thúc <= 10.")
		return
	}

	fmt.Println("\n==============================")
	fmt.Println("     BẢNG CỬU CHƯƠNG ĐÃ TẠO     ")
	fmt.Println("==============================\n")

	// Vòng lặp hiển thị bảng cửu chương từ start đến end
	for i := start; i <= end; i++ {
		fmt.Printf("--- BẢNG CỬU CHƯƠNG %d ---\n", i)
		for j := 1; j <= 10; j++ {
			fmt.Printf("%d x %2d = %2d\n", i, j, i*j)
		}
		fmt.Println() // Dòng trống giữa các bảng
	}
}
