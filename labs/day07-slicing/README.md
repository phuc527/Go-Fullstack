# Day 7 — Cắt slice & package slices

Hôm nay: **cắt slice, tìm và xóa ghế, sắp xếp**, và trả lời câu hỏi hôm qua: **vì sao `append` phải gán lại**. Mất khoảng 1 giờ.

## 1. Cắt slice `s[a:b]` (10')

```go
hangA := []string{"A1", "A2", "A3", "A4", "A5", "A6"}
//         vị trí:   0     1     2     3     4     5

hangA[:3]    // [A1 A2 A3]   từ đầu đến trước vị trí 3
hangA[3:]    // [A4 A5 A6]   từ vị trí 3 đến hết
hangA[2:4]   // [A3 A4]      từ vị trí 2 đến TRƯỚC vị trí 4
```

Giống `arr.slice(2, 4)` của JS: lấy từ `a`, **không lấy** `b`.

## 2. ⚠️ Bẫy lớn: cắt slice KHÔNG copy (10')

Trong JS, `arr.slice()` tạo **mảng mới**. Trong Go, `s[a:b]` **dùng chung dữ liệu** với slice gốc:

```go
vip := hangA[:3]
vip[0] = "ĐÃ BÁN"
fmt.Println(hangA)   // [ĐÃ BÁN A2 A3 A4 A5 A6]  ← hàng A cũng bị đổi!
```

Hình dung slice giống một **cửa sổ nhìn vào dãy ghế**. `hangA` và `vip` là 2 cửa sổ nhìn vào **cùng một dãy ghế thật**. Sửa qua cửa sổ nào thì dãy ghế cũng đổi.

```
dãy ghế thật:  [A1][A2][A3][A4][A5][A6]
hangA:          └──────────────────────┘
vip:            └──────────┘
```

Muốn bản copy riêng thì dùng `slices.Clone(s)`.

## 3. Vì sao `append` phải gán lại (10')

Mỗi slice có 2 con số:

- `len`: đang có bao nhiêu phần tử
- `cap`: dãy ghế thật phía sau **chứa được tối đa bao nhiêu** (sức chứa)

Khi `append`:

- **Còn chỗ** (`len < cap`): ghi thêm vào dãy cũ.
- **Hết chỗ**: Go tạo **dãy mới to gấp đôi**, copy hết sang, rồi trả về slice trỏ vào **dãy mới**.

Vì có thể nhận về dãy mới, bạn **bắt buộc** lấy kết quả: `s = append(s, x)`. Bài 6 sẽ cho bạn thấy tận mắt.

## 4. Package `slices` — hộp đồ nghề có sẵn (5')

| Muốn làm               | JavaScript                         | Go                             |
| ---------------------- | ---------------------------------- | ------------------------------ |
| Có phần tử này không?  | `arr.includes(x)`                  | `slices.Contains(s, x)`        |
| Phần tử ở vị trí nào?  | `arr.indexOf(x)` (-1 nếu không có) | `slices.Index(s, x)` (cũng -1) |
| Xóa phần tử vị trí `i` | `arr.splice(i, 1)`                 | `s = slices.Delete(s, i, i+1)` |
| Sắp xếp                | `arr.sort()`                       | `slices.Sort(s)`               |
| Copy                   | `[...arr]`                         | `slices.Clone(s)`              |

Cần `import "slices"`. Thường **VS Code tự thêm khi bạn lưu file**. Nếu không thì sửa tay:

```go
import (
	"fmt"
	"slices"
)
```

## 5. Bài tập (25')

Làm trong [main.go](main.go), rồi chạy `go run ./labs/day07-slicing`.

**Bài 1** — `hangA` gồm A1 đến A6. In ra: cả hàng, VIP (3 ghế đầu), Thường (phần còn lại), ghế giữa (A3 A4).

**Bài 2** — Bẫy chia sẻ:

- `vip := hangA[:3]`, sửa `vip[0] = "ĐÃ BÁN"`, rồi in `hangA`. Nó có đổi không?
- Tạo `hangB` gồm B1 B2 B3. `banSao := slices.Clone(hangB)`, sửa `banSao[0]`, rồi in `hangB`. Lần này có đổi không?

**Bài 3** — `gheTrong` gồm A1 đến A6 (tạo slice mới, đừng dùng `hangA` vì đã bị sửa ở Bài 2). In xem `"A4"` còn trống không.

**Bài 4** — Viết `datGhe(gheTrong, ma)` (khung có sẵn):

- Tìm vị trí ghế bằng `slices.Index`
- Không thấy (`-1`) → trả lại danh sách cũ và `false`
- Thấy → trả `slices.Delete(gheTrong, i, i+1), true`

Đặt A4, rồi **đặt A4 lần nữa** để thấy lần 2 thất bại.

**Bài 5** — Sắp xếp `[]string{"C2", "A5", "B1", "A1"}`.

**Bài 6** — Chạy đoạn này và quan sát `cap`:

```go
var s []string
for i := 1; i <= 5; i++ {
	s = append(s, fmt.Sprintf("C%d", i))
	fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}
```

Mỗi lần `cap` nhảy lên (2 → 4 → 8) là lúc Go tạo dãy mới và copy sang. Đó là lý do phải gán lại.

### Kết quả mong đợi

```
✂️ Chia ghế Seatly
Hàng A: [A1 A2 A3 A4 A5 A6]
VIP (3 ghế đầu): [A1 A2 A3]
Thường (còn lại): [A4 A5 A6]
Ghế giữa: [A3 A4]
--- Bẫy chia sẻ ---
Sửa vip[0] → hàng A: [ĐÃ BÁN A2 A3 A4 A5 A6]
Dùng Clone → hàng B: [B1 B2 B3]
--- Đặt ghế ---
Còn trống A4? true
Đặt A4: thành công=true, còn trống: [A1 A2 A3 A5 A6]
Đặt A4 lần nữa: thành công=false, còn trống: [A1 A2 A3 A5 A6]
--- Sắp xếp ---
Ghế đã chọn: [A1 A5 B1 C2]
--- len & cap ---
len=1 cap=2
len=2 cap=2
len=3 cap=4
len=4 cap=4
len=5 cap=8
```

> Con số `cap` có thể hơi khác giữa các phiên bản Go. Không sao, điều cần thấy là **nó nhảy lên từng nấc**.

<details>
<summary>Gợi ý Bài 4 (mở trước khi mở đáp án)</summary>

Khi gọi `datGhe`, bạn nhận về 2 giá trị và phải gán lại `gheTrong`:

```go
gheTrong, ok := datGhe(gheTrong, "A4")   // lần 1: dùng := (vì ok là biến mới)
gheTrong, ok = datGhe(gheTrong, "A4")    // lần 2: dùng =
```

Lần 1 dùng được `:=` dù `gheTrong` đã có sẵn, vì quy tắc là **chỉ cần ít nhất 1 biến bên trái là mới** (ở đây là `ok`).

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import (
	"fmt"
	"slices"
)

func datGhe(gheTrong []string, ma string) ([]string, bool) {
	i := slices.Index(gheTrong, ma)
	if i == -1 {
		return gheTrong, false
	}
	return slices.Delete(gheTrong, i, i+1), true
}

func main() {
	fmt.Println("✂️ Chia ghế Seatly")

	// Bài 1
	hangA := []string{"A1", "A2", "A3", "A4", "A5", "A6"}
	fmt.Println("Hàng A:", hangA)
	fmt.Println("VIP (3 ghế đầu):", hangA[:3])
	fmt.Println("Thường (còn lại):", hangA[3:])
	fmt.Println("Ghế giữa:", hangA[2:4])

	// Bài 2
	fmt.Println("--- Bẫy chia sẻ ---")
	vip := hangA[:3]
	vip[0] = "ĐÃ BÁN"
	fmt.Println("Sửa vip[0] → hàng A:", hangA)

	hangB := []string{"B1", "B2", "B3"}
	banSao := slices.Clone(hangB)
	banSao[0] = "ĐÃ BÁN"
	fmt.Println("Dùng Clone → hàng B:", hangB)

	// Bài 3 + 4
	fmt.Println("--- Đặt ghế ---")
	gheTrong := []string{"A1", "A2", "A3", "A4", "A5", "A6"}
	fmt.Println("Còn trống A4?", slices.Contains(gheTrong, "A4"))

	gheTrong, ok := datGhe(gheTrong, "A4")
	fmt.Printf("Đặt A4: thành công=%v, còn trống: %v\n", ok, gheTrong)

	gheTrong, ok = datGhe(gheTrong, "A4")
	fmt.Printf("Đặt A4 lần nữa: thành công=%v, còn trống: %v\n", ok, gheTrong)

	// Bài 5
	fmt.Println("--- Sắp xếp ---")
	gheDaChon := []string{"C2", "A5", "B1", "A1"}
	slices.Sort(gheDaChon)
	fmt.Println("Ghế đã chọn:", gheDaChon)

	// Bài 6
	fmt.Println("--- len & cap ---")
	var s []string
	for i := 1; i <= 5; i++ {
		s = append(s, fmt.Sprintf("C%d", i))
		fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
	}
}
```

</details>

## 6. Thử phá code (5')

1. Sắp xếp `[]string{"A10", "A2", "A1"}` → ra `[A1 A10 A2]`! Vì string so sánh **từng ký tự**: `"1" < "2"` nên `A10` đứng trước `A2`. JS cũng bị y hệt. Sau này lưu số ghế dạng `int` riêng sẽ hết lỗi này.
2. `hangA[2:10]` → crash: `slice bounds out of range`.
3. Viết `slices.Delete(gheTrong, i, i+1)` mà không gán lại → in `gheTrong` ra sẽ thấy dữ liệu lạ (phần cuối bị xóa trống). Luôn gán lại.

## Xong Day 7 khi

- [ ] Chạy ra đúng như "Kết quả mong đợi"
- [ ] Giải thích được bằng lời: vì sao sửa `vip[0]` lại làm đổi `hangA`
- [ ] Giải thích được bằng lời: vì sao phải viết `s = append(s, x)`

Ngày mai (Day 8): **map**, giống object của JS. Lúc đó sẽ lưu được trạng thái từng ghế (`"A1" → "trống"`, `"A2" → "đã đặt"`) thay vì xóa ghế ra khỏi slice.
