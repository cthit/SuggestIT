import { queryOptions, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api, type NewSuggestion, type Suggestion } from "@/lib/api";

export const meQuery = queryOptions({
  queryKey: ["me"],
  queryFn: api.me,
  staleTime: 60_000,
  retry: false,
});

export const suggestionsQuery = queryOptions({
  queryKey: ["suggestions"],
  queryFn: api.suggestions,
});

export function suggestionQuery(id: string) {
  return queryOptions({
    queryKey: ["suggestions", id],
    queryFn: () => api.suggestion(id),
  });
}

export function useMe() {
  return useQuery(meQuery);
}

export function useLogout() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: api.logout,
    onSuccess: async () => {
      client.removeQueries({ queryKey: suggestionsQuery.queryKey });
      await client.invalidateQueries({ queryKey: meQuery.queryKey });
    },
  });
}

export function useCreateSuggestion() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (suggestion: NewSuggestion) => api.createSuggestion(suggestion),
    onSuccess: () => client.invalidateQueries({ queryKey: suggestionsQuery.queryKey }),
  });
}

const deleteSuggestionKey = ["suggestions", "delete"] as const;

export function useDeleteSuggestion() {
  const client = useQueryClient();
  const key = suggestionsQuery.queryKey;

  return useMutation({
    mutationKey: deleteSuggestionKey,
    mutationFn: (suggestion: Suggestion) => api.deleteSuggestion(suggestion.id),
    // Remove the card straight away; it comes back if the request fails.
    onMutate: async (suggestion) => {
      await client.cancelQueries({ queryKey: key });
      client.setQueryData(key, (list) => list?.filter((item) => item.id !== suggestion.id));
    },
    // Put back only the card that failed. Restoring a snapshot from before this
    // delete would also bring back cards deleted since then.
    onError: (_error, suggestion) => {
      client.setQueryData(key, (list) =>
        list && !list.some((item) => item.id === suggestion.id)
          ? [...list, suggestion].sort(newestFirst)
          : list,
      );
    },
    // Refetch only once the last of several quick deletes has settled, so a
    // refetch cannot return a card whose DELETE is still in flight.
    onSettled: () =>
      client.isMutating({ mutationKey: deleteSuggestionKey }) === 1
        ? client.invalidateQueries({ queryKey: key })
        : undefined,
  });
}

/** Ids per bulk delete request; keeps each body well below the backend's 64 KiB limit. */
const CLEAR_BATCH_SIZE = 500;

export function useClearSuggestions() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: async (ids: string[]) => {
      for (let start = 0; start < ids.length; start += CLEAR_BATCH_SIZE) {
        await api.deleteSuggestions(ids.slice(start, start + CLEAR_BATCH_SIZE));
      }
    },
    onSettled: () => client.invalidateQueries({ queryKey: suggestionsQuery.queryKey }),
  });
}

function newestFirst(a: Suggestion, b: Suggestion) {
  return Date.parse(b.timestamp) - Date.parse(a.timestamp);
}
