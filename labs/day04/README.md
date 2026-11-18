# Day 4 — Hàm & trả nhiều giá trị

Hôm nay: **viết hàm, và hàm trả về 2 giá trị cùng lúc**. Đây là thứ Go dùng **ở khắp nơi**, nhất là để báo lỗi (sẽ học ở Day 13). Mất khoảng 1 giờ.

## 1. Viết hàm (10')

```go
func giaVe(hang string) int {
	return 1000000
}
```

Đọc là: hàm `giaVe`, nhận `hang` kiểu `string`, trả về `int`.

| JavaScript / TypeScript                         | Go                                    |
| ----------------------------------------------- | ------------------------------------- |
| `function giaVe(hang: string): number { }`      | `func giaVe(hang string) int { }`     |
| Kiểu viết **trước** tên: không có (TS: sau `:`) | Kiểu viết **sau** tên, không có `:`   |
| `(a: number, b: number)`                        | `(a, b int)` — cùng kiểu thì gộp được |
| Tham số thiếu → `undefined`                     | Gọi thiếu tham số → **lỗi compile**   |
| Arrow function `() => {}`                       | Không có, luôn viết `func`            |

Hàm viết **trên hay dưới** `main` đều được. Go không quan tâm thứ tự.

## 2. Trả nhiều giá trị — điểm mới nhất hôm nay (15')

Trong JS, muốn trả 2 thứ phải gói vào object hoặc mảng:

```js
function tinhTien(soVe, hang) {
  return { tong: 1800000, giamGia: 200000 };
}
const { tong, giamGia } = tinhTien(4, "Thường");
```

Go trả thẳng 2 giá trị, không cần gói:

```go
func tinhTien(soVe int, hang string) (int, int) {   // trả 2 int
	return 1800000, 200000
}

tong, giamGia := tinhTien(4, "Thường")   // nhận cả 2
tong, _ := tinhTien(4, "Thường")         // chỉ cần tổng → bỏ cái thứ 2 bằng _
```

- Hàm trả 2 giá trị thì **phải nhận đủ 2**. Không cần cái nào thì dùng `_`.
- Nhớ Day 2: biến tạo mà không dùng thì Go báo lỗi. `_` là cách nói "tôi cố tình bỏ cái này".

> 👀 Nhìn trước: sau này bạn sẽ gặp `ketQua, err := lamGiDo()` **rất nhiều**. Đó chính là trả 2 giá trị: kết quả và lỗi.

## 3. Bài tập (30')

File [main.go](main.go) đã có sẵn **khung hàm**. Bạn thay phần `// TODO` + `return 0` bằng code thật. Chạy `go run ./labs/day04`.

**Bài 1** — Viết `giaVe(hang)`. Dùng lại `switch` của Day 3, nhưng mỗi case `return` luôn giá trị.

**Bài 2** — Viết `tinhTien(soVe, hang)` trả về `(tong, giamGia)`:

- `tamTinh = giaVe(hang) * soVe` (hàm gọi hàm)
- Từ 4 vé trở lên thì `giamGia = tamTinh * 10 / 100`, không thì `0`
- `return tamTinh - giamGia, giamGia`

**Bài 3** — Trong `main`, gọi `tinhTien` cho VIP x2, Thường x4, Sinh viên x5 và in kết quả. Sau đó gọi lại Thường x4 nhưng **chỉ lấy tổng** bằng `_`.

**Bài 4** — Viết `kiemTraSoVe(soVe)` trả về `(bool, string)`:

- `<= 0` → `false, "Phải mua ít nhất 1 vé"`
- `> 10` → `false, "Tối đa 10 vé mỗi lần"`
- còn lại → `true, "OK"`

Gọi với 0, 12, 3 và in kết quả.

### Kết quả mong đợi

```
🧾 Tính tiền Seatly
VIP x2: tổng 2000000 đ, giảm 0 đ
Thường x4: tổng 1800000 đ, giảm 200000 đ
Sinh viên x5: tổng 1350000 đ, giảm 150000 đ
Chỉ lấy tổng: 1800000 đ
Mua 0 vé → hợp lệ: false, Phải mua ít nhất 1 vé
Mua 12 vé → hợp lệ: false, Tối đa 10 vé mỗi lần
Mua 3 vé → hợp lệ: true, OK
```

<details>
<summary>Gợi ý Bài 3 & 4 (mở trước khi mở đáp án)</summary>

Nếu viết thế này sẽ bị lỗi ở dòng thứ 2:

```go
tong, giam := tinhTien(2, "VIP")
tong, giam := tinhTien(4, "Thường")   // ❌ no new variables on left side of :=
```

2 cách sửa:

- Dòng thứ 2 trở đi dùng `=` thay vì `:=`.
- **Cách gọn hơn:** viết thêm 1 hàm nhỏ `inHoaDon(soVe int, hang string)`. Hàm này gọi `tinhTien` rồi in ra, sau đó trong `main` chỉ cần gọi `inHoaDon` 3 lần. Bài 4 làm tương tự với hàm `thuMua(soVe int)`.

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import "fmt"

func giaVe(hang string) int {
	switch hang {
	case "VIP":
		return 1000000
	case "Thường":
		return 500000
	case "Sinh viên":
		return 300000
	default:
		return 0
	}
}

func tinhTien(soVe int, hang string) (int, int) {
	tamTinh := giaVe(hang) * soVe
	giamGia := 0
	if soVe >= 4 {
		giamGia = tamTinh * 10 / 100
	}
	return tamTinh - giamGia, giamGia
}

func kiemTraSoVe(soVe int) (bool, string) {
	if soVe <= 0 {
		return false, "Phải mua ít nhất 1 vé"
	}
	if soVe > 10 {
		return false, "Tối đa 10 vé mỗi lần"
	}
	return true, "OK"
}

func inHoaDon(soVe int, hang string) {
	tong, giam := tinhTien(soVe, hang)
	fmt.Printf("%s x%d: tổng %d đ, giảm %d đ\n", hang, soVe, tong, giam)
}

func thuMua(soVe int) {
	hopLe, thongBao := kiemTraSoVe(soVe)
	fmt.Printf("Mua %d vé → hợp lệ: %v, %s\n", soVe, hopLe, thongBao)
}

func main() {
	fmt.Println("🧾 Tính tiền Seatly")

	// Bài 3
	inHoaDon(2, "VIP")
	inHoaDon(4, "Thường")
	inHoaDon(5, "Sinh viên")

	tong, _ := tinhTien(4, "Thường")
	fmt.Printf("Chỉ lấy tổng: %d đ\n", tong)

	// Bài 4
	thuMua(0)
	thuMua(12)
	thuMua(3)
}
```

Để ý `kiemTraSoVe`: gặp lỗi thì `return` **ngay**, phần cuối hàm chỉ còn trường hợp đúng. Go gọi kiểu này là **early return**, và đây là cách viết chuẩn mà bạn sẽ thấy ở mọi codebase Go.

</details>

## 4. Thử phá code (5')

1. Gọi `tinhTien(4)` (thiếu `hang`) → lỗi `not enough arguments`.
2. Viết `tong := tinhTien(4, "VIP")` (nhận 1 trong khi hàm trả 2) → lỗi `assignment mismatch`.
3. Xóa dòng `return` cuối của `kiemTraSoVe` → lỗi `missing return`.

## Xong Day 4 khi

- [ ] Chạy ra đúng như "Kết quả mong đợi"
- [ ] Viết được hàm trả 2 giá trị, và biết dùng `_` để bỏ 1 giá trị
- [ ] Hiểu "early return": gặp trường hợp sai thì `return` luôn

Ngày mai (Day 5) là **ôn tập**: gom Day 1–4 làm máy bán vé trên terminal, có nhập dữ liệu từ bàn phím.
