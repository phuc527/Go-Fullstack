# Day 1 — Chương trình Go đầu tiên

Hôm nay chỉ học 1 việc: **viết và chạy chương trình in chữ ra màn hình**. Mất khoảng 45 phút.

## 1. Chuẩn bị (10')

1. Mở terminal, gõ `go version`. Thấy `go version go1.27...` là đã cài xong.
2. VS Code: cài extension **Go** (tác giả: Go Team at Google).
3. Mở [main.go](main.go). Nếu VS Code hỏi cài thêm tools → bấm **Install All**.

## 2. Đọc hiểu file main.go (10')

```go
package main          // file này thuộc chương trình chạy được (giống "entry point")

import "fmt"          // giống: import fmt from "fmt" — fmt dùng để in ra màn hình

func main() {         // hàm main tự chạy khi chạy chương trình
	fmt.Println("🎫 Chào mừng đến Seatly!")   // giống console.log(...)
}
```

| JavaScript                     | Go                                             |
| ------------------------------ | ---------------------------------------------- |
| `console.log("hi")`            | `fmt.Println("hi")`                            |
| ``console.log(`Giá: ${gia}`)`` | `fmt.Printf("Giá: %d\n", gia)`                 |
| `node index.js`                | `go run ./labs/day01`                          |
| Không cần hàm main             | **Bắt buộc** có `package main` + `func main()` |
| Dấu `;` tùy ý                  | Không cần `;`                                  |

**Chạy thử** (ở thư mục gốc repo):

```bash
go run ./labs/day01
```

Thấy `🎫 Chào mừng đến Seatly!` là thành công.

## 3. Bài tập (20')

Làm trong [main.go](main.go), chỗ có `TODO`. Sau mỗi bài chạy lại `go run ./labs/day01`.

**Bài 2** — In thêm 1 dòng: `Sự kiện: Đêm nhạc Mùa Thu`

**Bài 3** — Dùng `fmt.Printf` in ra: `Ghế: A12 | Giá: 500000 đ`

> `Printf` giống template string: `%s` = chỗ điền chữ, `%d` = chỗ điền số nguyên, `\n` = xuống dòng.
> Ví dụ: `fmt.Printf("Xin chào %s, bạn %d tuổi\n", "Phúc", 30)`

**Bài 4** — Tạo 3 biến `suKien`, `ghe`, `gia`, rồi dùng chúng in ra: `Đêm nhạc Mùa Thu - Ghế A12 - 500000 đ`

```go
suKien := "Đêm nhạc Mùa Thu"   // := là khai báo biến, giống let trong JS
```

### Kết quả mong đợi

```
🎫 Chào mừng đến Seatly!
Sự kiện: Đêm nhạc Mùa Thu
Ghế: A12 | Giá: 500000 đ
Đêm nhạc Mùa Thu - Ghế A12 - 500000 đ
```

<details>
<summary>Đáp án (chỉ mở khi bí quá 20 phút)</summary>

```go
package main

import "fmt"

func main() {
	fmt.Println("🎫 Chào mừng đến Seatly!")

	fmt.Println("Sự kiện: Đêm nhạc Mùa Thu")

	fmt.Printf("Ghế: %s | Giá: %d đ\n", "A12", 500000)

	suKien := "Đêm nhạc Mùa Thu"
	ghe := "A12"
	gia := 500000
	fmt.Printf("%s - Ghế %s - %d đ\n", suKien, ghe, gia)
}
```

</details>

## 4. Thử phá code (5') — để quen thông báo lỗi

Mỗi lần thử 1 cái, chạy xem lỗi gì, rồi sửa lại:

1. Xóa dòng `import "fmt"` → lỗi `undefined: fmt`
2. Tạo biến `x := 1` nhưng không dùng → lỗi `declared and not used` (Go **bắt buộc** dùng mọi biến đã tạo)
3. Đổi `Println` thành `println` (chữ thường) → vẫn chạy nhưng đừng dùng, luôn dùng `fmt.Println`

## Xong Day 1 khi

- [ ] Chạy ra đúng 4 dòng như "Kết quả mong đợi"
- [ ] Hiểu `Println` khác `Printf` chỗ nào
