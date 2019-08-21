import React, {Component} from 'react';
import { Dialog, DialogContent, DialogActions } from '@material-ui/core';
import { DigitButton, DigitText } from '@cthit/react-digit-components';
import './ConfirmModal.css';

export class ConfirmModal extends Component{
  
  constructor(props){
    super(props);

    this.state = {
      open: props.open,
      onConfirm: props.onConfirm,
      onClose: props.onClose,
      text: props.text,
      title: props.title ? props.title : 'Är du säker?'
    }
  }

  render = () =>
    <Dialog open={this.state.open} aria-labelledby="form-dialog-title">
      <DigitText.Title className="dialog-title" text={this.state.title}/>
      <DialogContent>
        {this.state.text}
      </DialogContent>
      <DialogActions className="button-collection">
        <DigitButton onClick = {this.state.onClose}   text="Nej" primary/>
        <DigitButton onClick = {this.state.onConfirm} text="Ja" primary/>
      </DialogActions>
    </Dialog>
  
}

export default ConfirmModal;