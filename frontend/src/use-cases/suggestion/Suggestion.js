import React, { Component } from "react";
import { DigitText } from "@cthit/react-digit-components";
import "./Suggestion.css";
import { translateTimestamp } from "../common/methods";
import { getSuggestion } from "../../services/data.service";

class Suggestion extends Component {
    constructor(props) {
        super(props);

        this.state = {
            id: props.match.params.id,
            author: "",
            text: "",
            timestamp: "",
            title: "",
        };
        this.getSuggestionById();
    }

    render() {
        return (
            <div className="suggestionCard">
                <DigitText.Heading6 text={this.state.title} />
                <DigitText.Subtitle2
                    className="grayText"
                    text={"Inlagd av: " + this.state.author}
                />
                <DigitText.Subtitle2
                    className="grayText"
                    text={translateTimestamp(this.state.timestamp)}
                />
                <DigitText.Text text={this.state.text} />
            </div>
        );
    }

    getSuggestionById() {
        getSuggestion(this.state.id).then(res => {
            this.setState({
                author: res.data.author,
                text: res.data.text,
                timestamp: res.data.timestamp,
                title: res.data.title,
            });
        });
    }
}

export default Suggestion;
