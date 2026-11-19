import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getEntries, reverse } from "./api";
import { formatTime, formatVND, kindLabel } from "./format";

export default function History({ accountId }: { accountId: number }) {
  const queryClient = useQueryClient();
  const entriesQuery = useQuery({
    queryKey: ["entries", accountId],
    queryFn: () => getEntries(accountId),
  });

  const reversal = useMutation({
    mutationFn: reverse,
    onSettled: () => queryClient.invalidateQueries(),
  });

  const entries = entriesQuery.data ?? [];

  return (
    <section className="card stack" aria-labelledby="history-title">
      <h2 id="history-title">Sao kê</h2>

      {entriesQuery.isPending && <p className="muted">Đang tải sao kê…</p>}
      {entriesQuery.isError && (
        <p className="notice error" role="alert">
          Không tải được sao kê: {entriesQuery.error.message}
        </p>
      )}
      {entriesQuery.isSuccess && entries.length === 0 && (
        <p className="muted">Chưa có giao dịch nào.</p>
      )}
      {reversal.isError && (
        <p className="notice error" role="alert">
          {reversal.error.message}
        </p>
      )}

      {entries.length > 0 && (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Thời gian</th>
                <th>Giao dịch</th>
                <th className="num">Biến động</th>
                <th className="num">Số dư sau</th>
                <th>
                  <span className="sr-only">Thao tác</span>
                </th>
              </tr>
            </thead>
            <tbody>
              {entries.map((e) => (
                <tr key={e.id}>
                  <td className="muted small">{formatTime(e.createdAt)}</td>
                  <td>
                    #{e.transactionId} {kindLabel[e.kind]}
                    {e.reversed && <span className="badge">Đã đảo ngược</span>}
                  </td>
                  <td
                    className={`num ${e.direction === "CREDIT" ? "credit" : "debit"}`}
                  >
                    {e.direction === "CREDIT" ? "+" : "−"}
                    {formatVND(e.amount)}
                  </td>
                  <td className="num">{formatVND(e.balanceAfter)}</td>
                  <td>
                    {e.kind !== "REVERSAL" && !e.reversed && (
                      <button
                        className="btn small"
                        disabled={reversal.isPending}
                        onClick={() => reversal.mutate(e.transactionId)}
                      >
                        Đảo ngược
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
