import React from "react";
import { Dialog, DialogContent, DialogActions } from "@material-ui/core";
import { DigitButton, DigitText } from "@cthit/react-digit-components";
import "./ConfirmModal.css";

export const ConfirmModal = ({ open, onConfirm, onClose, text, title }) => (
    <Dialog open={open} aria-labelledby="form-dialog-title">
        <DigitText.Title
            className="dialog-title"
            text={title ? title : "Är du säker?"}
        />
        <DialogContent>
            <DigitText.Subtitle text={text} />
        </DialogContent>
        <DialogActions className="button-collection">
            <DigitButton onClick={onClose} text="Nej" primary />
            <DigitButton onClick={onConfirm} text="Ja" primary />
        </DialogActions>
    </Dialog>
);

export default ConfirmModal;
