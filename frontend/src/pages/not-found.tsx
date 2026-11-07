import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";

export function NotFoundPage() {
  return (
    <div className="grid justify-items-center gap-4 py-16 text-center">
      <p className="text-6xl font-bold tracking-tight text-primary">404</p>
      <h1 className="page-title">Site not found</h1>
      <p className="text-muted-foreground">The page you are looking for does not exist.</p>
      <Button asChild variant="outline">
        <Link to="/">Back to SuggestIT</Link>
      </Button>
    </div>
  );
}
