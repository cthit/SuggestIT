import { Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

export function AboutDialog() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="ghost" size="sm" aria-label="About SuggestIT">
          <Info aria-hidden />
          <span className="hidden sm:inline">About</span>
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>About suggestIT</DialogTitle>
          <DialogDescription className="sr-only">
            What SuggestIT is and how it works
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 text-sm leading-relaxed">
          <p>
            Do you have anything you want to discuss about Hubben 2.2 or the Software Engineering
            Student Division (IT)? Then you can write a suggestion on this page!
          </p>
          <p>
            All suggestions will be read at P.R.I.T.'s lunch meetings and those not to P.R.I.T. are
            forwarded to the right committee. If so, write which committee you want to reach.
          </p>
          <p>
            PS. Do you know that all committees and societies have their own email address that you
            can turn to for suggestions and opinions?
          </p>
        </div>
        <DialogFooter showCloseButton />
      </DialogContent>
    </Dialog>
  );
}
