import { toast } from "sonner";
import { isApiError } from "@/lib/api";

/** Shows an error toast, except for 401s which main.tsx already reports as an expired session. */
export function notifyError(error: Error) {
  if (!isApiError(error, 401)) {
    toast.error(error.message);
  }
}
