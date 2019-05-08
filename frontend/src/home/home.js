import React, { Component } from "react";
import {DigitHeader} from '@cthit/react-digit-components';

import "./home.css";
import Prompt from "./Prompt/Prompt";
import SuggestionBoard from "./SuggestionBoard/SuggestionBoard";
import SuggestITHeader from '../suggestitheader/suggestitheader';

class Home extends Component {
  render() {
    return (
      <SuggestITHeader renderMain = {()=>
      <div>
        <br/>
        <Prompt/>
        <SuggestionBoard/> 
      </div>}
      />
    );
  }
}

export default Home;
