import { AlertCircle, LoaderCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

export function Loading({ label = "Loading…" }: { label?: string }) {
  return (
    <output className="flex justify-center items-center gap-3 p-12 text-muted-foreground">
      <LoaderCircle className="size-5 animate-spin" aria-hidden />
      {label}
    </output>
  );
}

export function Failure({ error, retry }: { error: Error; retry?: () => void }) {
  return (
    <div
      role="alert"
      className="rounded-xl border border-destructive/25 bg-destructive/5 p-5 text-sm"
    >
      <div className="flex gap-2 items-start font-medium text-destructive">
        <AlertCircle className="size-5 shrink-0" aria-hidden />
        {error.message}
      </div>
      {retry && (
        <Button variant="outline" size="sm" className="mt-4" onClick={retry}>
          Try again
        </Button>
      )}
    </div>
  );
}
