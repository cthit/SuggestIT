import React, { Component } from "react";
import { DigitText } from "@cthit/react-digit-components";
import "./suggestion.style.css";
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
            <div className="suggestionCard main">
                <DigitText.Heading6 text={this.state.title} />
                <DigitText.Subtitle2
                    className="grayText"
                    text={"Posted by: " + this.state.author}
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
            if (!res) {
                //Redirects to / if no suggestion could be found
                window.location.href = "/";
                return;
            }
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
