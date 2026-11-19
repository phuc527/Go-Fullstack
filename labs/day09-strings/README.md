# Day 9 — String: xử lý mã ghế

Hôm nay: **package `strings`, đổi string ↔ số, và chữ tiếng Việt**. Bài tập: tách mã ghế `"A12"` thành hàng `A` và số `12`, xử lý input người dùng gõ. Mất khoảng 1 giờ.

## 1. Package `strings` — giống method của string trong JS (10')

| Muốn làm              | JavaScript               | Go                                |
| --------------------- | ------------------------ | --------------------------------- |
| Viết hoa / thường     | `s.toUpperCase()`        | `strings.ToUpper(s)`              |
| Bỏ khoảng trắng 2 đầu | `s.trim()`               | `strings.TrimSpace(s)`            |
| Có chứa không?        | `s.includes("A")`        | `strings.Contains(s, "A")`        |
| Bắt đầu bằng?         | `s.startsWith("A")`      | `strings.HasPrefix(s, "A")`       |
| Tách thành mảng       | `s.split(",")`           | `strings.Split(s, ",")`           |
| Nối mảng thành chuỗi  | `arr.join(" \| ")`       | `strings.Join(arr, " \| ")`       |
| Thay thế              | `s.replaceAll("a", "b")` | `strings.ReplaceAll(s, "a", "b")` |
| Cắt chuỗi con         | `s.slice(1)`             | `s[1:]` (giống cắt slice Day 7)   |

Khác JS: không viết `s.ToUpper()` mà viết `strings.ToUpper(s)`, tức là **gọi hàm của package rồi truyền string vào**. Nhớ `import "strings"` (VS Code tự thêm khi lưu).

String trong Go **không sửa được**. `ToUpper` không đổi `s` mà trả về string mới, nên phải gán: `s = strings.ToUpper(s)`.

## 2. Đổi string ↔ số: package `strconv` (10')

JS tự đổi kiểu: `"A" + 12` ra `"A12"`. Go thì **không bao giờ** tự đổi:

```go
so := 12
"A" + so                  // ❌ lỗi: mismatched types string and int
"A" + strconv.Itoa(so)    // ✅ "A12"     (Itoa = Integer to ASCII)

n, err := strconv.Atoi("12")    // ✅ 12, nil     (giống parseInt)
n, err := strconv.Atoi("12a")   // 0, lỗi "invalid syntax"
```

Khác `parseInt("12a")` của JS (ra `12`): `Atoi` **báo lỗi** nếu chuỗi không phải số hoàn toàn. Lại gặp `if err != nil` rồi.

> ⚠️ Đừng dùng `string(12)` để đổi số thành chữ. Nó **không** ra `"12"` mà ra ký tự có mã số 12 (ký tự lạ). `go vet` sẽ cảnh báo.

## 3. ⚠️ String là dãy byte, chữ tiếng Việt chiếm nhiều byte (15')

Go lưu string theo chuẩn UTF-8: chữ `A` chiếm 1 byte, nhưng `ê`, `ạ`, `Đ` chiếm 2–3 byte.

```go
ten := "Đêm nhạc Mùa Thu"
len(ten)                        // 21  ← số BYTE, không phải số ký tự!
utf8.RuneCountInString(ten)     // 16  ← số ký tự thật
ten[:3]                         // "Đ\xc3"  ← cắt giữa chữ "ê" → hỏng
```

Muốn làm việc theo **ký tự** thì đổi sang `[]rune` (mỗi `rune` = 1 ký tự):

```go
kyTu := []rune(ten)             // slice các ký tự
len(kyTu)                       // 16
string(kyTu[:3])                // "Đêm" ✅
```

**Quy tắc đơn giản:**

- Mã ghế, email, ID… (chỉ có chữ không dấu và số) → dùng `s[a:b]` và `len(s)` thoải mái.
- Tên sự kiện, tên người (có tiếng Việt) → muốn đếm hoặc cắt thì dùng `[]rune`.
- `strings.ToUpper`, `Split`, `Contains`… **luôn đúng** với tiếng Việt, không cần lo.

## 4. Bài tập (25')

Làm trong [main.go](main.go), rồi chạy `go run ./labs/day09-strings`.

**Bài 1** — Viết `chuanHoa(ma)`: bỏ khoảng trắng 2 đầu và viết hoa. In bằng `%q` để thấy rõ dấu cách: `"  a12 "` → `"A12"`.

**Bài 2** — Viết `tachMaGhe(ma)` trả `(hang, so, ok)`:

- `len(ma) < 2` → không hợp lệ
- `hang := ma[:1]`. Nếu `hang < "A" || hang > "Z"` → không hợp lệ (string so sánh được bằng `<`, `>`)
- `so, err := strconv.Atoi(ma[1:])`. Nếu `err != nil` → không hợp lệ

Thử với `"A12"`, `"B5"`, `"A"`, `"1A"`, `"AX"`.

**Bài 3** — Tách `"A12"`, rồi tạo mã **ghế bên cạnh** `"A13"` bằng `hang + strconv.Itoa(so+1)`.

**Bài 4** — Người dùng gõ nhiều ghế cùng lúc: `"a1, A2 ,b10"` (lộn xộn hoa thường và khoảng trắng).

- `strings.Split` theo dấu `,` → `chuanHoa` từng phần → `append` vào slice
- In slice, rồi nối lại bằng `strings.Join(dsGhe, " | ")`

**Bài 5** — Với `ten := "Đêm nhạc Mùa Thu"`: in `len`, số ký tự thật, bản viết hoa, kết quả cắt sai `ten[:3]` (in bằng `%q`). Sau đó viết `rutGon(s, toiDa)`: giữ tối đa `toiDa` **ký tự**, dài hơn thì cắt và thêm `"..."`.

### Kết quả mong đợi

```
🔤 Xử lý mã ghế Seatly
Chuẩn hóa "  a12 " → "A12"
--- Tách mã ghế ---
"A12" → hàng A, số 12
"B5" → hàng B, số 5
"A" → ❌ mã ghế không hợp lệ
"1A" → ❌ mã ghế không hợp lệ
"AX" → ❌ mã ghế không hợp lệ
Ghế bên cạnh A12: A13
--- Nhập nhiều ghế ---
Danh sách: [A1 A2 B10]
Vé của bạn: A1 | A2 | B10
--- Tiếng Việt ---
len: 21
Số ký tự: 16
Viết hoa: ĐÊM NHẠC MÙA THU
Cắt sai ten[:3]: "Đ\xc3"
Rút gọn: Đêm nhạc...
Không cần rút gọn: Hài kịch
```

<details>
<summary>Gợi ý (mở trước khi mở đáp án)</summary>

- Bài 2: in 5 trường hợp sẽ lặp code, nên viết thêm hàm nhỏ `inTach(ma string)` giống `inHoaDon` ở Day 4.
- Bài 3: `tachMaGhe` trả 3 giá trị, mà bạn chắc `"A12"` hợp lệ nên bỏ `ok` bằng `_`: `hang, so, _ := tachMaGhe("A12")`.
- Bài 4: `strings.Split("a1, A2 ,b10", ",")` ra `["a1" " A2 " "b10"]`. Các phần vẫn còn dấu cách, `chuanHoa` sẽ xử lý.
- Bài 5: `[]rune(s)` → so `len` với `toiDa` → `string(kyTu[:toiDa]) + "..."`.
- Cần import: `"strconv"`, `"strings"`, `"unicode/utf8"`.

</details>

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func chuanHoa(ma string) string {
	return strings.ToUpper(strings.TrimSpace(ma))
}

func tachMaGhe(ma string) (string, int, bool) {
	if len(ma) < 2 {
		return "", 0, false
	}
	hang := ma[:1]
	if hang < "A" || hang > "Z" {
		return "", 0, false
	}
	so, err := strconv.Atoi(ma[1:])
	if err != nil {
		return "", 0, false
	}
	return hang, so, true
}

func rutGon(s string, toiDa int) string {
	kyTu := []rune(s)
	if len(kyTu) <= toiDa {
		return s
	}
	return string(kyTu[:toiDa]) + "..."
}

func inTach(ma string) {
	hang, so, ok := tachMaGhe(ma)
	if !ok {
		fmt.Printf("%q → ❌ mã ghế không hợp lệ\n", ma)
		return
	}
	fmt.Printf("%q → hàng %s, số %d\n", ma, hang, so)
}

func main() {
	fmt.Println("🔤 Xử lý mã ghế Seatly")

	// Bài 1
	nhap := "  a12 "
	fmt.Printf("Chuẩn hóa %q → %q\n", nhap, chuanHoa(nhap))

	// Bài 2
	fmt.Println("--- Tách mã ghế ---")
	inTach("A12")
	inTach("B5")
	inTach("A")
	inTach("1A")
	inTach("AX")

	// Bài 3
	hang, so, _ := tachMaGhe("A12")
	gheBenCanh := hang + strconv.Itoa(so+1)
	fmt.Println("Ghế bên cạnh A12:", gheBenCanh)

	// Bài 4
	fmt.Println("--- Nhập nhiều ghế ---")
	nhapNhieu := "a1, A2 ,b10"
	var dsGhe []string
	for _, phan := range strings.Split(nhapNhieu, ",") {
		dsGhe = append(dsGhe, chuanHoa(phan))
	}
	fmt.Println("Danh sách:", dsGhe)
	fmt.Println("Vé của bạn:", strings.Join(dsGhe, " | "))

	// Bài 5
	fmt.Println("--- Tiếng Việt ---")
	ten := "Đêm nhạc Mùa Thu"
	fmt.Println("len:", len(ten))
	fmt.Println("Số ký tự:", utf8.RuneCountInString(ten))
	fmt.Println("Viết hoa:", strings.ToUpper(ten))
	fmt.Printf("Cắt sai ten[:3]: %q\n", ten[:3])
	fmt.Println("Rút gọn:", rutGon(ten, 8))
	fmt.Println("Không cần rút gọn:", rutGon("Hài kịch", 8))
}
```

</details>

## 5. Thử phá code (5')

| Thử viết                        | Kết quả                                                                                                     |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `"A" + so` (so là int)          | Lỗi compile: `mismatched types untyped string and int`                                                      |
| `string(so)`                    | Chạy được nhưng ra ký tự lạ. `go vet` cảnh báo: `conversion from int to string yields a string of one rune` |
| `strconv.Atoi("12a")`           | `0` và lỗi `parsing "12a": invalid syntax`                                                                  |
| `s := "A12"; fmt.Println(s[0])` | In ra `65`, không phải `A`! `s[0]` là **byte** (mã số của chữ A). Muốn lấy chữ thì dùng `s[:1]`             |
| `ten[:3]`                       | `"Đ\xc3"`: cắt giữa chữ có dấu                                                                              |

## Xong Day 9 khi

- [ ] Chạy ra đúng như "Kết quả mong đợi"
- [ ] Biết dùng `strconv.Atoi` / `strconv.Itoa` thay vì `+` hay `string(...)`
- [ ] Giải thích được: vì sao `len("Đêm")` ra `5` chứ không phải `3`

Ngày mai (Day 10) là **ôn tập Tuần 2**: sơ đồ phòng chiếu trên terminal, gõ mã ghế để đặt. Gom slice + map + string của tuần này lại.
