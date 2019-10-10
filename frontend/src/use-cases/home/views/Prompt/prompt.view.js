import React, { useState } from "react";
import {
    DigitTextField,
    DigitTextArea,
    DigitButton,
    DigitSwitch,
    DigitText,
    DigitToastActions,
} from "@cthit/react-digit-components";
import {
    addSuggestion,
    updateSuggestions,
} from "../../../../services/data.service";
import { connect } from "react-redux";
import "./prompt.style.css";

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
                <DigitTextField
                    error={errors.title_error}
                    errorMessage={title_error_message}
                    onChange={e => setTitle(e.target.value)}
                    value={title}
                    upperLabel="Title"
                />
                <br />
                {/*Change this tag to DigitTextArea when the*/}
                <DigitTextArea
                    error={errors.description_error}
                    errorMessage={description_error_message}
                    onChange={e => setText(e.target.value)}
                    value={text}
                    upperLabel="Suggestion"
                    rows={5}
                    rowsMax={10}
                />
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
