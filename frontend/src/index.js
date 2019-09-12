import React from "react";
import ReactDOM from "react-dom";
import App from "./App.jsx";
import * as serviceWorker from "./serviceWorker";
import {
    DigitProviders,
    DigitToast,
    DigitDialog,
} from "@cthit/react-digit-components";
import { Route } from "react-router";
import SuggestITHeader from "./common/suggestitheader/suggestitheader";
ReactDOM.render(
    <DigitProviders 
        children = {
            <div>
                <DigitToast />
                <DigitDialog />
                <SuggestITHeader renderMain={() => <Route component={App} />} />
            </div>
        }
    />,
    document.getElementById("root")
);

serviceWorker.unregister();
