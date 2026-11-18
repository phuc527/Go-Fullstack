// Package day01: string, rune, slice, map.
package day01

// Reverse đảo chuỗi theo ký tự (rune), KHÔNG theo byte.
// "Việt Nam" -> "maN tệiV".
func Reverse(s string) string {
	panic("TODO")
}

// WordFreq đếm số lần xuất hiện của mỗi từ.
// Không phân biệt hoa/thường; từ = chuỗi liên tiếp chữ cái hoặc chữ số.
// Gợi ý: strings.FieldsFunc, unicode.IsLetter, unicode.IsNumber, strings.ToLower.
func WordFreq(text string) map[string]int {
	panic("TODO")
}

type WordCount struct {
	Word  string
	Count int
}

// TopN trả về n từ nhiều nhất, giảm dần theo Count; bằng nhau thì theo alphabet.
// n <= 0 -> rỗng; n > số từ -> trả hết.
// Gợi ý: slices.SortFunc, cmp.Compare.
func TopN(freq map[string]int, n int) []WordCount {
	panic("TODO")
}
