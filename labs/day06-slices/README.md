# Day 6 — Slice: danh sách ghế

Hôm nay: **slice**, thứ Go dùng thay cho array của JS. Bài tập: quản lý danh sách ghế. Mất khoảng 1 giờ.

## 1. Tạo slice và đọc phần tử (10')

```go
gheVIP := []string{"A1", "A2", "A3"}   // slice các string

fmt.Println(gheVIP)       // [A1 A2 A3]
fmt.Println(len(gheVIP))  // 3
fmt.Println(gheVIP[0])    // A1
gheVIP[0] = "A0"          // sửa phần tử
```

`[]string` đọc là "slice của string". Mọi phần tử **phải cùng kiểu**. JS cho `[1, "a", true]` nhưng Go thì không.

| JavaScript                  | Go                                         |
| --------------------------- | ------------------------------------------ |
| `const ghe = ["A1", "A2"]`  | `ghe := []string{"A1", "A2"}`              |
| `ghe.length`                | `len(ghe)`                                 |
| `ghe.at(-1)` (phần tử cuối) | `ghe[len(ghe)-1]` (không có index âm)      |
| `ghe[10]` → `undefined`     | `ghe[10]` → **chương trình crash** (panic) |
| `let ghe = []`              | `var ghe []string` (rỗng, `len` = 0)       |

## 2. Thêm phần tử: append (10')

```go
gheVIP = append(gheVIP, "A4")          // thêm 1
gheVIP = append(gheVIP, "A5", "A6")    // thêm nhiều
```

⚠️ **Phải gán lại** `gheVIP = append(...)`. Khác `push()` của JS: `append` không sửa slice cũ mà **trả về slice mới**. Viết `append(gheVIP, "A4")` một mình thì Go báo lỗi luôn.

Vì sao lại thế thì Day 7 sẽ giải thích. Hôm nay chỉ cần nhớ: **`x = append(x, ...)`**.

## 3. Duyệt slice: for range (10')

```go
for i, ghe := range gheVIP {   // i = vị trí (từ 0), ghe = giá trị
	fmt.Println(i, ghe)
}

for _, ghe := range gheVIP {   // không cần vị trí → dùng _ (giống Day 4)
	fmt.Println(ghe)
}
```

| JavaScript                        | Go                              |
| --------------------------------- | ------------------------------- |
| `ghe.forEach((g, i) => ...)`      | `for i, g := range ghe { ... }` |
| `for (const g of ghe)`            | `for _, g := range ghe { ... }` |
| `ghe.reduce((t, x) => t + x, 0)`  | Không có, tự viết vòng `for`    |
| `ghe.filter(...)`, `ghe.map(...)` | Không có, tự viết vòng `for`    |

Go **cố tình** không có `map/filter/reduce` cho slice. Vòng `for` dài hơn vài dòng, nhưng ai đọc cũng hiểu ngay.

## 4. Thứ nhỏ mới: fmt.Sprintf

Giống `Printf` nhưng **không in ra**, mà trả về string. Nó tương đương template string của JS:

```go
ma := fmt.Sprintf("B%d", 3)   // "B3"   ←  JS: `B${3}`
```

## 5. Bài tập (30')

Làm trong [main.go](main.go), rồi chạy `go run ./labs/day06-slices`.

**Bài 1** — Tạo `gheVIP` gồm `"A1", "A2", "A3"`. In ra slice, số ghế, ghế đầu và ghế cuối.

**Bài 2** — Thêm `"A4"` và `"A5"` bằng **1 lần** `append`. In lại kèm số ghế.

**Bài 3** — Dùng `for range` in danh sách có đánh số **bắt đầu từ 1**: `1. A1`, `2. A2`, …

**Bài 4** — Tạo slice rỗng `var hangB []string`. Dùng vòng `for` từ 1 đến 5, mỗi lượt `append` một mã ghế tạo bằng `fmt.Sprintf`. Kết quả: `[B1 B2 B3 B4 B5]`.

**Bài 5** — Viết 2 hàm nhận slice làm tham số (khung đã có sẵn):

- `tongTien(danhSachGia []int) int`: cộng tất cả giá
- `demGheTu(danhSachGia []int, mucGia int) int`: đếm số ghế có giá `>= mucGia`

Gọi với `giaGhe := []int{1000000, 1000000, 500000, 500000, 500000}`.

### Kết quả mong đợi

```
💺 Danh sách ghế Seatly
Ghế VIP: [A1 A2 A3]
Số ghế: 3
Ghế đầu: A1, ghế cuối: A3
Sau khi thêm: [A1 A2 A3 A4 A5] (5 ghế)
--- Danh sách ---
1. A1
2. A2
3. A3
4. A4
5. A5
Hàng B: [B1 B2 B3 B4 B5]
Tổng doanh thu nếu bán hết: 3500000 đ
Số ghế giá từ 1 triệu: 2
```

<details>
<summary>Gợi ý (mở trước khi mở đáp án)</summary>

- Ghế cuối: `gheVIP[len(gheVIP)-1]`. Vì sao `-1`? Vì vị trí bắt đầu từ 0, slice 3 phần tử có vị trí 0, 1, 2.
- Bài 3: `i` bắt đầu từ 0, muốn in từ 1 thì in `i+1`.
- Bài 5 `tongTien`: tạo `tong := 0`, rồi `for _, gia := range danhSachGia { tong += gia }`.

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import "fmt"

func tongTien(danhSachGia []int) int {
	tong := 0
	for _, gia := range danhSachGia {
		tong += gia
	}
	return tong
}

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

	// Bài 1
	gheVIP := []string{"A1", "A2", "A3"}
	fmt.Println("Ghế VIP:", gheVIP)
	fmt.Println("Số ghế:", len(gheVIP))
	fmt.Printf("Ghế đầu: %s, ghế cuối: %s\n", gheVIP[0], gheVIP[len(gheVIP)-1])

	// Bài 2
	gheVIP = append(gheVIP, "A4", "A5")
	fmt.Printf("Sau khi thêm: %v (%d ghế)\n", gheVIP, len(gheVIP))

	// Bài 3
	fmt.Println("--- Danh sách ---")
	for i, ghe := range gheVIP {
		fmt.Printf("%d. %s\n", i+1, ghe)
	}

	// Bài 4
	var hangB []string
	for i := 1; i <= 5; i++ {
		hangB = append(hangB, fmt.Sprintf("B%d", i))
	}
	fmt.Println("Hàng B:", hangB)

	// Bài 5
	giaGhe := []int{1000000, 1000000, 500000, 500000, 500000}
	fmt.Printf("Tổng doanh thu nếu bán hết: %d đ\n", tongTien(giaGhe))
	fmt.Println("Số ghế giá từ 1 triệu:", demGheTu(giaGhe, 1000000))
}
```

</details>

## 6. Thử phá code (5')

| Thử viết                               | Kết quả                                                                   |
| -------------------------------------- | ------------------------------------------------------------------------- |
| `append(gheVIP, "A4")` (không gán lại) | Lỗi compile: `... is not used`                                            |
| `gheVIP[10]`                           | **Chạy được nhưng crash**: `panic: index out of range [10] with length 3` |
| `gheVIP[-1]`                           | Lỗi compile: `index -1 ... must not be negative`                          |
| `[]string{"A1", 2}`                    | Lỗi compile: số `2` không phải string                                     |

Lưu ý dòng thứ 2: đây là lỗi **lúc chạy** (panic), không phải lúc compile. Go không bắt được hết mọi lỗi trước khi chạy. Đọc slice mà không chắc có phần tử thì kiểm tra `len` trước.

## Xong Day 6 khi

- [ ] Chạy ra đúng như "Kết quả mong đợi"
- [ ] Nhớ `x = append(x, ...)` luôn phải gán lại
- [ ] Viết được vòng `for range` có `_` và không có `_`

Ngày mai (Day 7): cắt slice `s[1:3]`, package `slices` (tìm, sắp xếp), và **vì sao `append` phải gán lại**.
