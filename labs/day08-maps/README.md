# Day 8 — Map: trạng thái từng ghế

Hôm nay: **map**, giống object `{}` / `Map` của JS. Thay vì xóa ghế khỏi slice như Day 7, giờ ta lưu **trạng thái từng ghế**: `"A1" → "trống"`, `"A2" → "đã đặt"`. Mất khoảng 1 giờ.

## 1. Tạo, đọc, ghi, xóa (10')

```go
trangThai := map[string]string{   // key là string, value là string
	"A1": "trống",
	"A2": "đã đặt",                 // ⚠️ dòng cuối cũng phải có dấu phẩy
}

trangThai["A2"]                   // đọc: "đã đặt"
trangThai["A3"] = "trống"         // thêm mới hoặc sửa
delete(trangThai, "A3")           // xóa
len(trangThai)                    // số phần tử
```

`map[string]int` đọc là "map từ string sang int".

| JavaScript                  | Go                                              |
| --------------------------- | ----------------------------------------------- |
| `const m = { A1: "trống" }` | `m := map[string]string{"A1": "trống"}`         |
| `m["A1"]` / `m.A1`          | `m["A1"]` (không có `m.A1`)                     |
| `delete m["A1"]`            | `delete(m, "A1")`                               |
| `Object.keys(m).length`     | `len(m)`                                        |
| `"A1" in m`                 | `_, ok := m["A1"]` (xem phần 2)                 |
| Key thiếu → `undefined`     | Key thiếu → **zero value** (`""`, `0`, `false`) |

## 2. Kiểm tra key có tồn tại không: `v, ok` (10')

Vì key thiếu trả về zero value nên `trangThai["Z9"]` trả `""`. Bạn không phân biệt được "ghế không có" với "ghế có nhưng trạng thái rỗng". Go giải quyết bằng cách trả thêm giá trị thứ 2:

```go
tt, ok := trangThai["Z9"]
if !ok {
	fmt.Println("Z9 không tồn tại")
}
```

`ok` là `true` nếu key có trong map. Cách viết này gọi là **"comma ok"**, và bạn sẽ gặp lại nhiều lần trong Go.

Viết gọn được thành 1 dòng (`if` có câu lệnh khởi tạo trước dấu `;`):

```go
if _, ok := trangThai["Z9"]; !ok {
	fmt.Println("Z9 không tồn tại")
}
```

## 3. Mẹo đếm: zero value rất có ích (5')

Vì key thiếu trả `0`, đếm không cần kiểm tra trước:

```go
dem := map[string]int{}
dem["trống"]++     // chưa có → 0 → thành 1
dem["trống"]++     // 2
```

JS phải viết `dem[k] = (dem[k] || 0) + 1`. Go chỉ cần `dem[k]++`.

## 4. ⚠️ 2 bẫy (10')

**Bẫy 1 — Thứ tự duyệt map là ngẫu nhiên.**

```go
for ma, tt := range trangThai {
	fmt.Println(ma, tt)
}
```

Chạy 3 lần có thể ra 3 thứ tự khác nhau (mình đã thử: lần 3 ra `A3 A4 A5 A1 A2`). Go **cố tình** xáo trộn để bạn không viết code dựa vào thứ tự. Muốn in theo thứ tự thì sắp xếp key trước. Công thức (cần `import "maps"` và `"slices"`):

```go
for _, ma := range slices.Sorted(maps.Keys(trangThai)) {
	fmt.Println(ma, trangThai[ma])
}
```

Đọc là: lấy các key → sắp xếp → duyệt. Tạm thời cứ nhớ công thức này.

**Bẫy 2 — Map chưa tạo thì không ghi được.**

```go
var m map[string]int   // m là nil (chưa tạo)
fmt.Println(m["a"])     // đọc OK → 0
m["a"] = 1              // 💥 panic: assignment to entry in nil map
```

Luôn tạo map bằng `map[string]int{}` hoặc `make(map[string]int)`. Đừng chỉ dùng `var`.

## 5. Map truyền vào hàm thì hàm sửa được

Day 5 có nói: Go **copy** giá trị khi truyền vào hàm, nên muốn hàm sửa biến `int` thì phải dùng `&`. Map thì khác: bên trong map chỉ giữ "địa chỉ" của dữ liệu thật, nên hàm nhận map sửa vào là **bên ngoài thấy luôn**. Không cần `&`. Giống object trong JS.

## 6. Bài tập (25')

Làm trong [main.go](main.go), rồi chạy `go run ./labs/day08-maps`.

**Bài 1** — Tạo `trangThai`: A1 trống, A2 đã đặt, A3 trống, A4 trống, A5 đã đặt. In trạng thái A2 và tổng số ghế.

**Bài 2** — Dùng `_, ok := ...` kiểm tra ghế `"Z9"`. Không có thì in `Z9 không tồn tại`.

**Bài 3** — Viết `datGhe(trangThai, ma)` (khung có sẵn), trả về thông báo:

- Không có trong map → `❌ Ghế không tồn tại`
- Đang `"đã đặt"` → `❌ Ghế đã có người đặt`
- Còn lại → đổi thành `"đã đặt"` rồi trả `✅ Đặt thành công`

Đặt A1, A2, Z9.

**Bài 4** — In sơ đồ theo thứ tự A1 → A5 (dùng công thức ở phần 4). Trong cùng vòng lặp, gom các ghế `"trống"` vào 1 slice rồi in ra.

**Bài 5** — Tạo `thongKe := map[string]int{}`, duyệt `trangThai` và đếm số ghế theo trạng thái. In số ghế trống và đã đặt.

**Bài 6** — Gỡ ghế hỏng A3 bằng `delete`, in lại tổng số ghế.

### Kết quả mong đợi

```
🗺️ Trạng thái ghế Seatly
A2: đã đặt
Tổng số ghế: 5
Z9 không tồn tại
--- Đặt ghế ---
Đặt A1: ✅ Đặt thành công
Đặt A2: ❌ Ghế đã có người đặt
Đặt Z9: ❌ Ghế không tồn tại
--- Sơ đồ (đã sắp xếp) ---
A1: đã đặt
A2: đã đặt
A3: trống
A4: trống
A5: đã đặt
Ghế còn trống: [A3 A4]
--- Thống kê ---
Trống: 2, đã đặt: 3
--- Gỡ ghế hỏng A3 ---
Tổng số ghế: 4
```

Để ý: A1 hiện `đã đặt` vì Bài 3 đã đặt nó **bên trong hàm**. Đó chính là phần 5: hàm sửa map thì bên ngoài thấy luôn.

<details>
<summary>Gợi ý Bài 3 (mở trước khi mở đáp án)</summary>

Dùng early return như Day 4:

```go
tt, ok := trangThai[ma]
if !ok {
	return "❌ Ghế không tồn tại"
}
if tt == "đã đặt" {
	return ...
}
// tới đây chắc chắn ghế còn trống
```

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import (
	"fmt"
	"maps"
	"slices"
)

func datGhe(trangThai map[string]string, ma string) string {
	tt, ok := trangThai[ma]
	if !ok {
		return "❌ Ghế không tồn tại"
	}
	if tt == "đã đặt" {
		return "❌ Ghế đã có người đặt"
	}
	trangThai[ma] = "đã đặt"
	return "✅ Đặt thành công"
}

func main() {
	fmt.Println("🗺️ Trạng thái ghế Seatly")

	// Bài 1
	trangThai := map[string]string{
		"A1": "trống",
		"A2": "đã đặt",
		"A3": "trống",
		"A4": "trống",
		"A5": "đã đặt",
	}
	fmt.Println("A2:", trangThai["A2"])
	fmt.Println("Tổng số ghế:", len(trangThai))

	// Bài 2
	if _, ok := trangThai["Z9"]; !ok {
		fmt.Println("Z9 không tồn tại")
	}

	// Bài 3
	fmt.Println("--- Đặt ghế ---")
	fmt.Println("Đặt A1:", datGhe(trangThai, "A1"))
	fmt.Println("Đặt A2:", datGhe(trangThai, "A2"))
	fmt.Println("Đặt Z9:", datGhe(trangThai, "Z9"))

	// Bài 4
	fmt.Println("--- Sơ đồ (đã sắp xếp) ---")
	var gheTrong []string
	for _, ma := range slices.Sorted(maps.Keys(trangThai)) {
		fmt.Printf("%s: %s\n", ma, trangThai[ma])
		if trangThai[ma] == "trống" {
			gheTrong = append(gheTrong, ma)
		}
	}
	fmt.Println("Ghế còn trống:", gheTrong)

	// Bài 5
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
```

</details>

## 7. Thử phá code (5')

1. Ở Bài 4, bỏ `slices.Sorted(maps.Keys(...))`, viết thẳng `for ma, tt := range trangThai`. Chạy 3–4 lần, thứ tự có giống nhau không?
2. Đổi Bài 5 thành `var thongKe map[string]int` → `panic: assignment to entry in nil map`.
3. Xóa dấu phẩy sau `"A5": "đã đặt"` → lỗi `syntax error: unexpected newline`.

## Xong Day 8 khi

- [ ] Chạy ra đúng như "Kết quả mong đợi"
- [ ] Biết dùng `v, ok := m[key]` để kiểm tra key
- [ ] Nhớ 2 bẫy: thứ tự ngẫu nhiên, và map `nil` không ghi được

Ngày mai (Day 9): **string**: tách mã ghế `"A12"` thành hàng `A` và số `12`, và xử lý chữ tiếng Việt.
