import React, { useState } from "react";
import {
    DigitButton,
    DigitSwitch,
    DigitText,
    DigitToastActions,
} from "@cthit/react-digit-components";
import TextField from "@material-ui/core/TextField";
import {
    addSuggestion,
    updateSuggestions,
} from "../../../services/data.service";
import { connect } from "react-redux";
import "./Prompt.css";

const title_error_message = "The title is not filled in";
const description_error_message = "The description is not filled in";

const PromptView = ({ toastOpen }) => {
    const [text, setText] = useState("");
    const [title, setTitle] = useState("");
    const [author, setAuthor] = useState("");

    const [errors, setErrors] = useState({
        title_error: false,
        description_error: false,
    });
    const [anonymous_author, setAnonymousAuthor] = useState(false);

    const sendNewSuggestion = suggestion => {
        addSuggestion(suggestion).then(res => {
            updateSuggestions();
            setTitle("");
            setAuthor("");
            setText("");
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
                <TextField
                    error={errors.title_error}
                    helperText={errors.title_error && title_error_message}
                    onChange={e => setTitle(e.target.value)}
                    value={title}
                    label="Title"
                    inputProps={{
                        maxlength: 50
                    }}
                />
                <br />
                {/*Change this tag to DigitTextArea when the*/}
                <TextField
                    multiline
                    rows={5}
                    rowsMax={10}
                    error={errors.description_error}
                    helperText={errors.description_error && description_error_message}
                    onChange={e => setText(e.target.value)}
                    value={text}
                    label="Suggestion"
                    inputProps={{
                        maxlength: 500
                    }}
                />
                <TextField
                    onChange={e => setAuthor(e.target.value)}
                    value={author}
                    label="CID"
                    disabled={anonymous_author}
                    inputProps={{
                        maxlength: 50
                    }}
                />
                <DigitSwitch
                    value={anonymous_author}
                    label="Anonymous"
                    primary
                    onChange={() => {
                        setAnonymousAuthor(!anonymous_author);
                    }}
                />
                <DigitButton
                    text="Send"
                    primary
                    raised
                    onClick={() => {
                        setErrors({
                            title_error: title === "",
                            description_error: text === "",
                        });

                        if (title === "" || text === "") return;

                        sendNewSuggestion({
                            title: title,
                            text: text,
                            author:
                                author === "" || anonymous_author
                                    ? "Anonymous"
                                    : author,
                        });
                    }}
                />
            </div>
        </div>
    );
};

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
