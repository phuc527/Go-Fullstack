# Day 2 — Biến & kiểu dữ liệu

Hôm nay: **tạo biến, 4 kiểu cơ bản, và đổi kiểu**. Bài tập: in hóa đơn 3 vé có phí dịch vụ. Mất khoảng 1 giờ.

## 1. Tạo biến (10')

```go
giaVe := 500000          // cách hay dùng nhất — Go tự đoán kiểu (int)
var soVe int = 3         // ghi rõ kiểu
var conVe bool           // không gán → Go tự cho "giá trị rỗng" (false)
const phiDichVu = 0.05   // hằng số, không đổi được
```

| JavaScript                   | Go                                          |
| ---------------------------- | ------------------------------------------- |
| `let gia = 500000`           | `gia := 500000`                             |
| `const PHI = 0.05`           | `const phiDichVu = 0.05`                    |
| `let x` → `undefined`        | `var x int` → `0` (không có `undefined`)    |
| `gia = "abc"` vẫn chạy       | **Lỗi** — biến đã là `int` thì mãi là `int` |
| `number` (1 kiểu cho tất cả) | `int` (số nguyên), `float64` (số thập phân) |

`:=` chỉ dùng **lần đầu** tạo biến. Gán lại thì dùng `=`:

```go
soVe := 3
soVe = 5     // đúng
soVe := 5    // lỗi: biến đã tồn tại
```

## 2. 4 kiểu cơ bản

| Kiểu      | Ví dụ                | Giá trị rỗng (zero value) |
| --------- | -------------------- | ------------------------- |
| `string`  | `"Đêm nhạc Mùa Thu"` | `""`                      |
| `int`     | `500000`             | `0`                       |
| `float64` | `0.05`               | `0`                       |
| `bool`    | `true`               | `false`                   |

## 3. Điều quan trọng nhất hôm nay: int × float64 = LỖI

Trong JS `500000 * 0.05` chạy bình thường. Trong Go thì **không** — 2 kiểu khác nhau không tính chung được:

```go
tamTinh := 1500000            // int
phi := tamTinh * 0.05         // ❌ lỗi
phi := float64(tamTinh) * 0.05 // ✅ đổi int → float64 trước
```

Đổi kiểu: `float64(x)`, `int(x)` — giống `Number(x)` trong JS nhưng bắt buộc phải viết.

## 4. In ra với Printf

| Ký hiệu | Ý nghĩa                                | Ví dụ output |
| ------- | -------------------------------------- | ------------ |
| `%v`    | in giá trị bất kỳ (dùng khi lười nghĩ) | `500000`     |
| `%T`    | in **kiểu** của biến                   | `int`        |
| `%d`    | số nguyên                              | `500000`     |
| `%.0f`  | số thập phân, 0 chữ số sau dấu phẩy    | `75000`      |
| `%q`    | chuỗi có dấu ngoặc kép                 | `""`         |
| `%%`    | in ra dấu `%`                          | `%`          |

## 5. Bài tập (30')

Làm trong [main.go](main.go), rồi chạy `go run ./labs/day02`.

**Bài 1** — Tạo 4 biến bằng `:=`: `tenSuKien`, `giaVe`, `soVe`, `conVe`. In mỗi biến kèm kiểu bằng `%v` và `%T`.

**Bài 2** — Tính `tamTinh = giaVe * soVe` và in ra.

**Bài 3** — Tạo `const phiDichVu = 0.05`. Tính phí và tổng thanh toán. In bằng `%.0f`.
Thử viết `tamTinh * phiDichVu` trước để thấy lỗi, rồi mới sửa.

**Bài 4** — Khai báo `var a int`, `var b string`, `var c bool` mà không gán gì. In ra xem Go cho giá trị gì.

### Kết quả mong đợi

```
🎫 Hóa đơn Seatly
Sự kiện: Đêm nhạc Mùa Thu (kiểu string)
Giá vé: 500000 (kiểu int)
Số vé: 3 (kiểu int)
Còn vé: true (kiểu bool)
Tạm tính: 1500000 đ
Phí dịch vụ (5%): 75000 đ
Tổng thanh toán: 1575000 đ
Zero value: 0 "" false
```

<details>
<summary>Gợi ý Bài 3 (mở trước khi mở đáp án)</summary>

- Phải đổi `tamTinh` sang `float64` cả khi tính phí lẫn khi cộng tổng.
- Muốn in chữ `5%` trong Printf thì viết `5%%`.

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import "fmt"

func main() {
	fmt.Println("🎫 Hóa đơn Seatly")

	// Bài 1
	tenSuKien := "Đêm nhạc Mùa Thu"
	giaVe := 500000
	soVe := 3
	conVe := true
	fmt.Printf("Sự kiện: %v (kiểu %T)\n", tenSuKien, tenSuKien)
	fmt.Printf("Giá vé: %v (kiểu %T)\n", giaVe, giaVe)
	fmt.Printf("Số vé: %v (kiểu %T)\n", soVe, soVe)
	fmt.Printf("Còn vé: %v (kiểu %T)\n", conVe, conVe)

	// Bài 2
	tamTinh := giaVe * soVe
	fmt.Printf("Tạm tính: %d đ\n", tamTinh)

	// Bài 3
	const phiDichVu = 0.05
	phi := float64(tamTinh) * phiDichVu
	tong := float64(tamTinh) + phi
	fmt.Printf("Phí dịch vụ (5%%): %.0f đ\n", phi)
	fmt.Printf("Tổng thanh toán: %.0f đ\n", tong)

	// Bài 4
	var a int
	var b string
	var c bool
	fmt.Printf("Zero value: %v %q %v\n", a, b, c)
}
```

</details>

## 6. Thử phá code (5')

1. Gán `giaVe = "năm trăm"` → lỗi gì? (Go không cho đổi kiểu biến)
2. Viết `soVe := 5` lần thứ 2 → lỗi `no new variables on left side of :=`
3. Xóa `float64(...)` ở Bài 3 → lỗi `mismatched types` hoặc `truncated`

> 💡 Ghi chú cho sau này: app thật **không dùng float cho tiền** vì dễ sai số (giống `0.1 + 0.2` trong JS). Sẽ học cách làm đúng khi tới phần thanh toán.

## Xong Day 2 khi

- [ ] Chạy ra đúng 9 dòng như "Kết quả mong đợi"
- [ ] Giải thích được `:=` khác `=` chỗ nào
- [ ] Biết vì sao phải viết `float64(tamTinh)`
