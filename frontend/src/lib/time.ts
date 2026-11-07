import { useSyncExternalStore } from "react";

const units: [Intl.RelativeTimeFormatUnit, number][] = [
  ["second", 60],
  ["minute", 60],
  ["hour", 24],
  ["day", 365.25 / 12],
  ["month", 12],
  ["year", Number.POSITIVE_INFINITY],
];

const relative = new Intl.RelativeTimeFormat("en", { numeric: "always" });

/** Formats a timestamp as e.g. "Posted 5 minutes ago". */
export function postedAgo(timestamp: string, now: number = Date.now()): string {
  let value = Math.max(0, (now - new Date(timestamp).getTime()) / 1000);

  if (value < 10) {
    return "Posted just now";
  }

  for (const [unit, size] of units) {
    if (value < size) {
      return `Posted ${relative.format(-Math.floor(value), unit)}`;
    }

    value /= size;
  }

  return "Posted a long time ago";
}

const absolute = new Intl.DateTimeFormat("en-GB", { dateStyle: "long", timeStyle: "short" });

export function formatDateTime(timestamp: string): string {
  return absolute.format(new Date(timestamp));
}

// A single shared clock so every relative timestamp re-renders together.
const listeners = new Set<() => void>();
let now = Date.now();
let timer: ReturnType<typeof setInterval> | undefined;

function subscribe(listener: () => void) {
  listeners.add(listener);

  if (timer === undefined) {
    now = Date.now();
    timer = setInterval(() => {
      now = Date.now();
      for (const notify of listeners) {
        notify();
      }
    }, 30_000);
  }

  return () => {
    listeners.delete(listener);

    if (listeners.size === 0 && timer !== undefined) {
      clearInterval(timer);
      timer = undefined;
    }
  };
}

export function useNow(): number {
  return useSyncExternalStore(subscribe, () => now);
}
