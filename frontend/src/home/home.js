import React, { Component } from "react";
import {DigitHeader} from '@cthit/react-digit-components';

import "./home.css";
import Prompt from "./Prompt";

class Home extends Component {
  render() {
    return (
      <DigitHeader renderMain = {()=>
      <div>
          <Prompt/>
      </div>}>
      </DigitHeader>
    );
  }
}

export default Home;
