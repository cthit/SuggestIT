import React, { Component } from "react";
import "./SuggestionBoard.css";
import { suggestions } from "../../../redux/SuggestionStore";
import { SuggestionItem } from "../../common/SuggestionItem/SuggestionItem";
import {
    updateSuggestions,
    deleteSuggestions,
} from "../../../services/data.service";
import {
    DigitButton,
    DigitText,
    DigitDialogActions,
} from "@cthit/react-digit-components";
import { connect } from "react-redux";

class SuggestionBoardView extends Component {
    constructor(props) {
        super(props);
        this.state = {
            suggestions: [],
            clearButton: null,
            dialogOpen: props.dialogOpen,
        };
        suggestions.subscribe(() => {
            this.setState({
                suggestions: suggestions.getState(),
                clearButton:
                    suggestions.getState().length > 0 ? (
                        this.ClearButton()
                    ) : (
                        <DigitText.Title text="The suggestion box is empty." />
                    ),
            });
        });
        this.getData();
    }

    getData() {
        updateSuggestions();
    }

    render() {
        return (
            <div>
                <div className="grid">{this.state.clearButton}</div>
                <div className="grid">
                    {this.state.suggestions.map(obj => (
                        <SuggestionItem
                            key={obj.id}
                            suggestion={obj}
                            ts={obj.timestamp}
                        />
                    ))}
                </div>
            </div>
        );
    }

    ClearButton = () => (
        <DigitButton
            className="clear-button"
            text="Clear suggestions"
            primary
            raised
            onClick={() => {
                this.state.dialogOpen({
                    title: "Are you sure",
                    description: `Are you sure you want to delete ${this.state.suggestions.length} suggestions?`,
                    cancelButtonText: "No",
                    confirmButtonText: "Yes",
                    onConfirm: () => this.clearSuggestions(),
                });
            }}
        />
    );

    clearSuggestions = () => {
        this.setState({
            confirmOpen: false,
        });
        deleteSuggestions(this.state.suggestions).then(res =>
            updateSuggestions()
        );
    };

    openConfirm = () => {
        this.setState({
            confirmOpen: true,
        });
    };
}

const mapStateToProps = (state, ownProps) => ({});

const mapDispatchToProps = dispatch => ({
    dialogOpen: dialogData =>
        dispatch(DigitDialogActions.digitDialogOpen(dialogData)),
});

const SuggestionBoard = connect(
    mapStateToProps,
    mapDispatchToProps
)(SuggestionBoardView);

export default SuggestionBoard;
