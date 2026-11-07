export interface Suggestion {
  id: string;
  /** RFC 3339 timestamp in UTC. */
  timestamp: string;
  title: string;
  text: string;
  author: string;
}

export type NewSuggestion = Pick<Suggestion, "title" | "text" | "author">;

export interface Me {
  cid: string;
  nickname: string;
  name: string;
  picture: string;
}

export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export function isApiError(error: unknown, status?: number): error is ApiError {
  return error instanceof ApiError && (status === undefined || error.status === status);
}

async function request<T>(
  path: string,
  init: { method?: string; body?: unknown } = {},
): Promise<T> {
  const hasBody = init.body !== undefined;
  let response: Response;

  try {
    response = await fetch(path, {
      method: init.method ?? "GET",
      credentials: "same-origin",
      headers: {
        Accept: "application/json",
        ...(hasBody ? { "Content-Type": "application/json" } : {}),
      },
      body: hasBody ? JSON.stringify(init.body) : undefined,
    });
  } catch {
    throw new ApiError(0, "Could not reach SuggestIT. Check your connection and try again.");
  }

  if (!response.ok) {
    throw new ApiError(response.status, await errorMessage(response));
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

async function errorMessage(response: Response): Promise<string> {
  try {
    const body = (await response.json()) as { error?: unknown };

    if (typeof body.error === "string" && body.error) {
      return body.error;
    }
  } catch {
    // Not a JSON error body; fall through to a generic message.
  }

  return `Request failed (${response.status}). Please try again.`;
}

export const api = {
  /** The signed-in user, or `null` when nobody is signed in. */
  me: async (): Promise<Me | null> => {
    try {
      return await request<Me>("/api/me");
    } catch (error) {
      if (isApiError(error, 401)) {
        return null;
      }

      throw error;
    }
  },
  logout: () => request<void>("/api/auth/logout", { method: "POST" }),
  suggestions: () => request<Suggestion[]>("/api/suggestions"),
  suggestion: (id: string) => request<Suggestion>(`/api/suggestions/${encodeURIComponent(id)}`),
  createSuggestion: (suggestion: NewSuggestion) =>
    request<Suggestion>("/api/suggestions", { method: "POST", body: suggestion }),
  deleteSuggestion: (id: string) =>
    request<void>(`/api/suggestions/${encodeURIComponent(id)}`, { method: "DELETE" }),
  deleteSuggestions: (ids: string[]) =>
    request<void>("/api/suggestions/delete", { method: "POST", body: { ids } }),
};

export const LOGIN_URL = "/api/auth/login";
