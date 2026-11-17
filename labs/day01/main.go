package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	names := os.Args[1:] // os.Args[0] là đường dẫn tới chính binary

	// Guard clause: xử lý case lỗi trước rồi thoát sớm, không dùng else.
	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "usage: hello <name>...")
		os.Exit(1)
	}

	fmt.Printf("Xin chào, %s!\n", strings.Join(names, ", "))
}
