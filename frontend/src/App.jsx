import React, { Component } from "react";
import { BrowserRouter, Switch } from "react-router-dom";
import { Route } from "react-router";

import Home from "./use-cases/home/home";
import NotFound from "./use-cases/notfound/notfound";
import Suggestion from "./use-cases/suggestion/Suggestion";
import "./App.css";

class App extends Component {
    render() {
        return (
            <BrowserRouter>
                <Switch>
                    <Route exact path="/" component={Home} />
                    <Route path="/suggestion/:id" component={Suggestion} />
                    <Route path="/" component={NotFound} />
                </Switch>
            </BrowserRouter>
        );
    }
}

export default App;
