package main

import "fmt"

// Bài 5: cộng tất cả giá trong slice
func tongTien(danhSachGia []int) int {
	tong := 0
	for _, gia := range danhSachGia {
		tong += gia
	}
	return tong
}

// Bài 5: đếm có bao nhiêu giá >= mucGia
func demGheTu(danhSachGia []int, mucGia int) int {
	dem := 0
	for _, gia := range danhSachGia {
		if gia >= mucGia {
			dem++
		}
	}
	return dem
}

func main() {
	fmt.Println("💺 Danh sách ghế Seatly")

	// TODO Bài 1: tạo slice gheVIP gồm "A1", "A2", "A3". In slice, số ghế, ghế đầu và ghế cuối
	gheVIP := []string{"A1", "A2", "A3"}
	fmt.Println("Ghế VIP:", gheVIP)
	fmt.Println("Số ghế VIP:", len(gheVIP))
	fmt.Println("Ghế đầu:", gheVIP[0])
	fmt.Println("Ghế cuối:", gheVIP[len(gheVIP)-1])

	// TODO Bài 2: thêm "A4" và "A5" vào gheVIP bằng append, in lại
	gheVIP = append(gheVIP, "A4", "A5")
	fmt.Println("Ghế VIP sau khi thêm:", gheVIP)

	// TODO Bài 3: dùng for range in danh sách có đánh số: "1. A1", "2. A2", ...
	for i, ghe := range gheVIP {
		fmt.Printf("%d. %s\n", i+1, ghe)
	}

	// TODO Bài 4: tạo slice rỗng hangB, dùng vòng for + append + fmt.Sprintf để có [B1 B2 B3 B4 B5]
	hangB := []string{}
	for i := 1; i <= 5; i++ {
		hangB = append(hangB, fmt.Sprintf("B%d", i))
	}
	fmt.Println("Hạng B:", hangB)

	// TODO Bài 5: giaGhe := []int{1000000, 1000000, 500000, 500000, 500000}
	//             in tổng doanh thu (tongTien) và số ghế giá từ 1 triệu (demGheTu)
	giaGhe := []int{1000000, 1000000, 500000, 500000, 500000}
	fmt.Println("Tổng doanh thu:", tongTien(giaGhe))
	fmt.Println("Số ghế giá từ 1 triệu:", demGheTu(giaGhe, 1000000))
}
