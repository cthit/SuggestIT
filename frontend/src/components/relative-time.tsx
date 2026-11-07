import { formatDateTime, postedAgo, useNow } from "@/lib/time";
import { cn } from "@/lib/utils";

export function RelativeTime({ timestamp, className }: { timestamp: string; className?: string }) {
  const now = useNow();

  return (
    <time dateTime={timestamp} title={formatDateTime(timestamp)} className={cn(className)}>
      {postedAgo(timestamp, now)}
    </time>
  );
}
