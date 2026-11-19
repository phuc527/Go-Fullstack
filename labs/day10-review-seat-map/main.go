package main

import "fmt"

// Ôn tập Tuần 2 (slice + map + string): sơ đồ phòng chiếu, gõ mã ghế để đặt.
// Xem các bước trong README.md.
//
// Các hàm nên viết:
//
//	taoPhong(gheDaBan []string) map[string]bool   — Bước 1
func taoPhong(gheDaBan []string) map[string]bool {
	daDat := make(map[string]bool)
	for _, ghe := range gheDaBan {
		daDat[ghe] = true
	}
	return daDat
}

// inSoDo(daDat map[string]bool)                 — Bước 2
func inSoDo(daDat map[string]bool) {
	fmt.Println("  1 2 3 4 5")
	for _, hang := range []string{"A", "B", "C"} {
		fmt.Print(hang, " ")
		for so := 1; so <= 5; so++ {
			ma := fmt.Sprintf("%s%d", hang, so)
			if daDat[ma] {
				fmt.Print("X ")
			} else {
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}
}

// datGhe(daDat map[string]bool, ma string) (bool, string) — Bước 3
func datGhe(daDat map[string]bool, ma string) (bool, string) {
	if daDat[ma] {
		return false, "Ghế đã được đặt"
	}
	daDat[ma] = true
	return true, "Đặt ghế thành công"
}

// giaVe(ma string) int                          — Bước 4
func giaVe(ma string) int {
	hang := ma[:1]
	switch hang {
	case "A":
		return 100
	case "B":
		return 80
	case "C":
		return 50
	default:
		return 0
	}
}

func main() {
	fmt.Println("🎬 Phòng chiếu Seatly")

	// TODO Bước 1: tạo phòng 3 hàng (A, B, C) x 5 ghế; A3, B2, B3 đã bán sẵn
	gheDaBan := []string{"A3", "B2", "B3"}
	daDat := taoPhong(gheDaBan)

	// TODO Bước 2: in sơ đồ
	inSoDo(daDat)

	// TODO Bước 3: vòng lặp: nhập mã ghế → đặt → in lại sơ đồ; gõ 0 để thoát
	fmt.Println("Nhập mã ghế để đặt (0 để thoát):")
	for {
		var ma string
		fmt.Scanln(&ma)
		if ma == "0" {
			break
		}
		_, msg := datGhe(daDat, ma)
		fmt.Println(msg)
		inSoDo(daDat)
	}

	// TODO Bước 4: in vé của bạn + tổng tiền
	fmt.Println("Vé của bạn:")
	tongTien := 0
	for ma, dat := range daDat {
		if dat {
			fmt.Printf("- %s: %d\n", ma, giaVe(ma))
			tongTien += giaVe(ma)
		}
	}
	fmt.Printf("Tổng tiền: %d\n", tongTien)

}
