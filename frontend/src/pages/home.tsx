import { useEffect } from "react";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { toast } from "sonner";
import { SuggestionBoard } from "@/components/suggestion-board";
import { SuggestionForm } from "@/components/suggestion-form";
import { useMe } from "@/lib/queries";

export function HomePage() {
  const { login } = useSearch({ from: "/" });
  const navigate = useNavigate();
  const me = useMe();

  // The backend sends the browser back to "/?login=failed" when the Gamma login is rejected.
  useEffect(() => {
    if (login === "failed") {
      toast.error("You do not have the authority to view suggestions", { id: "login-failed" });
      void navigate({ to: "/", search: {}, replace: true });
    }
  }, [login, navigate]);

  return (
    <div className="grid gap-12">
      <div className="mx-auto w-full max-w-2xl">
        <SuggestionForm />
      </div>
      {me.data && <SuggestionBoard />}
    </div>
  );
}
