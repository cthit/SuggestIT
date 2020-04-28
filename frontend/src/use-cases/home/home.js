import React, { useState, useEffect } from "react";

import Prompt from "./views/Prompt";
import SuggestionBoard from "./views/SuggestionBoard/SuggestionBoard";
import { checkLogin } from "../../services/data.service";
import { DigitDesign, DigitLayout } from "@cthit/react-digit-components";

const Home = () => {
    const [suggestionboard, setSuggestionBoard] = useState(<div />);

    useEffect(() => {
        checkLogin().then(() => setSuggestionBoard(<SuggestionBoard />));
        return () => {};
    }, []);

    return (
        <div
            style={{
                width: "100%",
                "flex-wrap": "wrap",
                display: "flex",
                justifyContent: "center",
            }}
        >
            <Prompt />
            <div style={{ "flex-basis": "100%", height: 0 }}></div>
            {suggestionboard}
        </div>
    );
};

export default Home;
