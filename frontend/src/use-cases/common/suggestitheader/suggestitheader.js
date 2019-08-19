import React, { Component } from "react";
import {DigitHeader, DigitButton} from '@cthit/react-digit-components';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import axios from "axios";
import Cookies from 'universal-cookie';

const cookies = new Cookies();

class SuggestITHeader extends Component {
    constructor(props){
        super(props);
        this.state = {
            renderMain: props.renderMain,
            isLoggedIn: false,
            loginOpen: false,
            passTextField: ''
        }

        this.checkLogin();
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
                      onChange = {(event) => this.setState({passTextField: event.target.value})}
                    />
                  </DialogContent>
                  <DialogActions>
                    <Button onClick={this.handleClose} color="primary">
                      Cancel
                    </Button>
                    <Button onClick={this.login} color="primary">
                      Login
                    </Button>
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

  login = () => {
    axios.put('http://localhost:5000/authenticate', {
      password: this.state.passTextField
    }).then(res => {
          cookies.set('PRIT_AUTH_KEY', res.data.key);
          this.handleClose();
          this.setState({
            isLoggedIn: true
          });
        }
      ).catch(res => 
        console.log("Fail")
      );
  }

  checkLogin = () => {
    if(cookies.get('PRIT_AUTH_KEY'))
      axios.get("http://localhost:5000/authenticate", { headers: {
        Authorization: cookies.get('PRIT_AUTH_KEY')
      }}).then(res =>
          {
            console.log("You are logged in")
            this.setState({
            isLoggedIn: true
          })
        }
        ).catch(res => 
          this.setState({
            isLoggedIn: false
          })  
        )
    else
      this.setState({isLoggedIn: false})
  }

  handleClickOpen = () =>
    this.setState({
      loginOpen: true
    })
  
}

export default SuggestITHeader;
