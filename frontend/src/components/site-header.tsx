import { Link } from "@tanstack/react-router";
import { Lightbulb, LogIn, LogOut } from "lucide-react";
import { AboutDialog } from "@/components/about-dialog";
import { Button } from "@/components/ui/button";
import { LOGIN_URL } from "@/lib/api";
import { notifyError } from "@/lib/notify";
import { useLogout, useMe } from "@/lib/queries";

export function SiteHeader() {
  const me = useMe();
  const logout = useLogout();
  const user = me.data;

  return (
    <header className="sticky top-0 z-40 border-b bg-white/90 backdrop-blur supports-[backdrop-filter]:bg-white/75">
      <div className="shell flex min-h-16 items-center justify-between gap-3 sm:min-h-20">
        <Link to="/" aria-label="SuggestIT home" className="flex items-center gap-2.5 rounded-lg">
          <span className="grid size-9 place-items-center rounded-lg bg-brand text-primary-foreground">
            <Lightbulb className="size-5" aria-hidden />
          </span>
          <span className="text-xl font-bold tracking-tight sm:text-2xl">
            Suggest<span className="text-accent-foreground">IT</span>
          </span>
        </Link>
        <div className="flex items-center gap-1 sm:gap-2">
          <AboutDialog />
          {me.isPending ? (
            <span className="h-8 w-20" aria-hidden />
          ) : user ? (
            <>
              <span className="ml-2 hidden border-l pl-4 text-sm md:inline">
                {user.nickname || user.name || user.cid}
              </span>
              <Button
                variant="ghost"
                size="sm"
                disabled={logout.isPending}
                onClick={() =>
                  logout.mutate(undefined, {
                    onError: notifyError,
                  })
                }
              >
                <LogOut aria-hidden />
                Logout
              </Button>
            </>
          ) : (
            <Button variant="outline" size="sm" asChild>
              <a href={LOGIN_URL}>
                <LogIn aria-hidden />
                Login
              </a>
            </Button>
          )}
        </div>
      </div>
    </header>
  );
}
