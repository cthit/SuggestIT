import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "@tanstack/react-router";
import { ArrowLeft, LogIn, UserRound } from "lucide-react";
import { Failure, Loading } from "@/components/feedback";
import { RelativeTime } from "@/components/relative-time";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { isApiError, LOGIN_URL } from "@/lib/api";
import { suggestionQuery, useMe } from "@/lib/queries";

export function SuggestionPage() {
  const { id } = useParams({ from: "/suggestion/$id" });
  const me = useMe();

  return (
    <div className="mx-auto grid max-w-3xl gap-6">
      <Link
        to="/"
        className="flex w-fit items-center gap-2 rounded-sm text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" aria-hidden />
        Back to suggestions
      </Link>
      {me.isPending ? (
        <Loading />
      ) : me.isError ? (
        <Failure error={me.error} retry={() => void me.refetch()} />
      ) : me.data ? (
        <SuggestionDetails id={id} />
      ) : (
        <Notice title="Log in to view suggestions">
          <Button asChild>
            <a href={LOGIN_URL}>
              <LogIn aria-hidden />
              Login
            </a>
          </Button>
        </Notice>
      )}
    </div>
  );
}

function SuggestionDetails({ id }: { id: string }) {
  const query = useQuery(suggestionQuery(id));

  if (query.isPending) {
    return <Loading />;
  }

  if (query.isError) {
    if (isApiError(query.error, 404)) {
      return <Notice title="This suggestion does not exist. It may have been deleted." />;
    }

    return <Failure error={query.error} retry={() => void query.refetch()} />;
  }

  const suggestion = query.data;

  return (
    <Card className="gap-5">
      <CardHeader className="gap-3 border-b pb-5">
        <h1 className="text-2xl leading-tight font-bold tracking-tight break-words">
          {suggestion.title}
        </h1>
        <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-muted-foreground">
          <span className="flex items-center gap-1.5">
            <UserRound className="size-4" aria-hidden />
            Posted by: {suggestion.author}
          </span>
          <RelativeTime timestamp={suggestion.timestamp} />
        </div>
      </CardHeader>
      <CardContent>
        <p className="leading-relaxed whitespace-pre-wrap break-words">{suggestion.text}</p>
      </CardContent>
    </Card>
  );
}

function Notice({ title, children }: { title: string; children?: ReactNode }) {
  return (
    <div className="flex flex-col items-center gap-4 rounded-xl border border-dashed bg-card/60 px-6 py-12 text-center">
      <p className="font-medium">{title}</p>
      {children}
    </div>
  );
}
