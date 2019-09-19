import React, { Component } from "react";

import "./home.css";
import Prompt from "./Prompt/Prompt";
import SuggestionBoard from "./SuggestionBoard/SuggestionBoard";
import { checkLogin } from "../../services/data.service";

class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            suggestionboard: <div></div>,
            description: "",
        };

        checkLogin().then(res =>
            this.setState({
                suggestionboard: <SuggestionBoard />,
            })
        );
        console.log(process.env);
    }

    render() {
        return (
            <div className="main">
                <Prompt />
                {this.state.suggestionboard}
            </div>
        );
    }
}

export default Home;
