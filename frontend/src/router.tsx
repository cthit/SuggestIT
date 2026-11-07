import { createRootRoute, createRoute, createRouter, Outlet } from "@tanstack/react-router";
import { Failure } from "@/components/feedback";
import { SiteHeader } from "@/components/site-header";
import { CallbackPage } from "@/pages/callback";
import { HomePage } from "@/pages/home";
import { NotFoundPage } from "@/pages/not-found";
import { SuggestionPage } from "@/pages/suggestion";

function Layout() {
  return (
    <>
      <a
        href="#content"
        className="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:rounded-md focus:bg-white focus:p-3"
      >
        Skip to content
      </a>
      <SiteHeader />
      <main id="content" className="shell app-shell">
        <Outlet />
      </main>
    </>
  );
}

const rootRoute = createRootRoute({
  component: Layout,
  notFoundComponent: NotFoundPage,
  errorComponent: ({ error, reset }) => (
    <Failure error={error instanceof Error ? error : new Error(String(error))} retry={reset} />
  ),
});

const homeRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  validateSearch: (search: Record<string, unknown>): { login?: "failed" } =>
    search.login === "failed" ? { login: "failed" } : {},
  component: HomePage,
});

const suggestionRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/suggestion/$id",
  component: SuggestionPage,
});

const callbackRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/callback",
  component: CallbackPage,
});

export const router = createRouter({
  scrollRestoration: true,
  routeTree: rootRoute.addChildren([homeRoute, suggestionRoute, callbackRoute]),
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
