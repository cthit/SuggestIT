import React, { Component } from "react";
import {
    DigitHeader,
    DigitButton,
    DigitDialogActions,
    DigitLayout,
    DigitText,
} from "@cthit/react-digit-components";
import { connect } from "react-redux";
import About from "./elements/about";
import { checkLogin, logOut, loginRedirect } from "../../services/data.service";

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
                if (error == undefined) return;

                props.dialogOpen({
                    title: "",
                    renderMain: () => <DigitText.Text text={error} />,
                    renderButtons: (confirm, cancel) => (
                        <DigitButton text="Close" onClick={cancel} />
                    ),
                    onConfirm: (confirm, reject) => confirm(),
                });
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
                                        onClick={() => loginRedirect()}
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
