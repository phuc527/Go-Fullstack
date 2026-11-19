import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, type FormEvent } from "react";
import { deposit, transfer, type Account } from "./api";
import { formatVND } from "./format";
import { useIdempotencyKey } from "./useIdempotencyKey";

type Mode = "deposit" | "transfer";
type Props = { account: Account; accounts: Account[] };

export default function MoneyForm({ account, accounts }: Props) {
  const queryClient = useQueryClient();
  const [mode, setMode] = useState<Mode>("transfer");
  const [amount, setAmount] = useState("");
  const [toId, setToId] = useState("");
  const [key, renewKey] = useIdempotencyKey();

  const recipients = accounts.filter((a) => !a.isSystem && a.id !== account.id);

  const post = useMutation({
    mutationFn: () =>
      mode === "deposit"
        ? deposit(key, account.id, Number(amount))
        : transfer(key, account.id, Number(toId), Number(amount)),
    onSuccess: () => {
      setAmount("");
      renewKey(); // giao dịch xong → lần sau là ý định mới
    },
    // Thành công hay lỗi đều tải lại số dư + sao kê mới nhất
    onSettled: () => queryClient.invalidateQueries(),
  });

  // Sửa nội dung = ý định khác → key mới (gửi key cũ với nội dung khác sẽ bị 422)
  function edit(update: () => void) {
    update();
    renewKey();
    post.reset();
  }

  function submit(e: FormEvent) {
    e.preventDefault();
    post.mutate();
  }

  return (
    <form className="card stack" onSubmit={submit}>
      <fieldset className="tabs">
        <legend className="sr-only">Loại giao dịch</legend>
        {(["transfer", "deposit"] as const).map((m) => (
          <label key={m}>
            <input
              type="radio"
              name="mode"
              checked={mode === m}
              onChange={() => edit(() => setMode(m))}
            />
            {m === "transfer" ? "Chuyển khoản" : "Nạp tiền"}
          </label>
        ))}
      </fieldset>

      {mode === "transfer" && (
        <label className="field">
          Người nhận
          <select
            value={toId}
            onChange={(e) => edit(() => setToId(e.target.value))}
            required
          >
            <option value="">Chọn tài khoản…</option>
            {recipients.map((a) => (
              <option key={a.id} value={a.id}>
                #{a.id} · {a.owner}
              </option>
            ))}
          </select>
        </label>
      )}

      <label className="field">
        Số tiền (đồng)
        <input
          type="number"
          inputMode="numeric"
          min={1}
          step={1}
          required
          value={amount}
          onChange={(e) => edit(() => setAmount(e.target.value))}
          placeholder="100000"
        />
        {Number(amount) > 0 && (
          <span className="muted small">{formatVND(Number(amount))}</span>
        )}
      </label>

      <button className="btn primary" disabled={post.isPending}>
        {post.isPending
          ? "Đang xử lý…"
          : mode === "transfer"
            ? "Chuyển tiền"
            : "Nạp tiền"}
      </button>
      <p className="muted small mono">Idempotency-Key: {key}</p>

      <div aria-live="polite">
        {post.isSuccess && (
          <p className="notice success">
            Giao dịch #{post.data.id} thành công.
          </p>
        )}
        {post.isError && (
          <p className="notice error" role="alert">
            {post.error.message}
          </p>
        )}
      </div>
    </form>
  );
}
