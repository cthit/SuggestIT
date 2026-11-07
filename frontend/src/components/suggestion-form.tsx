import { useId, useRef, useState, type FormEvent, type ReactNode } from "react";
import { LoaderCircle, Send } from "lucide-react";
import { toast } from "sonner";
import { GdprNotice } from "@/components/gdpr-dialog";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { notifyError } from "@/lib/notify";
import { useCreateSuggestion } from "@/lib/queries";

const MAX_LENGTH = 24;
const TITLE_ERROR = "The title is not filled in";
const DESCRIPTION_ERROR = "The description is not filled in";

interface Errors {
  title?: string;
  text?: string;
}

export function SuggestionForm() {
  const ids = useId();
  const titleRef = useRef<HTMLInputElement>(null);
  const textRef = useRef<HTMLTextAreaElement>(null);
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");
  const [author, setAuthor] = useState("");
  const [anonymous, setAnonymous] = useState(false);
  const [errors, setErrors] = useState<Errors>({});
  const create = useCreateSuggestion();

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const next: Errors = {
      title: title.trim() ? undefined : TITLE_ERROR,
      text: text.trim() ? undefined : DESCRIPTION_ERROR,
    };
    setErrors(next);

    if (next.title || next.text) {
      (next.title ? titleRef : textRef).current?.focus();

      return;
    }

    create.mutate(
      {
        title: title.trim(),
        text: text.trim(),
        author: anonymous || !author.trim() ? "Anonymous" : author.trim(),
      },
      {
        onSuccess: () => {
          setTitle("");
          setText("");
          setAuthor("");
          toast.success("Thank you! The suggestion has been sent to P.R.I.T.");
        },
        onError: notifyError,
      },
    );
  }

  return (
    <Card className="gap-5">
      <CardHeader>
        <CardTitle>
          <h1 className="text-xl tracking-tight">New suggestion</h1>
        </CardTitle>
        <CardDescription>
          Ideas, complaints or questions about Hubben 2.2 or the IT student division go straight to
          P.R.I.T.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form noValidate onSubmit={submit} className="grid gap-5">
          <Field
            id={`${ids}-title`}
            label="Title"
            error={errors.title}
            counter={`${title.length}/${MAX_LENGTH}`}
          >
            {(describedBy) => (
              <Input
                ref={titleRef}
                id={`${ids}-title`}
                name="title"
                autoComplete="off"
                maxLength={MAX_LENGTH}
                required
                value={title}
                aria-invalid={errors.title ? true : undefined}
                aria-describedby={describedBy}
                onChange={(event) => {
                  setTitle(event.target.value);
                  setErrors((current) => ({ ...current, title: undefined }));
                }}
              />
            )}
          </Field>

          <Field id={`${ids}-text`} label="Suggestion" error={errors.text}>
            {(describedBy) => (
              <Textarea
                ref={textRef}
                id={`${ids}-text`}
                name="text"
                rows={5}
                required
                className="min-h-32 max-h-80"
                value={text}
                aria-invalid={errors.text ? true : undefined}
                aria-describedby={describedBy}
                onChange={(event) => {
                  setText(event.target.value);
                  setErrors((current) => ({ ...current, text: undefined }));
                }}
              />
            )}
          </Field>

          <Field
            id={`${ids}-author`}
            label="Name/Nick"
            counter={anonymous ? undefined : `${author.length}/${MAX_LENGTH}`}
          >
            {(describedBy) => (
              <Input
                id={`${ids}-author`}
                name="author"
                autoComplete="nickname"
                maxLength={MAX_LENGTH}
                placeholder={anonymous ? "Anonymous" : "Optional"}
                disabled={anonymous}
                value={anonymous ? "" : author}
                aria-describedby={describedBy}
                onChange={(event) => setAuthor(event.target.value)}
              />
            )}
          </Field>

          <div className="flex items-center gap-3">
            <Switch id={`${ids}-anonymous`} checked={anonymous} onCheckedChange={setAnonymous} />
            <Label htmlFor={`${ids}-anonymous`}>Anonymous</Label>
          </div>

          <div className="flex flex-col-reverse gap-4 border-t pt-5 sm:flex-row sm:items-center sm:justify-between">
            <GdprNotice />
            <Button type="submit" size="lg" disabled={create.isPending} className="sm:min-w-32">
              {create.isPending ? (
                <LoaderCircle className="animate-spin" aria-hidden />
              ) : (
                <Send aria-hidden />
              )}
              Send
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

function Field({
  id,
  label,
  error,
  counter,
  children,
}: {
  id: string;
  label: string;
  error?: string;
  counter?: string;
  children: (describedBy: string | undefined) => ReactNode;
}) {
  const errorId = `${id}-error`;
  const counterId = `${id}-counter`;
  const describedBy =
    [error ? errorId : undefined, counter ? counterId : undefined].filter(Boolean).join(" ") ||
    undefined;

  return (
    <div className="field">
      <div className="flex items-baseline justify-between gap-3">
        <Label htmlFor={id}>{label}</Label>
        {counter && (
          <span id={counterId} className="text-xs tabular-nums text-muted-foreground">
            <span className="sr-only">Characters used: </span>
            {counter}
          </span>
        )}
      </div>
      {children(describedBy)}
      {error && (
        <p id={errorId} className="text-sm font-medium text-destructive">
          {error}
        </p>
      )}
    </div>
  );
}
