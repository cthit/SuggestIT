import React, { Component } from "react";
import {
    DigitHeader,
    DigitButton,
    DigitDialogActions,
    DigitLayout,
} from "@cthit/react-digit-components";
import TextField from "@material-ui/core/TextField";
import { connect } from "react-redux";
import About from "./elements/about";
import { login, checkLogin, logOut } from "../../services/data.service";
import { ErrorBox } from "./suggestitheader.style.js";

class SuggestITHeaderView extends Component {
    constructor(props) {
        super(props);
        this.state = {
            renderMain: props.renderMain,
            isLoggedIn: false,
            cidTextField: "",
            passTextField: "",
            loginErrorMessage: "",
            dialogOpen: props.dialogOpen,
        };

        checkLogin()
            .then(res =>
                this.setState({
                    isLoggedIn: true,
                })
            )
            .catch(error => {
                //User might not be logged in
            });
    }

    render() {
        return (
            <DigitHeader
                renderMain={this.state.renderMain}
                title="SuggestIT"
                renderHeader={() => (
                    <div>
                        <DigitLayout.Row>
                            <DigitButton
                                text="About"
                                onClick={values =>
                                    this.state.dialogOpen({
                                        title: "About suggestIT",
                                        renderMain: () => <About />,
                                        renderButtons: (confirm, cancel) => (
                                            <DigitButton
                                                text="Close"
                                                onClick={cancel}
                                            />
                                        ),
                                    })
                                }
                            />
                            {!this.state.isLoggedIn ? (
                                <div>
                                    <DigitButton
                                        text="Login"
                                        outlined
                                        onClick={values =>
                                            this.state.dialogOpen({
                                                title:
                                                    "Enter CID and password for chalmers.it",
                                                renderMain: () => (
                                                    <div>
                                                        {this.state
                                                            .loginErrorMessage ===
                                                        "" ? null : (
                                                            <ErrorBox>
                                                                {
                                                                    this.state
                                                                        .loginErrorMessage
                                                                }
                                                            </ErrorBox>
                                                        )}
                                                        {this.cidTextField()}
                                                        {this.passTextField()}
                                                    </div>
                                                ),
                                                renderButtons: (
                                                    confirm,
                                                    cancel
                                                ) => (
                                                    <>
                                                        <DigitButton
                                                            text={"cancel"}
                                                            onClick={cancel}
                                                        />
                                                        <DigitButton
                                                            text={"Confirm"}
                                                            onClick={confirm}
                                                        />
                                                    </>
                                                ),
                                                onCancel: e => this.login(),
                                                onConfirm: e =>
                                                    this.handleClose(),
                                            })
                                        }
                                    />
                                </div>
                            ) : (
                                <DigitButton
                                    outlined
                                    text="logout"
                                    onClick={() => {
                                        logOut();
                                        window.location.reload(false);
                                    }}
                                />
                            )}
                        </DigitLayout.Row>
                    </div>
                )}
            />
        );
    }

    cidTextField = () => (
        <TextField
            autoFocus
            margin="dense"
            id="cid"
            label="CID"
            type="username"
            fullWidth
            onChange={event =>
                this.setState({
                    cidTextField: event.target.value,
                })
            }
            onKeyPress={event => {
                if (event.key === "Enter") this.login();
            }}
        />
    );

    passTextField = () => (
        <TextField
            margin="dense"
            id="password"
            label="Password"
            type="password"
            fullWidth
            onChange={event =>
                this.setState({
                    passTextField: event.target.value,
                })
            }
            onKeyPress={event => {
                if (event.key === "Enter") this.login();
            }}
        />
    );

    handleClose = () =>
        this.setState({
            cidTextField: "",
            passTextField: "",
            loginErrorMessage: "",
        });

    login = () => {
        login(this.state.cidTextField, this.state.passTextField)
            .then(res => {
                this.handleClose();
                this.setState({
                    isLoggedIn: true,
                });
                window.location.reload(false);
            })
            .catch(error =>
                this.setState({
                    loginErrorMessage: error,
                })
            );
    };
}

const mapStateToProps = (state, ownProps) => ({});

const mapDispatchToProps = dispatch => ({
    dialogOpen: dialogData =>
        dispatch(DigitDialogActions.digitDialogCustomOpen(dialogData)),
});

const SuggestITHeader = connect(
    mapStateToProps,
    mapDispatchToProps
)(SuggestITHeaderView);

export default SuggestITHeader;
