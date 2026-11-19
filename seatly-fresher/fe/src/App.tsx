import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { bookSeat, getSeats, type Seat } from "./api";

const formatPrice = (vnd: number) =>
  new Intl.NumberFormat("vi-VN", { style: "currency", currency: "VND" }).format(
    vnd,
  );

function groupByRow(seats: Seat[]) {
  const rows = new Map<string, Seat[]>();
  for (const seat of seats) {
    rows.set(seat.row, [...(rows.get(seat.row) ?? []), seat]);
  }
  return [...rows];
}

export default function App() {
  const queryClient = useQueryClient();
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const seatsQuery = useQuery({ queryKey: ["seats"], queryFn: getSeats });

  const booking = useMutation({
    mutationFn: bookSeat,
    onSuccess: () => setSelectedId(null),
    // Thành công hay lỗi (vd 409 có người đặt trước) đều tải lại sơ đồ mới nhất
    onSettled: () => queryClient.invalidateQueries({ queryKey: ["seats"] }),
  });

  const seats = seatsQuery.data ?? [];
  // Ghế vừa bị người khác đặt mất (409) thì không còn tính là đang chọn
  const selected = seats.find((s) => s.id === selectedId && !s.booked);

  function selectSeat(id: string) {
    booking.reset();
    setSelectedId(id === selectedId ? null : id);
  }

  return (
    <main className="page">
      <header>
        <h1>Seatly</h1>
        <p className="muted">Đêm nhạc Mùa Thu · chọn 1 ghế để đặt</p>
      </header>

      {seatsQuery.isPending && <p className="muted">Đang tải sơ đồ ghế…</p>}

      {seatsQuery.isError && (
        <div className="notice error" role="alert">
          <p>Không tải được sơ đồ ghế: {seatsQuery.error.message}</p>
          <p className="muted">Kiểm tra BE đã chạy chưa.</p>
          <button className="btn" onClick={() => seatsQuery.refetch()}>
            Thử lại
          </button>
        </div>
      )}

      {seatsQuery.isSuccess && seats.length === 0 && (
        <p className="muted">Chưa có ghế nào.</p>
      )}

      {seats.length > 0 && (
        <>
          <div className="screen" aria-hidden="true">
            Sân khấu
          </div>

          <div className="seat-map" role="group" aria-label="Sơ đồ ghế">
            {groupByRow(seats).map(([row, rowSeats]) => (
              <div className="seat-row" key={row}>
                <span className="row-label">{row}</span>
                {rowSeats.map((seat) => (
                  <button
                    key={seat.id}
                    className="seat"
                    disabled={seat.booked}
                    aria-pressed={seat.id === selected?.id}
                    aria-label={`Ghế ${seat.id}, ${formatPrice(seat.price)}${seat.booked ? ", đã đặt" : ""}`}
                    onClick={() => selectSeat(seat.id)}
                  >
                    {seat.number}
                  </button>
                ))}
              </div>
            ))}
          </div>

          <ul className="legend" aria-label="Chú thích">
            <li>
              <span className="swatch" /> Trống
            </li>
            <li>
              <span className="swatch selected" /> Đang chọn
            </li>
            <li>
              <span className="swatch booked" /> Đã đặt
            </li>
          </ul>

          <section className="summary">
            <div>
              <p className="muted">Ghế đang chọn</p>
              <p className="summary-value">
                {selected
                  ? `${selected.id} · ${formatPrice(selected.price)}`
                  : "Chưa chọn"}
              </p>
            </div>
            <button
              className="btn primary"
              disabled={!selected || booking.isPending}
              onClick={() => selected && booking.mutate(selected.id)}
            >
              {booking.isPending ? "Đang đặt…" : "Đặt ghế"}
            </button>
          </section>

          <div aria-live="polite">
            {booking.isSuccess && (
              <p className="notice success">
                Đã đặt ghế {booking.data.id}. Hẹn gặp bạn ở sự kiện!
              </p>
            )}
            {booking.isError && (
              <p className="notice error" role="alert">
                {booking.error.message}
              </p>
            )}
          </div>
        </>
      )}
    </main>
  );
}
