import { useId, useState } from "react";
import { Link } from "@tanstack/react-router";
import { ChevronDown, Trash2, UserRound } from "lucide-react";
import { Collapsible } from "radix-ui";
import { RelativeTime } from "@/components/relative-time";
import { Button } from "@/components/ui/button";
import type { Suggestion } from "@/lib/api";

export function SuggestionCard({
  suggestion,
  onDelete,
}: {
  suggestion: Suggestion;
  onDelete: (suggestion: Suggestion) => void;
}) {
  const [open, setOpen] = useState(false);
  const titleId = useId();

  return (
    <Collapsible.Root open={open} onOpenChange={setOpen} asChild>
      <article
        aria-labelledby={titleId}
        className="flex flex-col gap-3 rounded-xl border bg-card p-5 text-card-foreground shadow-sm transition-shadow hover:shadow-md"
      >
        <div className="flex items-start justify-between gap-3">
          <div className="grid min-w-0 gap-1">
            <h3 id={titleId} className="text-base leading-snug font-semibold break-words">
              <Link
                to="/suggestion/$id"
                params={{ id: suggestion.id }}
                className="rounded-sm underline-offset-4 hover:text-accent-foreground hover:underline"
              >
                {suggestion.title}
              </Link>
            </h3>
            <p className="flex min-w-0 items-center gap-1.5 text-sm text-muted-foreground">
              <UserRound className="size-3.5 shrink-0" aria-hidden />
              <span className="truncate">{suggestion.author}</span>
            </p>
          </div>
          <Button
            variant="ghost"
            size="icon-sm"
            className="-mt-1 -mr-2 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
            aria-label={`Delete suggestion "${suggestion.title}"`}
            data-delete-suggestion={suggestion.id}
            onClick={() => onDelete(suggestion)}
          >
            <Trash2 aria-hidden />
          </Button>
        </div>

        <Collapsible.Content className="data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:animate-in data-[state=open]:fade-in-0">
          <p className="rounded-md bg-muted p-3 text-sm leading-relaxed whitespace-pre-wrap break-words">
            {suggestion.text}
          </p>
        </Collapsible.Content>

        <div className="mt-auto flex items-center justify-between gap-2">
          <RelativeTime
            timestamp={suggestion.timestamp}
            className="text-xs text-muted-foreground"
          />
          <Collapsible.Trigger asChild>
            <Button variant="ghost" size="sm" className="-mr-2 text-accent-foreground">
              {open ? "Hide" : "Show"}
              <span className="sr-only"> suggestion text</span>
              <ChevronDown
                aria-hidden
                className={open ? "rotate-180 transition-transform" : "transition-transform"}
              />
            </Button>
          </Collapsible.Trigger>
        </div>
      </article>
    </Collapsible.Root>
  );
}
