import React, { Component } from "react";
import Test from "../use-cases/test/";
import "./App.css";
import { DigitHeader } from "@cthit/react-digit-components";

class App extends Component {
  render() {
    return (
      <div>
        <DigitHeader dense title="SuggestIT" renderMain={() => <Test />} />
      </div>
    );
  }
}

export default App;
