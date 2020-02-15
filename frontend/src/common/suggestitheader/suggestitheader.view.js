import React, { Component } from "react";
import {
    DigitHeader,
    DigitButton,
    DigitDialogActions,
    DigitLayout,
} from "@cthit/react-digit-components";
import { connect } from "react-redux";
import About from "./elements/about";
import { checkLogin, logOut } from "../../services/data.service";

const auth_url = process.env.REACT_APP_LDAP_AUTH_URL;
const client_id = process.env.REACT_APP_CLIENT_ID;

class SuggestITHeaderView extends Component {
    constructor(props) {
        super(props);
        this.state = {
            renderMain: props.renderMain,
            isLoggedIn: false,
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
                                        onConfirm: (confirm, reject) =>
                                            confirm(),
                                    })
                                }
                            />
                            {!this.state.isLoggedIn ? (
                                <div>
                                    <DigitButton
                                        text="Login"
                                        outlined
                                        onClick={() =>
                                            window.location.replace(
                                                `${auth_url}/authenticate?client_id=${client_id}`
                                            )
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
