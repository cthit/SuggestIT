import { useRef } from "react";
import { Inbox, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { Failure, Loading } from "@/components/feedback";
import { SuggestionCard } from "@/components/suggestion-card";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import type { Suggestion } from "@/lib/api";
import { notifyError } from "@/lib/notify";
import {
  suggestionsQuery,
  useClearSuggestions,
  useCreateSuggestion,
  useDeleteSuggestion,
} from "@/lib/queries";
import { useQuery } from "@tanstack/react-query";

export function SuggestionBoard() {
  const headingRef = useRef<HTMLHeadingElement>(null);
  const query = useQuery(suggestionsQuery);
  const remove = useDeleteSuggestion();
  const restore = useCreateSuggestion();
  const clear = useClearSuggestions();
  const suggestions = query.data ?? [];

  function deleteWithUndo(suggestion: Suggestion) {
    moveFocusAwayFrom(suggestion);
    // mutateAsync (not mutate) so every call gets its own toast when deleting several quickly.
    remove.mutateAsync(suggestion).then(
      () =>
        toast("The suggestion has been deleted", {
          action: {
            label: "Undo",
            onClick: () =>
              void restore
                .mutateAsync({
                  title: suggestion.title,
                  text: suggestion.text,
                  author: suggestion.author,
                })
                .catch(notifyError),
          },
        }),
      notifyError,
    );
  }

  // The card disappears as soon as the delete starts. Move focus to the next
  // card's delete button (or the previous one, or the heading) so keyboard and
  // screen reader users are not sent back to the top of the page.
  function moveFocusAwayFrom(suggestion: Suggestion) {
    const index = suggestions.findIndex((item) => item.id === suggestion.id);
    const neighbour = suggestions[index + 1] ?? suggestions[index - 1];
    const target = neighbour
      ? document.querySelector<HTMLElement>(
          `[data-delete-suggestion="${CSS.escape(neighbour.id)}"]`,
        )
      : null;

    (target ?? headingRef.current)?.focus();
  }

  function clearAll() {
    const count = suggestions.length;

    clear.mutate(
      suggestions.map((suggestion) => suggestion.id),
      {
        onSuccess: () => toast.success(`Deleted ${count} ${plural(count)}`),
        onError: notifyError,
      },
    );
  }

  return (
    <section aria-labelledby="suggestions-heading" className="grid gap-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h2
          ref={headingRef}
          id="suggestions-heading"
          tabIndex={-1}
          className="flex items-center gap-2.5 rounded-sm text-xl font-semibold outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50"
        >
          Suggestions
          {query.isSuccess && (
            <span className="rounded-full bg-secondary px-2.5 py-0.5 text-sm font-medium text-secondary-foreground tabular-nums">
              {suggestions.length}
            </span>
          )}
        </h2>
        {suggestions.length > 0 && (
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="outline" size="sm" disabled={clear.isPending}>
                <Trash2 aria-hidden />
                Clear suggestions
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                <AlertDialogDescription>
                  Are you sure you want to delete {suggestions.length} {plural(suggestions.length)}?
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>No</AlertDialogCancel>
                <AlertDialogAction
                  className={buttonVariants({ variant: "destructive" })}
                  onClick={clearAll}
                >
                  Yes
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        )}
      </div>

      {query.isPending ? (
        <Loading label="Loading suggestions…" />
      ) : query.isError ? (
        <Failure error={query.error} retry={() => void query.refetch()} />
      ) : suggestions.length === 0 ? (
        <div className="flex flex-col items-center gap-3 rounded-xl border border-dashed bg-card/60 px-6 py-14 text-center">
          <span className="grid size-12 place-items-center rounded-full bg-secondary text-secondary-foreground">
            <Inbox className="size-6" aria-hidden />
          </span>
          <p className="font-medium">The suggestion box is empty.</p>
        </div>
      ) : (
        <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {suggestions.map((suggestion) => (
            <li key={suggestion.id} className="grid">
              <SuggestionCard suggestion={suggestion} onDelete={deleteWithUndo} />
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}

function plural(count: number) {
  return count === 1 ? "suggestion" : "suggestions";
}
