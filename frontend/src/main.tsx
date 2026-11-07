import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider } from "@tanstack/react-router";
import { Toaster, toast } from "sonner";
import { isApiError } from "@/lib/api";
import { meQuery } from "@/lib/queries";
import { router } from "./router";
import "./styles.css";

// A 401 from any protected endpoint means the session cookie expired: forget the user so the
// suggestion board disappears instead of showing a stale list.
function onError(error: unknown) {
  if (isApiError(error, 401) && client.getQueryData(meQuery.queryKey)) {
    client.setQueryData(meQuery.queryKey, null);
    toast.error("Your session has expired. Please log in again.", { id: "session-expired" });
  }
}

const client = new QueryClient({
  queryCache: new QueryCache({ onError }),
  mutationCache: new MutationCache({ onError }),
  defaultOptions: {
    queries: {
      staleTime: 20_000,
      // Retry once on network/server errors; 4xx answers will not change on a retry.
      retry: (count, error) =>
        count < 1 && (!isApiError(error) || error.status === 0 || error.status >= 500),
    },
  },
});

const root = document.getElementById("root");

if (!root) {
  throw new Error("Application root is missing");
}

createRoot(root).render(
  <StrictMode>
    <QueryClientProvider client={client}>
      <RouterProvider router={router} />
      <Toaster richColors closeButton position="bottom-right" duration={5000} />
    </QueryClientProvider>
  </StrictMode>,
);
