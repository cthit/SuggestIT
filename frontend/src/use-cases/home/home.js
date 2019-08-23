import React, { Component } from "react";

import "./home.css";
import Prompt from "./Prompt/Prompt";
import SuggestionBoard from "./SuggestionBoard/SuggestionBoard";
import { checkLogin } from "../../services/data.service";
import { DigitToast } from "@cthit/react-digit-components";

class Home extends Component {
  constructor(props){
    super(props)
    this.state = {
      suggestionboard: <div></div>
    }

    checkLogin().then( res=>
      this.setState({
        suggestionboard: <SuggestionBoard/>
      })
    )
  }

  render() {
    return (
      <div>
        <br/>
        <DigitToast />
        <Prompt/>
        {this.state.suggestionboard}
      </div>
    );
  }
}

export default Home;
