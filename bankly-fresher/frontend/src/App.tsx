import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import AccountList from "./AccountList";
import AccountPanel from "./AccountPanel";
import { getAccounts } from "./api";
import CreateAccountForm from "./CreateAccountForm";

export default function App() {
  const accountsQuery = useQuery({
    queryKey: ["accounts"],
    queryFn: getAccounts,
  });
  const [selectedId, setSelectedId] = useState<number | null>(null);

  const accounts = accountsQuery.data ?? [];
  const selected = accounts.find((a) => a.id === selectedId);

  return (
    <main className="page">
      <header>
        <h1>Bankly</h1>
        <p className="muted">
          Sổ cái kế toán kép · mọi biến động đều ghi lại, không sửa, không xóa
        </p>
      </header>

      <div className="layout">
        <aside className="stack">
          {accountsQuery.isPending && (
            <p className="muted">Đang tải tài khoản…</p>
          )}

          {accountsQuery.isError && (
            <div className="notice error" role="alert">
              <p>Không tải được tài khoản: {accountsQuery.error.message}</p>
              <p className="muted">
                Kiểm tra backend và Postgres đã chạy chưa.
              </p>
              <button className="btn" onClick={() => accountsQuery.refetch()}>
                Thử lại
              </button>
            </div>
          )}

          {accountsQuery.isSuccess && (
            <AccountList
              accounts={accounts}
              selectedId={selectedId}
              onSelect={setSelectedId}
            />
          )}

          <CreateAccountForm onCreated={setSelectedId} />
        </aside>

        <section>
          {selected ? (
            <AccountPanel account={selected} accounts={accounts} />
          ) : (
            <div className="card empty">
              <p className="muted">
                Chọn một tài khoản để nạp tiền, chuyển khoản và xem sao kê.
              </p>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
