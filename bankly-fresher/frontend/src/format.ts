export const formatVND = (dong: number) =>
  new Intl.NumberFormat("vi-VN", { style: "currency", currency: "VND" }).format(
    dong,
  );

export const formatTime = (iso: string) =>
  new Intl.DateTimeFormat("vi-VN", {
    dateStyle: "short",
    timeStyle: "medium",
  }).format(new Date(iso));

export const kindLabel = {
  DEPOSIT: "Nạp tiền",
  TRANSFER: "Chuyển khoản",
  REVERSAL: "Đảo ngược",
} as const;
