import React, { useState, useEffect } from "react";

import Prompt from "./views/Prompt";
import SuggestionBoard from "./views/SuggestionBoard/SuggestionBoard";
import { checkLogin } from "../../services/data.service";
import { HomeWrapper, Divider } from "./home.style";

const Home = () => {
    const [suggestionboard, setSuggestionBoard] = useState(<div />);

    useEffect(() => {
        checkLogin().then(() => setSuggestionBoard(<SuggestionBoard />));
        return () => {};
    }, []);

    return (
        <HomeWrapper>
            <Prompt />
            <Divider />
            {suggestionboard}
        </HomeWrapper>
    );
};

export default Home;
