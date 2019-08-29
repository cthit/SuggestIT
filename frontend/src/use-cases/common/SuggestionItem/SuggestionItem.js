import React, { useState } from "react";
import {
    DigitNavLink,
    DigitText,
    DigitIconButton,
    DigitToastActions,
} from "@cthit/react-digit-components";
import "./SuggestionItem.css";
import { translateTimestamp } from "../methods";
import { Clear, ExpandMore, ExpandLess } from "@material-ui/icons";
import { Grid } from "@material-ui/core";
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
            text: "Förslaget raderades",
            duration: 5000,
            actionText: "Ångra",
            actionHandler: () => {
                addSuggestion(suggestion).then(res => updateSuggestions());
            },
        });
        deleteSuggestion(suggestion.id).then(res => updateSuggestions());
    };

    return (
        <div className="card">
            <div className="innerCard">
                <Grid container>
                    <Grid item xs={10}>
                        <DigitNavLink
                            text={suggestion.title}
                            link={"/suggestion/" + suggestion.id}
                        />
                    </Grid>
                    <Grid item xs={1}>
                        <DigitIconButton
                            icon={Clear}
                            onClick={() => deleteWithToast()}
                        />
                    </Grid>
                </Grid>
                <section>
                    <Grid container>
                        <Grid item xs={10}>
                            <DigitText.Subtitle2
                                text={suggestion.author}
                                className="authorLabel"
                            />
                            <DigitText.Subtitle2
                                text={translateTimestamp(ts)}
                                className="timeStampLabel"
                            />
                        </Grid>
                        <Grid item xs={1}>
                            <DigitIconButton
                                icon={text ? ExpandLess : ExpandMore}
                                onClick={() => toggleText()}
                            />
                        </Grid>
                    </Grid>
                    {text}
                </section>
            </div>
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
