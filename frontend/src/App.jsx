import React, { useEffect } from "react";
import { BrowserRouter, Switch, Redirect } from "react-router-dom";
import { Route } from "react-router";
import { SuggestionsProvider } from "./common/context/suggestion-context";

import Home from "./use-cases/home/home";
import NotFound from "./use-cases/notfound";
import Suggestion from "./use-cases/suggestion";
import Callback from "./use-cases/callback";
import "./App.css";

const API = () => {
    useEffect(() => {
        window.location = "http://localhost:5000/api/login";
    });

    return <></>;
};

const App = () => (
    <SuggestionsProvider>
        <BrowserRouter>
            <Switch>
                <Route exact path="/" component={Home} />
                <Route path="/suggestion/:id" component={Suggestion} />
                <Route path="/callback" component={Callback} />
                <Route path="/api" component={API} />
                <Route path="/" component={NotFound} />
            </Switch>
        </BrowserRouter>
    </SuggestionsProvider>
);

export default App;
