import React, { Component } from "react";
import { BrowserRouter, Switch } from "react-router-dom";
import { Route } from "react-router";

import Home from "./use-cases/home/home";
import NotFound from "./use-cases/notfound";
import Suggestion from "./use-cases/suggestion";
import Callback from "./use-cases/callback";
import "./App.css";

class App extends Component {
    render() {
        return (
            <BrowserRouter>
                <Switch>
                    <Route exact path="/" component={Home} />
                    <Route path="/suggestion/:id" component={Suggestion} />
                    <Route path="/callback" component={Callback} />
                    <Route path="/" component={NotFound} />
                </Switch>
            </BrowserRouter>
        );
    }
}

export default App;
