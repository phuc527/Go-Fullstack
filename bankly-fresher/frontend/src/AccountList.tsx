import type { Account } from "./api";
import { formatVND } from "./format";

type Props = {
  accounts: Account[];
  selectedId: number | null;
  onSelect: (id: number) => void;
};

export default function AccountList({ accounts, selectedId, onSelect }: Props) {
  const customers = accounts.filter((a) => !a.isSystem);
  const system = accounts.find((a) => a.isSystem);

  return (
    <nav className="card stack" aria-label="Tài khoản">
      <h2>Tài khoản</h2>

      {customers.length === 0 && (
        <p className="muted">
          Chưa có tài khoản nào. Mở tài khoản đầu tiên bên dưới.
        </p>
      )}

      <ul className="account-list">
        {customers.map((a) => (
          <li key={a.id}>
            <button
              className="account"
              aria-current={a.id === selectedId}
              onClick={() => onSelect(a.id)}
            >
              <span>
                <span className="account-owner">{a.owner}</span>
                <span className="muted small">
                  #{a.id}
                  {a.status === "FROZEN" && (
                    <span className="badge">Đóng băng</span>
                  )}
                </span>
              </span>
              <span className="amount">{formatVND(a.balance)}</span>
            </button>
          </li>
        ))}
      </ul>

      {system && (
        <button
          className="account system"
          aria-current={system.id === selectedId}
          onClick={() => onSelect(system.id)}
        >
          <span>
            <span className="account-owner">{system.owner}</span>
            <span className="muted small">
              Luôn âm đúng bằng tổng tiền khách đã nạp → tổng mọi số dư = 0
            </span>
          </span>
          <span className="amount">{formatVND(system.balance)}</span>
        </button>
      )}
    </nav>
  );
}
