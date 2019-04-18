import React, { Component } from "react";
import {BrowserRouter, Switch} from 'react-router-dom';
import {Route} from 'react-router';
import { DigitHeader } from "@cthit/react-digit-components";

import Home from './home/home';
import NotFound from './notfound/notfound';
import Test from "./use-cases/test";
import "./App.css";

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <Switch>
          <Route exact path="/" component = {Home} />
          <Route path="/" component = {NotFound} />
        </Switch>
      </BrowserRouter>
    );
  }
}

export default App;
