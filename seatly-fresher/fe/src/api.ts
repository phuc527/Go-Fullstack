// Khớp với struct Seat trong be/main.go
export type Seat = {
  id: string;
  row: string;
  number: number;
  price: number;
  booked: boolean;
};

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init);
  const data = await res.json().catch(() => null);
  if (!res.ok) {
    // BE trả lỗi dạng {"error": "..."}
    throw new Error(data?.error ?? `Lỗi máy chủ (${res.status})`);
  }
  return data as T;
}

export const getSeats = () => request<Seat[]>("/api/seats");

export const bookSeat = (seatId: string) =>
  request<Seat>("/api/bookings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ seatId }),
  });
