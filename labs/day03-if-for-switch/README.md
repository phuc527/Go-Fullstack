# Day 3 — if, switch, for

Hôm nay: **rẽ nhánh và vòng lặp**. Bài tập: quầy vé tính giá theo hạng, giảm giá, in sơ đồ ghế. Mất khoảng 1 giờ.

## 1. if / else (5')

```go
if soVe >= 4 {
	fmt.Println("Được giảm giá")
} else if soVe > 0 {
	fmt.Println("Giá thường")
} else {
	fmt.Println("Chưa chọn vé")
}
```

Khác JS 2 chỗ:

- **Không** có ngoặc tròn `( )` quanh điều kiện.
- **Bắt buộc** có ngoặc nhọn `{ }`, kể cả khi chỉ có 1 dòng.

Điều kiện phải là `bool` thật. JS cho `if (soVe)` nhưng Go thì phải viết `if soVe > 0`.

## 2. switch (10')

```go
switch hang {
case "VIP":
	gia = 1000000
case "Thường", "Phổ thông":   // nhiều giá trị trong 1 case
	gia = 500000
default:
	fmt.Println("Hạng không hợp lệ")
}
```

| JavaScript                        | Go                             |
| --------------------------------- | ------------------------------ |
| Phải viết `break` cuối mỗi case   | **Tự dừng**, không cần `break` |
| Quên `break` → chạy lọt case dưới | Không bao giờ lọt              |
| `case "A": case "B":`             | `case "A", "B":`               |

## 3. for — vòng lặp DUY NHẤT của Go (15')

Go không có `while`, không có `do...while`. Chỉ có `for`, nhưng viết được 3 kiểu:

```go
// Kiểu 1: giống for trong JS (bỏ ngoặc tròn)
for i := 1; i <= 5; i++ {
	fmt.Println(i)          // 1 2 3 4 5
}

// Kiểu 2: giống while trong JS
conLai := 3
for conLai > 0 {
	conLai--
}

// Kiểu 3: lặp n lần (i chạy từ 0 đến n-1)
for i := range 5 {
	fmt.Println(i)          // 0 1 2 3 4
}
```

`break` (thoát vòng lặp) và `continue` (bỏ qua lượt này) dùng giống JS.

## 4. Bài tập (30')

Làm trong [main.go](main.go), rồi chạy `go run ./labs/day03-if-for-switch`.

**Bài 1 (switch)** — `hang := "VIP"`. Dùng `switch` tìm giá:
VIP = 1000000, Thường = 500000, Sinh viên = 300000, hạng khác thì in `Hạng không hợp lệ`.
In ra `Hạng VIP: 1000000 đ/vé`. Đổi `hang` sang `"Thường"` và `"abc"` để thử.

**Bài 2 (if/else)** — `soVe := 4`. Tính `tong = gia * soVe`.
Nếu mua từ 4 vé trở lên thì giảm 10% và in `Mua 4 vé → giảm 10%`, không thì in `Không có giảm giá`. Cuối cùng in tổng.

> Mẹo: giảm 10% mà vẫn giữ kiểu `int` thì viết `tong * 90 / 100`. Không cần đổi sang `float64` như hôm qua.

**Bài 3 (for lồng nhau)** — In sơ đồ 3 hàng, mỗi hàng 5 ghế. Ghế số 3 ở hàng 2 đã bán thì in `X`.

> Gợi ý: `fmt.Printf(" %d", ghe)` in trên cùng 1 dòng, còn `fmt.Println()` để xuống dòng khi hết hàng.

**Bài 4 (for kiểu while)** — `conLai := 3`. Mỗi lượt in `Còn N vé` rồi giảm `conLai` đi 1. Khi hết thì in `Hết vé!`.

### Kết quả mong đợi

```
🎟️ Quầy vé Seatly
Hạng VIP: 1000000 đ/vé
Mua 4 vé → giảm 10%
Tổng: 3600000 đ
--- Sơ đồ ghế ---
Hàng 1: 1 2 3 4 5
Hàng 2: 1 2 X 4 5
Hàng 3: 1 2 3 4 5
--- Bán vé ---
Còn 3 vé
Còn 2 vé
Còn 1 vé
Hết vé!
```

<details>
<summary>Gợi ý Bài 3 (mở trước khi mở đáp án)</summary>

- Vòng ngoài chạy theo hàng (1 → 3), vòng trong chạy theo ghế (1 → 5).
- Đặt tên biến hàng là `hangGhe` vì `hang` đã dùng ở Bài 1 rồi.
- Bên trong vòng trong: `if hangGhe == 2 && ghe == 3 { ... } else { ... }`

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import "fmt"

func main() {
	fmt.Println("🎟️ Quầy vé Seatly")

	// Bài 1
	hang := "VIP"
	gia := 0
	switch hang {
	case "VIP":
		gia = 1000000
	case "Thường":
		gia = 500000
	case "Sinh viên":
		gia = 300000
	default:
		fmt.Println("Hạng không hợp lệ")
	}
	fmt.Printf("Hạng %s: %d đ/vé\n", hang, gia)

	// Bài 2
	soVe := 4
	tong := gia * soVe
	if soVe >= 4 {
		tong = tong * 90 / 100
		fmt.Printf("Mua %d vé → giảm 10%%\n", soVe)
	} else {
		fmt.Println("Không có giảm giá")
	}
	fmt.Printf("Tổng: %d đ\n", tong)

	// Bài 3
	fmt.Println("--- Sơ đồ ghế ---")
	for hangGhe := 1; hangGhe <= 3; hangGhe++ {
		fmt.Printf("Hàng %d:", hangGhe)
		for ghe := 1; ghe <= 5; ghe++ {
			if hangGhe == 2 && ghe == 3 {
				fmt.Print(" X")
			} else {
				fmt.Printf(" %d", ghe)
			}
		}
		fmt.Println()
	}

	// Bài 4
	fmt.Println("--- Bán vé ---")
	conLai := 3
	for conLai > 0 {
		fmt.Printf("Còn %d vé\n", conLai)
		conLai--
	}
	fmt.Println("Hết vé!")
}
```

</details>

## 5. Thử phá code (5')

1. Viết `if (soVe >= 4) {` có ngoặc tròn → vẫn chạy, nhưng lưu file thì VS Code tự xóa ngoặc (Go có 1 kiểu format duy nhất).
2. Bỏ ngoặc nhọn: `if soVe >= 4 fmt.Println("x")` → lỗi.
3. Viết `if soVe {` → lỗi `non-boolean condition in if statement`.
4. Thử thêm `break` vào trong vòng lặp ghế khi gặp ghế đã bán. Hàng 2 bị cắt còn mấy ghế?

## Xong Day 3 khi

- [ ] Chạy ra đúng như "Kết quả mong đợi"
- [ ] Viết được 3 kiểu `for` mà không nhìn tài liệu
- [ ] Biết vì sao `switch` trong Go không cần `break`
