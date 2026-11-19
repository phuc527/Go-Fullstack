# Day 10 — Ôn tập Tuần 2: Sơ đồ phòng chiếu

Hôm nay **không có kiến thức mới**. Gom slice + map + string của tuần này thành 1 chương trình: in sơ đồ ghế, người dùng gõ mã ghế để đặt, cuối cùng in vé. Mất khoảng 1–1.5 giờ.

```
🎬 Phòng chiếu Seatly
     1  2  3  4  5
  A [ ][ ][X][ ][ ]
  B [ ][X][X][ ][ ]
  C [ ][ ][ ][ ][ ]
Nhập mã ghế (0 để thoát): a1
✅ Đã đặt A1
```

## Ôn nhanh tuần 2 (5') — trả lời trong đầu, không nhìn lại

1. Vì sao phải viết `s = append(s, x)` mà không chỉ `append(s, x)`?
2. `vip := hangA[:3]` rồi sửa `vip[0]` thì `hangA` có đổi không? Vì sao?
3. Đọc key không có trong map thì được gì? Làm sao biết key có tồn tại?
4. Vì sao không in sơ đồ ghế bằng `for ma := range map`?
5. `len("Đêm")` bằng mấy? Đổi `12` thành `"12"` thế nào?

## Ý tưởng chính: map để TRA CỨU, slice để GIỮ THỨ TỰ

- `daDat := map[string]bool{}`: `"A1" → false` (trống), `"A3" → true` (đã đặt). Tra 1 ghế rất nhanh.
- `cacHang := []string{"A", "B", "C"}` và vòng `for so := 1; so <= 5`: dùng để **in sơ đồ đúng thứ tự** (vì map duyệt ngẫu nhiên, Day 8).

Dùng `bool` cho gọn: zero value `false` = trống, đúng như ý ta muốn.

## Thứ nhỏ mới: biến dùng chung cho cả file

Khai báo **bên ngoài** hàm thì mọi hàm trong file đều dùng được:

```go
var cacHang = []string{"A", "B", "C"}   // ngoài hàm: phải dùng var, không dùng := được
const soGheMoiHang = 5

func main() { ... }
```

## Bài tập

Viết trong [main.go](main.go), rồi chạy `go run ./labs/day10-review-seat-map`.

**Bước 1 — `taoPhong(gheDaBan []string) map[string]bool`**

- 2 vòng `for` lồng nhau (hàng × số ghế) → gán `daDat[hang+strconv.Itoa(so)] = false`
- Duyệt `gheDaBan` → gán `true`
- Gọi: `daDat := taoPhong([]string{"A3", "B2", "B3"})`

**Bước 2 — `inSoDo(daDat map[string]bool)`**

- Dòng đầu: 4 dấu cách, rồi `" 1 "`, `" 2 "`, … `" 5 "`
- Mỗi hàng: `"  A "` rồi `[X]` hoặc `[ ]` cho từng ghế (dùng `fmt.Print` để không xuống dòng, hết hàng mới `fmt.Println()`)

**Bước 3 — Vòng lặp đặt ghế**

- `for { ... }` vô hạn, mỗi lượt: `fmt.Scan(&nhap)` → chuẩn hóa (trim + viết hoa, Day 9)
- Gõ `0` thì `break`
- Viết `datGhe(daDat, ma) (bool, string)` theo kiểu Day 8: không tồn tại / đã có người đặt / đặt thành công
- Đặt thành công thì `append` vào `gheCuaBan` và in lại sơ đồ

> Lần này `Scan` đọc vào **string** nên gõ gì cũng được, không bị kẹt như khi đọc số ở Day 5.

**Bước 4 — In vé**

- Viết `giaVe(ma)`: ghế hàng A giá `1000000`, các hàng khác `500000` (gợi ý: `strings.HasPrefix`)
- Chưa đặt ghế nào thì in `Bạn chưa đặt ghế nào`, không thì in từng ghế kèm giá và tổng

### Kết quả mong đợi

Gõ lần lượt: `a1`, `a1`, `z9`, `c5`, `0`

```
🎬 Phòng chiếu Seatly
     1  2  3  4  5
  A [ ][ ][X][ ][ ]
  B [ ][X][X][ ][ ]
  C [ ][ ][ ][ ][ ]
Nhập mã ghế (0 để thoát): a1
✅ Đã đặt A1
     1  2  3  4  5
  A [X][ ][X][ ][ ]
  B [ ][X][X][ ][ ]
  C [ ][ ][ ][ ][ ]
Nhập mã ghế (0 để thoát): a1
❌ A1 đã có người đặt
Nhập mã ghế (0 để thoát): z9
❌ Ghế Z9 không tồn tại
Nhập mã ghế (0 để thoát): c5
✅ Đã đặt C5
     1  2  3  4  5
  A [X][ ][X][ ][ ]
  B [ ][X][X][ ][ ]
  C [ ][ ][ ][ ][X]
Nhập mã ghế (0 để thoát): 0
===== VÉ CỦA BẠN =====
A1 (1000000 đ)
C5 (500000 đ)
Tổng: 1500000 đ
```

Thử thêm: chạy lại, gõ `0` ngay → phải thấy `Bạn chưa đặt ghế nào`.

<details>
<summary>Gợi ý (mở trước khi mở đáp án)</summary>

- Ghép mã ghế: `hang + strconv.Itoa(so)`, ví dụ `"A" + "1"` → `"A1"`.
- Trong `inSoDo`: `if daDat[ma] { [X] } else { [ ] }`. Không cần `v, ok` vì chắc chắn ghế có trong map.
- Trong `datGhe` thì **cần** `v, ok`, vì người dùng có thể gõ `Z9`:

  ```go
  dat, ok := daDat[ma]
  if !ok { ... không tồn tại }
  if dat { ... đã có người đặt }
  ```

- Nối chuỗi thông báo: `"✅ Đã đặt " + ma`.
- Cần import: `"strconv"`, `"strings"`.

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

var cacHang = []string{"A", "B", "C"}

const soGheMoiHang = 5

func taoPhong(gheDaBan []string) map[string]bool {
	daDat := map[string]bool{}
	for _, hang := range cacHang {
		for so := 1; so <= soGheMoiHang; so++ {
			daDat[hang+strconv.Itoa(so)] = false
		}
	}
	for _, ma := range gheDaBan {
		daDat[ma] = true
	}
	return daDat
}

func inSoDo(daDat map[string]bool) {
	fmt.Print("    ")
	for so := 1; so <= soGheMoiHang; so++ {
		fmt.Printf(" %d ", so)
	}
	fmt.Println()
	for _, hang := range cacHang {
		fmt.Printf("  %s ", hang)
		for so := 1; so <= soGheMoiHang; so++ {
			if daDat[hang+strconv.Itoa(so)] {
				fmt.Print("[X]")
			} else {
				fmt.Print("[ ]")
			}
		}
		fmt.Println()
	}
}

func giaVe(ma string) int {
	if strings.HasPrefix(ma, "A") {
		return 1000000
	}
	return 500000
}

func datGhe(daDat map[string]bool, ma string) (bool, string) {
	dat, ok := daDat[ma]
	if !ok {
		return false, "❌ Ghế " + ma + " không tồn tại"
	}
	if dat {
		return false, "❌ " + ma + " đã có người đặt"
	}
	daDat[ma] = true
	return true, "✅ Đã đặt " + ma
}

func main() {
	fmt.Println("🎬 Phòng chiếu Seatly")
	daDat := taoPhong([]string{"A3", "B2", "B3"})
	inSoDo(daDat)

	var gheCuaBan []string
	for {
		var nhap string
		fmt.Print("Nhập mã ghế (0 để thoát): ")
		fmt.Scan(&nhap)
		ma := strings.ToUpper(strings.TrimSpace(nhap))
		if ma == "0" || ma == "" {
			break
		}

		ok, thongBao := datGhe(daDat, ma)
		fmt.Println(thongBao)
		if ok {
			gheCuaBan = append(gheCuaBan, ma)
			inSoDo(daDat)
		}
	}

	fmt.Println("===== VÉ CỦA BẠN =====")
	if len(gheCuaBan) == 0 {
		fmt.Println("Bạn chưa đặt ghế nào")
		return
	}
	tong := 0
	for _, ma := range gheCuaBan {
		fmt.Printf("%s (%d đ)\n", ma, giaVe(ma))
		tong += giaVe(ma)
	}
	fmt.Printf("Tổng: %d đ\n", tong)
}
```

Vì sao có `|| ma == ""`? Nếu input bị đóng (ví dụ bấm Ctrl+Z rồi Enter trên Windows), `Scan` không đọc được gì và `nhap` vẫn là `""`. Không có điều kiện này thì vòng lặp sẽ chạy mãi.

</details>

## Bài nâng cao (không bắt buộc) — hủy ghế

Gõ `-A1` để **hủy** ghế A1 mà **chính bạn** đã đặt trong phiên này:

- `strings.HasPrefix(ma, "-")` → bỏ dấu `-` bằng `ma[1:]`
- Tìm trong `gheCuaBan` bằng `slices.Index` (Day 7). Không có thì in `❌ Bạn không đặt ghế B2`. Không cho hủy ghế người khác!
- Có thì xóa khỏi `gheCuaBan` bằng `slices.Delete`, đặt `daDat[ma] = false`, in `↩️ Đã hủy A1` và sơ đồ

<details>
<summary>Đáp án bài nâng cao</summary>

Thêm `"slices"` vào import, rồi chèn đoạn này **ngay sau** `if ma == "0" ... break`:

```go
if strings.HasPrefix(ma, "-") {
	ma = ma[1:]
	i := slices.Index(gheCuaBan, ma)
	if i == -1 {
		fmt.Println("❌ Bạn không đặt ghế " + ma)
		continue
	}
	gheCuaBan = slices.Delete(gheCuaBan, i, i+1)
	daDat[ma] = false
	fmt.Println("↩️ Đã hủy " + ma)
	inSoDo(daDat)
	continue
}
```

Thử: `a1`, `c5`, `-a1`, `-b2`, `0` → hủy A1 thành công, B2 bị từ chối (của người khác), vé còn `C5 (500000 đ)`.

</details>

## Xong Tuần 2 khi

- [ ] Chạy đúng như "Kết quả mong đợi"
- [ ] Trả lời được 5 câu "Ôn nhanh" mà không nhìn lại
- [ ] Giải thích được vì sao chương trình dùng **cả** map lẫn slice

Tuần sau (Day 11–15): **struct, method, error**. Thay vì `map[string]bool` rời rạc, bạn sẽ có `Event` và `Seat` thật sự, với `event.Book("A1")` trả về lỗi đàng hoàng.
