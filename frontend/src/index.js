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
import SuggestITHeader from "./common/suggestitheader";
import { UserProvider } from "./common/context/user-context/index.js";

ReactDOM.render(
    <DigitProviders
        children={
            <div>
                <DigitToast />
                <DigitDialog />
                <UserProvider>
                    <SuggestITHeader
                        renderMain={() => <Route component={App} />}
                    />
                </UserProvider>
            </div>
        }
    />,
    document.getElementById("root")
);

serviceWorker.unregister();
