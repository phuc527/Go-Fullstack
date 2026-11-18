package task

import (
	"errors"
	"time"
)

// Store lưu task vào một file JSON.
type Store struct {
	path   string
	tasks  []Task
	nextID int
	now    func() time.Time // inject đồng hồ để test được thời gian
}

// fileData là format được ghi xuống đĩa.
type fileData struct {
	NextID int    `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}

// Open đọc store từ path. File chưa tồn tại thì trả store rỗng (không phải lỗi).
//
// TODO(lab1):
//   - os.ReadFile(path)
//   - Nếu errors.Is(err, fs.ErrNotExist) -> store rỗng, nextID = 1
//   - Lỗi khác -> wrap: fmt.Errorf("open store %s: %w", path, err)
//   - json.Unmarshal vào fileData; lỗi parse cũng phải wrap
func Open(path string) (*Store, error) {
	return nil, errors.New("TODO: implement Open")
}

// Add tạo task mới, gán ID tăng dần. Chỉ thay đổi bộ nhớ; gọi Save để ghi file.
//
// TODO(lab1): dùng New(...) để tái sử dụng validation, chỉ tăng nextID khi thành công.
func (s *Store) Add(title string) (Task, error) {
	return Task{}, errors.New("TODO: implement Add")
}

// List trả task theo filter. filter == "" nghĩa là lấy tất cả.
//
// TODO(lab1): trả về một slice MỚI, không trả s.tasks trực tiếp
// (vì sao? caller sửa slice sẽ sửa luôn dữ liệu bên trong Store).
func (s *Store) List(filter Status) []Task {
	return nil
}

// Done đánh dấu task hoàn thành.
//
// TODO(lab1):
//   - Không tìm thấy -> ErrNotFound (wrap kèm id: fmt.Errorf("task %d: %w", id, ErrNotFound))
//   - Cẩn thận: `for _, t := range s.tasks` cho bạn một BẢN SAO.
//     Phải dùng index: s.tasks[i].MarkDone(...)
func (s *Store) Done(id int) error {
	return errors.New("TODO: implement Done")
}

// Delete xóa task theo ID. Không tìm thấy -> ErrNotFound.
//
// TODO(lab1): gợi ý slices.DeleteFunc (Go 1.21+) hoặc tự viết bằng append.
func (s *Store) Delete(id int) error {
	return errors.New("TODO: implement Delete")
}

// Save ghi store xuống đĩa một cách atomic.
//
// TODO(lab1):
//   - json.MarshalIndent(fileData{...}, "", "  ")
//   - Ghi vào file tạm path+".tmp", sau đó os.Rename sang path
//     (rename là atomic -> crash giữa chừng không làm hỏng file cũ)
func (s *Store) Save() error {
	return errors.New("TODO: implement Save")
}
