# Day 5 — Ôn tập: Máy bán vé trên terminal

Hôm nay **không có kiến thức mới**, chỉ thêm 1 thứ nhỏ là nhập dữ liệu từ bàn phím. Bạn gom Day 1–4 lại thành 1 chương trình thật sự dùng được. Mất khoảng 1–1.5 giờ.

## Ôn nhanh tuần 1 (5') — trả lời trong đầu, không nhìn lại

1. `:=` khác `=` chỗ nào?
2. Vì sao `1500000 * 0.05` lỗi trong Go?
3. Go có mấy kiểu vòng lặp? Viết `while` trong Go thế nào?
4. `switch` trong Go có cần `break` không?
5. Hàm trả 2 giá trị mà chỉ cần 1 thì làm sao?

## Thứ mới duy nhất: nhập từ bàn phím (10')

```go
var soVe int                 // 1. tạo biến rỗng
fmt.Print("Nhập số vé: ")    // 2. in câu hỏi (Print không xuống dòng)
_, err := fmt.Scan(&soVe)    // 3. đọc số người dùng gõ, ghi vào soVe
if err != nil {              // 4. gõ chữ thay vì số → err khác nil
	fmt.Println("Vui lòng nhập một con số")
	return
}
```

Giống `readline` / `prompt()` bên JS, nhưng có 2 thứ lạ:

- **`&soVe`**: dấu `&` đưa cho `Scan` **địa chỉ** của biến để nó ghi giá trị vào. Tạm thời cứ nhớ "Scan thì thêm `&`". Day 12 sẽ học kỹ.
- **`_, err :=`**: chính là trả 2 giá trị của Day 4! `Scan` trả về (số giá trị đọc được, lỗi). Ta bỏ cái đầu bằng `_` và giữ `err`. `err` bằng `nil` nghĩa là **không có lỗi** (`nil` giống `null` của JS).

> 👀 Đây là lần đầu bạn thấy `if err != nil`. Bạn sẽ viết dòng này hàng nghìn lần với Go.

## Bài tập

Viết trong [main.go](main.go). Các hàm `giaVe`, `tinhTien`, `kiemTraSoVe` thì **copy từ [Day 4](../day04-functions/main.go)** sang.

Chạy: `go run ./labs/day05-review-ticket-machine`, rồi gõ số trên terminal và Enter.

**Bước 1 — Menu.** Viết hàm `tenHang(luaChon int) string`: 1 → `"VIP"`, 2 → `"Thường"`, 3 → `"Sinh viên"`, số khác → `""`.
Dùng vòng `for` từ 1 đến 3 để in menu kèm giá.

**Bước 2 — Chọn hạng.** Nhập `luaChon`. Nếu `err != nil` hoặc `tenHang` trả `""` thì in `❌ Lựa chọn không hợp lệ` rồi `return`.

> Vì sao dùng số 1-3 mà không cho gõ chữ "Sinh viên"? Vì `Scan` dừng đọc ở dấu cách, nên gõ "Sinh viên" nó chỉ đọc được "Sinh". Đọc cả dòng sẽ học sau.

**Bước 3 — Nhập số vé.** Gõ chữ thì in `❌ Vui lòng nhập một con số`. Gọi `kiemTraSoVe` mà không hợp lệ thì in thông báo lỗi rồi dừng.

**Bước 4 — Hóa đơn.** Viết hàm `inHoaDon(hang string, soVe int)`:

- Tạm tính và giảm giá: lấy từ `tinhTien`
- **Phí dịch vụ 5%** tính trên số tiền **sau giảm**
- `TỔNG = sau giảm + phí`

### Kết quả mong đợi

Chọn `2`, nhập `4` vé:

```
🎫 Máy bán vé Seatly
Chọn hạng vé:
  1. VIP (1000000 đ)
  2. Thường (500000 đ)
  3. Sinh viên (300000 đ)
Nhập lựa chọn (1-3): 2
Nhập số vé: 4
===== HÓA ĐƠN =====
Hạng: Thường
Số vé: 4
Tạm tính: 2000000 đ
Giảm giá: 200000 đ
Phí dịch vụ 5%: 90000 đ
TỔNG: 1890000 đ
===================
Cảm ơn bạn!
```

Tự thử thêm các trường hợp sau, mỗi lần chạy lại chương trình:

| Bạn gõ         | Phải thấy                     |
| -------------- | ----------------------------- |
| `1` rồi `2`    | VIP, giảm 0 đ, TỔNG 2100000 đ |
| `5`            | `❌ Lựa chọn không hợp lệ`    |
| `abc`          | `❌ Lựa chọn không hợp lệ`    |
| `3` rồi `0`    | `❌ Phải mua ít nhất 1 vé`    |
| `3` rồi `muoi` | `❌ Vui lòng nhập một con số` |

<details>
<summary>Gợi ý (mở trước khi mở đáp án)</summary>

- Hàm `main` nên ngắn, chỉ gọi các hàm nhỏ: `inMenu()`, `tenHang()`, `kiemTraSoVe()`, `inHoaDon()`.
- In menu: `hang := tenHang(i)` rồi `giaVe(hang)` để lấy giá, không cần gõ tay giá lần nữa.
- Lần `Scan` thứ 2 dùng `_, err = ...` (dấu `=`, không phải `:=`), vì `err` đã được tạo ở lần đầu.
- `if !hopLe` nghĩa là "nếu không hợp lệ" (`!` giống JS).

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import "fmt"

func tenHang(luaChon int) string {
	switch luaChon {
	case 1:
		return "VIP"
	case 2:
		return "Thường"
	case 3:
		return "Sinh viên"
	default:
		return ""
	}
}

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

func inMenu() {
	fmt.Println("Chọn hạng vé:")
	for i := 1; i <= 3; i++ {
		hang := tenHang(i)
		fmt.Printf("  %d. %s (%d đ)\n", i, hang, giaVe(hang))
	}
}

func inHoaDon(hang string, soVe int) {
	sauGiam, giamGia := tinhTien(soVe, hang)
	phi := sauGiam * 5 / 100
	fmt.Println("===== HÓA ĐƠN =====")
	fmt.Printf("Hạng: %s\n", hang)
	fmt.Printf("Số vé: %d\n", soVe)
	fmt.Printf("Tạm tính: %d đ\n", sauGiam+giamGia)
	fmt.Printf("Giảm giá: %d đ\n", giamGia)
	fmt.Printf("Phí dịch vụ 5%%: %d đ\n", phi)
	fmt.Printf("TỔNG: %d đ\n", sauGiam+phi)
	fmt.Println("===================")
}

func main() {
	fmt.Println("🎫 Máy bán vé Seatly")

	inMenu()

	var luaChon int
	fmt.Print("Nhập lựa chọn (1-3): ")
	_, err := fmt.Scan(&luaChon)
	hang := tenHang(luaChon)
	if err != nil || hang == "" {
		fmt.Println("❌ Lựa chọn không hợp lệ")
		return
	}

	var soVe int
	fmt.Print("Nhập số vé: ")
	_, err = fmt.Scan(&soVe)
	if err != nil {
		fmt.Println("❌ Vui lòng nhập một con số")
		return
	}
	hopLe, thongBao := kiemTraSoVe(soVe)
	if !hopLe {
		fmt.Println("❌", thongBao)
		return
	}

	inHoaDon(hang, soVe)
	fmt.Println("Cảm ơn bạn!")
}
```

</details>

## Bài nâng cao (không bắt buộc) — bán nhiều lượt

Cho máy bán vé chạy liên tục:

- Bọc Bước 2–4 trong `for { ... }` (vòng lặp vô hạn, giống `while (true)`).
- Nhập sai thì dùng `continue` để hỏi lại từ đầu, thay vì `return`.
- Sau mỗi hóa đơn hỏi `Mua tiếp? (1 = có, 0 = không)`. Không phải `1` thì `break`.
- Cộng dồn `doanhThu`. Muốn vậy thì `inHoaDon` phải **trả về** tổng tiền: đổi thành `func inHoaDon(...) int`.
- Thoát vòng lặp thì in `💰 Tổng doanh thu: ... đ`.

> ⚠️ Ở bài này chỉ gõ **số**. Nếu gõ chữ, `Scan` để lại chữ đó trong bộ đệm và vòng lặp có thể báo lỗi liên tục. Cách đọc input chắc chắn hơn sẽ học sau.

<details>
<summary>Đáp án bài nâng cao</summary>

Thêm `return sauGiam + phi` cuối hàm `inHoaDon` (và đổi chữ ký thành `func inHoaDon(hang string, soVe int) int`), rồi thay `main`:

```go
func main() {
	fmt.Println("🎫 Máy bán vé Seatly")

	doanhThu := 0
	for {
		inMenu()

		var luaChon int
		fmt.Print("Nhập lựa chọn (1-3): ")
		_, err := fmt.Scan(&luaChon)
		hang := tenHang(luaChon)
		if err != nil || hang == "" {
			fmt.Println("❌ Lựa chọn không hợp lệ")
			continue
		}

		var soVe int
		fmt.Print("Nhập số vé: ")
		_, err = fmt.Scan(&soVe)
		if err != nil {
			fmt.Println("❌ Vui lòng nhập một con số")
			continue
		}
		hopLe, thongBao := kiemTraSoVe(soVe)
		if !hopLe {
			fmt.Println("❌", thongBao)
			continue
		}

		doanhThu += inHoaDon(hang, soVe)

		var tiep int
		fmt.Print("Mua tiếp? (1 = có, 0 = không): ")
		fmt.Scan(&tiep)
		if tiep != 1 {
			break
		}
	}
	fmt.Printf("💰 Tổng doanh thu: %d đ\n", doanhThu)
}
```

Thử: mua Thường x4, rồi chọn `9` (sai), rồi VIP x2, rồi `0` → doanh thu `3990000 đ`.

</details>

## Xong Tuần 1 khi

- [ ] Chạy đúng cả 6 trường hợp ở trên
- [ ] Trả lời được 5 câu "Ôn nhanh" mà không nhìn lại
- [ ] `main` của bạn ngắn, phần lớn logic nằm trong các hàm nhỏ

Tuần sau (Day 6–10): **slice & map**. Đây là lúc lưu được **nhiều ghế**, và bắt đầu đặt ghế thật thay vì chỉ đếm số vé.
