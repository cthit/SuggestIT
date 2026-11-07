import { useEffect } from "react";
import { Loading } from "@/components/feedback";

/**
 * Gamma redirects here (the registered redirect URI). Hand the authorization code over to the
 * backend, which finishes the login and sends the browser back to "/".
 */
export function CallbackPage() {
  useEffect(() => {
    window.location.replace(`/api/auth/callback${window.location.search}`);
  }, []);

  return <Loading label="Logging in…" />;
}
