// Khớp với struct trong backend/internal/domain/bank.go
export type Account = {
  id: number;
  owner: string;
  currency: string;
  status: "ACTIVE" | "FROZEN";
  balance: number; // đồng, số nguyên
  isSystem: boolean;
  createdAt: string;
};

export type Entry = {
  id: number;
  transactionId: number;
  accountId: number;
  direction: "DEBIT" | "CREDIT";
  amount: number;
  balanceAfter: number;
  createdAt: string;
  kind: "DEPOSIT" | "TRANSFER" | "REVERSAL";
  reversed: boolean;
};

export type Transaction = {
  id: number;
  kind: Entry["kind"];
  amount: number;
  reversesId: number | null;
  createdAt: string;
  entries: Entry[];
};

// Mọi response của backend có dạng này (xem backend/internal/httpapi/response.go)
type ApiResponse<T> = {
  success: boolean;
  data: T;
  message: string;
};

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init);
  const body: ApiResponse<T> | null = await res.json().catch(() => null);
  if (!res.ok || !body?.success) {
    throw new Error(body?.message || `Lỗi máy chủ (${res.status})`);
  }
  return body.data;
}

function send<T>(
  method: string,
  url: string,
  body?: unknown,
  idempotencyKey?: string,
) {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  // Cùng 1 ý định chuyển tiền → cùng 1 key, kể cả khi bấm lại sau lỗi mạng.
  // Backend thấy key cũ thì trả lại kết quả cũ, KHÔNG trừ tiền lần 2.
  if (idempotencyKey) headers["Idempotency-Key"] = idempotencyKey;
  return request<T>(url, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  });
}

export const getAccounts = () => request<Account[]>("/api/accounts");

export const getEntries = (accountId: number) =>
  request<Entry[]>(`/api/accounts/${accountId}/entries`);

export const createAccount = (owner: string) =>
  send<Account>("POST", "/api/accounts", { owner });

export const setStatus = (id: number, status: Account["status"]) =>
  send<Account>("PUT", `/api/accounts/${id}/status`, { status });

export const deposit = (key: string, accountId: number, amount: number) =>
  send<Transaction>("POST", "/api/deposits", { accountId, amount }, key);

export const transfer = (
  key: string,
  fromAccountId: number,
  toAccountId: number,
  amount: number,
) =>
  send<Transaction>(
    "POST",
    "/api/transfers",
    { fromAccountId, toAccountId, amount },
    key,
  );

// Key cố định theo giao dịch gốc: bấm "Đảo ngược" 2 lần vẫn chỉ đảo 1 lần.
export const reverse = (transactionId: number) =>
  send<Transaction>(
    "POST",
    `/api/transactions/${transactionId}/reversal`,
    undefined,
    `reverse-${transactionId}`,
  );
