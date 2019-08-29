import React, { Component } from "react";
import "./SuggestionBoard.css";
import { suggestions } from "../../../redux/SuggestionStore";
import { SuggestionItem } from "../../common/SuggestionItem/SuggestionItem";
import {
    updateSuggestions,
    deleteSuggestions,
} from "../../../services/data.service";
import { DigitButton, DigitText } from "@cthit/react-digit-components";
import ConfirmModal from "../../common/ConfirmModal/ConfirmModal";

class SuggestionBoard extends Component {
    constructor(props) {
        super(props);
        this.state = {
            suggestions: [],
            confirmOpen: false,
            clearButton: <div></div>,
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
                <ConfirmModal
                    open={this.state.confirmOpen}
                    onConfirm={this.clearSuggestions}
                    onClose={() => this.setState({ confirmOpen: false })}
                    text={`Are you sure you want to delete ${this.state.suggestions.length} suggestions?`}
                />
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
            onClick={() => this.openConfirm()}
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

export default SuggestionBoard;
