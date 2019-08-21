import React from "react";
import ReactDOM from "react-dom";
import App from "./App.jsx";
import * as serviceWorker from "./serviceWorker";
import { DigitProviders } from "@cthit/react-digit-components"
import { Route } from "react-router";
import SuggestITHeader from './common/suggestitheader/suggestitheader';
ReactDOM.render(
<DigitProviders>
  <SuggestITHeader renderMain = {()=>
      <Route component={App}/>  
    }
    />
</DigitProviders>, document.getElementById("root"));

serviceWorker.unregister();
