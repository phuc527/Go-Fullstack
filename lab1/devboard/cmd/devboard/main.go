// Command devboard: CLI quản lý task.
//
//	devboard add "Học goroutine"
//	devboard list [-status todo|done]
//	devboard done 1
//	devboard delete 1
//
// File dữ liệu mặc định: ./tasks.json (đổi bằng biến môi trường DEVBOARD_FILE).
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/yourname/devboard/internal/task"
)

func main() {
	// main chỉ làm 1 việc: gọi run và xử lý exit code.
	// Logic nằm trong run(...) để có thể test bằng cách truyền args + buffer.
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New(usage)
	}

	path := os.Getenv("DEVBOARD_FILE")
	if path == "" {
		path = "tasks.json"
	}

	store, err := task.Open(path)
	if err != nil {
		return err
	}

	cmd, rest := args[0], args[1:]
	switch cmd {
	case "add":
		title := strings.Join(rest, " ")
		t, err := store.Add(title)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Added #%d: %s\n", t.ID, t.Title)
		return store.Save()

	case "list":
		fs := flag.NewFlagSet("list", flag.ContinueOnError)
		status := fs.String("status", "", "filter: todo | done")
		if err := fs.Parse(rest); err != nil {
			return err
		}
		printTasks(out, store.List(task.Status(*status)))
		return nil

	case "done", "delete":
		id, err := parseID(rest)
		if err != nil {
			return err
		}
		if cmd == "done" {
			err = store.Done(id)
		} else {
			err = store.Delete(id)
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s #%d ok\n", cmd, id)
		return store.Save()

	default:
		return fmt.Errorf("unknown command %q\n%s", cmd, usage)
	}
}

func parseID(args []string) (int, error) {
	if len(args) != 1 {
		return 0, errors.New("expected exactly one task ID")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid ID %q: %w", args[0], err)
	}
	return id, nil
}

func printTasks(out io.Writer, tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Fprintln(out, "No tasks.")
		return
	}
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tTITLE\tCREATED")
	for _, t := range tasks {
		mark := "[ ]"
		if t.Status == task.StatusDone {
			mark = "[x]"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, mark, t.Title, t.CreatedAt.Format("2006-01-02 15:04"))
	}
	w.Flush()
}

const usage = `usage:
  devboard add <title>
  devboard list [-status todo|done]
  devboard done <id>
  devboard delete <id>`
