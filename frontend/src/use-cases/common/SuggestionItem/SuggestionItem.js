import React, { useState } from "react";
import {
    DigitNavLink,
    DigitText,
    DigitIconButton,
    DigitToastActions,
    DigitLayout,
} from "@cthit/react-digit-components";
import "./SuggestionItem.css";
import { translateTimestamp } from "../methods";
import { Clear, ExpandMore, ExpandLess } from "@material-ui/icons";
import {
    updateSuggestions,
    deleteSuggestion,
    addSuggestion,
} from "../../../services/data.service";
import { connect } from "react-redux";

/*
  For this component to work you need to have DigitProvider around it and a DigitToast in 
  
  Minimal example:
  <DigitProviders>
    <DigitToast />
    <SuggestionItem />
  </DigitProviders>
  */

const SuggestionItemView = ({ suggestion, ts, ...props }) => {
    const [text, dispatchToggle] = useState(null);
    const toastOpen = props["toastOpen"];

    const toggleText = () => {
        dispatchToggle(
            text ? null : (
                <DigitText.Subtitle
                    className="suggestionText"
                    text={suggestion.text}
                />
            )
        );
    };

    const deleteWithToast = () => {
        toastOpen({
            text: "The suggestion has been deleted",
            duration: 5000,
            actionText: "Undo",
            actionHandler: () => {
                addSuggestion(suggestion).then(res => updateSuggestions());
            },
        });
        deleteSuggestion(suggestion.id).then(res => updateSuggestions());
    };

    return (
        <div className="card">
            <DigitLayout.Grid>
                <DigitLayout.Row>
                    <DigitNavLink
                        text={suggestion.title}
                        link={"/suggestion/" + suggestion.id}
                    />
                    <div className="float-right">
                        <DigitIconButton
                            icon={Clear}
                            onClick={() => deleteWithToast()}
                        />
                    </div>
                </DigitLayout.Row>
                <DigitLayout.Row>
                    <DigitText.Subtitle2
                        text={suggestion.author}
                        className="authorLabel"
                    />
                    <div className="float-right">
                        <DigitIconButton
                            icon={text ? ExpandLess : ExpandMore}
                            onClick={() => toggleText()}
                        />
                    </div>
                </DigitLayout.Row>
                <DigitLayout.Row>
                    <DigitText.Subtitle2
                        text={translateTimestamp(ts)}
                        className="timeStampLabel"
                    />
                </DigitLayout.Row>
                {text}
            </DigitLayout.Grid>
        </div>
    );
};

const mapStateToProps = (state, ownProps) => ({});

const mapDispatchToProps = dispatch => ({
    toastOpen: toastData =>
        dispatch(DigitToastActions.digitToastOpen(toastData)),
});

export const SuggestionItem = connect(
    mapStateToProps,
    mapDispatchToProps
)(SuggestionItemView);

export default SuggestionItem;
