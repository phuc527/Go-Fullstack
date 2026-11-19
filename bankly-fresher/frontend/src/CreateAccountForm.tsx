import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, type FormEvent } from "react";
import { createAccount } from "./api";

export default function CreateAccountForm({
  onCreated,
}: {
  onCreated: (id: number) => void;
}) {
  const queryClient = useQueryClient();
  const [owner, setOwner] = useState("");

  const create = useMutation({
    mutationFn: createAccount,
    onSuccess: (acc) => {
      setOwner("");
      onCreated(acc.id);
      return queryClient.invalidateQueries({ queryKey: ["accounts"] });
    },
  });

  function submit(e: FormEvent) {
    e.preventDefault();
    create.mutate(owner);
  }

  return (
    <form className="card stack" onSubmit={submit}>
      <h2>Mở tài khoản</h2>
      <label className="field">
        Tên chủ tài khoản
        <input
          value={owner}
          onChange={(e) => setOwner(e.target.value)}
          maxLength={100}
          required
          placeholder="Nguyễn Văn A"
        />
      </label>
      <button className="btn primary" disabled={create.isPending}>
        {create.isPending ? "Đang mở…" : "Mở tài khoản"}
      </button>
      {create.isError && (
        <p className="notice error" role="alert">
          {create.error.message}
        </p>
      )}
    </form>
  );
}
