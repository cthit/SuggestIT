import React, { Component } from "react";
import {DigitHeader} from '@cthit/react-digit-components';

import "./home.css";
import Prompt from "./Prompt/Prompt";
import SuggestionBoard from "./SuggestionBoard/SuggestionBoard";

class Home extends Component {
  render() {
    return (
      <DigitHeader renderMain = {()=>
      <div>
        <br/>
        <Prompt/>
        <SuggestionBoard/> 
      </div>}
      title ="SuggestIT"
      />
    );
  }
}

export default Home;
