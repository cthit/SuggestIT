import React, { useContext } from "react";

import Prompt from "./views/Prompt/prompt.view";
import SuggestionBoard from "./views/SuggestionBoard/SuggestionBoard";
import UserContext from "common/context/user-context";

const Home = () => {
    const [user] = useContext(UserContext);

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
            {user ? <SuggestionBoard /> : null}
        </div>
    );
};

export default Home;
