import { useState } from "react";

// 1 key cho 1 ý định giao dịch:
// - Bấm "Gửi" lại sau lỗi mạng/timeout → GIỮ key → backend không ghi sổ 2 lần.
// - Giao dịch xong, hoặc người dùng sửa số tiền/người nhận → `renew()` lấy key mới.
export function useIdempotencyKey() {
  const [key, setKey] = useState(() => crypto.randomUUID());
  return [key, () => setKey(crypto.randomUUID())] as const;
}
