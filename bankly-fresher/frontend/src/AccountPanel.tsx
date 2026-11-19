import { useMutation, useQueryClient } from "@tanstack/react-query";
import { setStatus, type Account } from "./api";
import { formatVND } from "./format";
import History from "./History";
import MoneyForm from "./MoneyForm";

type Props = {
  account: Account;
  accounts: Account[];
};

export default function AccountPanel({ account, accounts }: Props) {
  const queryClient = useQueryClient();
  const frozen = account.status === "FROZEN";

  const toggle = useMutation({
    mutationFn: () => setStatus(account.id, frozen ? "ACTIVE" : "FROZEN"),
    onSettled: () => queryClient.invalidateQueries({ queryKey: ["accounts"] }),
  });

  return (
    <div className="stack">
      <div className="card balance-card">
        <div>
          <p className="muted small">
            #{account.id} · {account.owner}
            {frozen && <span className="badge">Đóng băng</span>}
          </p>
          <p className="balance">{formatVND(account.balance)}</p>
          <p className="muted small">Số dư khả dụng</p>
        </div>
        {!account.isSystem && (
          <button
            className="btn"
            disabled={toggle.isPending}
            onClick={() => toggle.mutate()}
          >
            {frozen ? "Mở băng" : "Đóng băng"}
          </button>
        )}
      </div>
      {toggle.isError && (
        <p className="notice error" role="alert">
          {toggle.error.message}
        </p>
      )}

      {!account.isSystem && (
        // key: đổi tài khoản → form mới tinh (xóa số tiền, key cũ, thông báo cũ)
        <MoneyForm
          key={`money-${account.id}`}
          account={account}
          accounts={accounts}
        />
      )}

      <History key={`history-${account.id}`} accountId={account.id} />
    </div>
  );
}
