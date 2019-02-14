import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import { DigitProviders } from "@cthit/react-digit-components"
import { Route } from "react-router";

ReactDOM.render(
<DigitProviders>
    <Route component={App}/>
</DigitProviders>, document.getElementById("root"));

serviceWorker.unregister();
