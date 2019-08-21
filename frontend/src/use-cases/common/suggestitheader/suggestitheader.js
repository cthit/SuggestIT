import React, { Component } from "react";
import {DigitHeader, DigitButton} from '@cthit/react-digit-components';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import { login, checkLogin } from "../../../services/data.service";

class SuggestITHeader extends Component {
    constructor(props){
        super(props);
        this.state = {
            renderMain: props.renderMain,
            isLoggedIn: false,
            loginOpen: false,
            passTextField: 'Hello',
            passError: false
        }

        checkLogin().then(res=>
            this.setState({
            isLoggedIn: true
          })
        )
        .catch(error => {
          //User might not be logged in
        })
    }

  render() {
    return (
      <DigitHeader
        renderMain = {this.state.renderMain}
        title = "SuggestIT"
        renderHeader = {() => !this.state.isLoggedIn ?
              <div>
                <DigitButton text="Login" outlined onClick={this.handleClickOpen}/>
                <Dialog open={this.state.loginOpen} onClose={this.handleClose} aria-labelledby="form-dialog-title">
                  <DialogTitle id="form-dialog-title">P.R.I.T login</DialogTitle>
                  <DialogContent>
                    <TextField
                      autoFocus
                      margin="dense"
                      id="name"
                      label="Password"
                      type="password"
                      fullWidth
                      error= {this.state.passError}
                      onChange = {(event) => this.setState({passTextField: event.target.value})}
                      onKeyPress = {event => {
                        if(event.key === 'Enter')
                          this.login()
                      }}
                      //Why not using DigitTextField, 1. No onKeyPress event, 2. No autofocus
                    />
                  </DialogContent>
                  <DialogActions>
                    <DigitButton onClick={this.handleClose} text="Cancel" primary/>
                    <DigitButton onClick={this.login} text="Login" primary/>
                  </DialogActions>
                </Dialog>
              </div>
              
          : <dir></dir>
        }/>
    );
  }

  handleClose = () => 
    this.setState({
      loginOpen: false,
      passTextField: ''
    })

  handleClickOpen = () =>
    this.setState({
      loginOpen: true
    })

  login = () => {
    login(this.state.passTextField)
    .then(res => {
      this.handleClose();
      this.setState({
        isLoggedIn: true
      });
      window.location.reload(false);
    })
    .catch(error => 
      this.setState({
        passError: true
      })
    );
  }  
}

export default SuggestITHeader;
