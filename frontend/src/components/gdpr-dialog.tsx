import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

export function GdprNotice() {
  return (
    <Dialog>
      <p className="text-sm text-muted-foreground">
        By clicking Send I agree to{" "}
        <DialogTrigger className="font-medium text-accent-foreground underline underline-offset-4 hover:no-underline rounded-sm">
          this GDPR agreement
        </DialogTrigger>
      </p>
      <DialogContent className="max-h-[calc(100dvh-2rem)] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>GDPR agreement</DialogTitle>
          <DialogDescription className="sr-only">
            How the data you submit with a suggestion is handled
          </DialogDescription>
        </DialogHeader>
        <p className="text-sm leading-relaxed">
          The data is collected by the Software Engineering Student Division (IT) with the aim of
          knowing how to improve Hubben 2.2. The data that will be saved is title, the suggestion,
          the time it was sent and name if not anonymous. You have the right to withdraw your
          consent to the management, request to have all data that the organization has about you
          and to complain to the Data Inspectorate in case of dissatisfaction. The data will be kept
          until P.R.I.T. has solved the issue or deemed it invalid and will only be handled by the
          Student Division Information Technology, ultimately responsible for the data management of
          the section is William Levén, who can be reached at{" "}
          <a className="text-accent-foreground underline" href="mailto:ordforande@chalmers.it">
            ordforande@chalmers.it
          </a>
          . If you have any questions, please contact P.R.I.T.{" "}
          <a className="text-accent-foreground underline" href="mailto:prit@chalmers.it">
            prit@chalmers.it
          </a>{" "}
          or the section's data protection representative{" "}
          <a className="text-accent-foreground underline" href="mailto:dpo@chalmers.it">
            dpo@chalmers.it
          </a>
          .
        </p>
        <DialogFooter showCloseButton />
      </DialogContent>
    </Dialog>
  );
}
