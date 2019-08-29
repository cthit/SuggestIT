import React, { Component, useState, useReducer } from "react";
import {
    DigitTextField,
    /*DigitTextArea,*/
    DigitButton,
    DigitSwitch,
    DigitLayout,
    DigitText,
    DigitToastActions,
} from "@cthit/react-digit-components";
import {
    addSuggestion,
    updateSuggestions,
} from "../../../services/data.service";
import { connect } from "react-redux";
import { TextField } from "@material-ui/core";
import "./Prompt.css";

const initSuggestion = {
    title: "",
    text: "",
    author: "",
};

const suggestionChanged = (state, action) => {
    console.log("Setting suggestion");
    return { ...state, ...action.payload };
};

const PromptView = ({ toastOpen }) => {
    const [suggestion, setSuggestion] = useReducer(
        suggestionChanged,
        initSuggestion
    );
    const title_error_message = "The title is not filled in";
    const description_error_message = "The description is not filled in";
    const [errors, setErrors] = useState({
        title_error: false,
        description_error: false,
    });
    const [author, setAuthor] = useState("");
    const [anonymous_author, setAnonymousAuthor] = useState(false);

    const sendNewSuggestion = () => {
        addSuggestion(suggestion).then(res => {
            updateSuggestions();
            setSuggestion({
                title: "",
                text: "",
                author: "",
            });
            toastOpen({
                text: "Thank you! The suggestion has been sent to P.R.I.T.",
                duration: 5000,
            });
        });
    };

    return (
        <div className="prompt">
            <div className="innerPrompt">
                <DigitText.Heading6 text="New Suggestion" />
                <DigitTextField
                    error={errors.title_error}
                    errorMessage={title_error_message}
                    onChange={e =>
                        setSuggestion({ payload: { title: e.target.value } })
                    }
                    value={suggestion.title}
                    upperLabel="Title"
                />
                <br />
                <TextField
                    id="description"
                    type="text"
                    rows={6}
                    fullWidth={true}
                    label="Description"
                    multiline
                    error={errors.description_error}
                    errorMessage={description_error_message}
                    onChange={e =>
                        setSuggestion({ payload: { text: e.target.value } })
                    }
                />
                {
                    //DigitTextArea will be used when it has been updated,
                    //There is currently a bug causing an infinite loop
                    /*<DigitTextArea
                  error={this.state.description_isempty}
                  errorMessage = {this.state.description_error_message}
                  onChange={e => {
                    console.log("Hello")
                  }}
                  value={this.state.description}
                  upperLabel="Förslag"
                  rows={5}
                  rowsMax={10}
                />*/
                }
                <DigitLayout.Row>
                    <DigitTextField
                        onChange={e => setAuthor(e.target.value)}
                        value={author}
                        disabled={anonymous_author}
                        upperLabel="CID"
                    />

                    <DigitSwitch
                        value={anonymous_author}
                        label="Anonymous"
                        primary
                        onChange={e => setAnonymousAuthor(e.target.value)}
                    />
                </DigitLayout.Row>
                <DigitButton
                    text="Send"
                    primary
                    raised
                    onClick={() => {
                        setErrors({
                            title_error: suggestion.title === "",
                            description_error: suggestion.text === "",
                        });

                        if (suggestion.title === "" || suggestion.text === "")
                            return;

                        setAuthor(author === "" ? "Anonymous" : author);
                        sendNewSuggestion();
                    }}
                />
            </div>
        </div>
    );
};

/*class PromptView extends Component {
    constructor(props) {
        super(props);
        this.state = {
            title: "",
            text: "",
            author: "",
            anonymus_author: false,
            title_isempty: false,
            title_error_message: "The title is not filled in",
            description_isempty: false,
            description_error_message: "You need to discribe your suggestion",
            toastOpen: this.props["toastOpen"],
        };
    }

    sendNewSuggestion() {
        addSuggestion({
            title: this.state.title,
            text: this.state.description,
            author:
                this.state.author === "" || this.state.anonymus_author
                    ? "Anonymous"
                    : this.state.author,
        }).then(res => {
            updateSuggestions();
            this.setState({
                title: "",
                text: "",
                author: "",
            });
            this.state.toastOpen({
                text: "Thank you! The suggestion has been sent to P.R.I.T.",
                duration: 5000,
            });
        });
    }

    render() {
        return (
            <div className="prompt">
                <div className="innerPrompt">
                    <DigitText.Heading6 text="New Suggestion" />
                    <DigitTextField
                        error={this.state.title_isempty}
                        errorMessage={this.state.title_error_message}
                        onChange={e => {
                            this.setState({
                                title: e.target.value,
                            });
                        }}
                        value={this.state.title}
                        upperLabel="Title"
                    />
                    <br />
                    <TextField
                        id="description"
                        type="text"
                        rows={6}
                        fullWidth={true}
                        label="Description"
                        multiline
                        error={this.state.description_isempty}
                        errorMessage={this.state.description_error_message}
                        onChange={e =>
                            this.setState({
                                text: e.target.value,
                            })
                        }
                    />
                    {
                        //DigitTextArea will be used when it has been updated,
                        //There is currently a bug causing an infinite loop
                        /*<DigitTextArea
                      error={this.state.description_isempty}
                      errorMessage = {this.state.description_error_message}
                      onChange={e => {
                        console.log("Hello")
                      }}
                      value={this.state.description}
                      upperLabel="Förslag"
                      rows={5}
                      rowsMax={10}
                    />* /
                    }
                    <DigitLayout.Row>
                        <DigitTextField
                            onChange={e => {
                                this.setState({
                                    author: e.target.value,
                                });
                            }}
                            value={this.state.author}
                            disabled={this.state.anonymus_author}
                            upperLabel="CID"
                        />

                        <DigitSwitch
                            value={this.state.anonymus_author}
                            label="Anonymous"
                            primary
                            onChange={e => {
                                this.setState({
                                    anonymus_author: e.target.checked,
                                });
                            }}
                        />
                    </DigitLayout.Row>
                    <DigitButton
                        text="Send"
                        primary
                        raised
                        onClick={() => {
                            this.setState({
                                title_isempty: this.state.title === "",
                                description_isempty:
                                    this.state.description === "",
                            });

                            if (
                                this.state.title === "" ||
                                this.state.description === ""
                            ) {
                                return;
                            }
                            this.sendNewSuggestion();
                        }}
                    />
                </div>
            </div>
        );
    }
}*/

const mapStateToProps = (state, ownProps) => ({});

const mapDispatchToProps = dispatch => ({
    toastOpen: toastData =>
        dispatch(DigitToastActions.digitToastOpen(toastData)),
});

export const Prompt = connect(
    mapStateToProps,
    mapDispatchToProps
)(PromptView);

export default Prompt;
