import React, { useContext } from "react";

import Prompt from "./views/Prompt";
import SuggestionBoard from "./views/SuggestionBoard/SuggestionBoard";
import UserContext from "common/context/user-context";
import { HomeWrapper, Divider } from "./home.style";

const Home = () => {
    const [user] = useContext(UserContext);

    return (
        <HomeWrapper>
            <Prompt />
            <Divider />
            {user && <SuggestionBoard />}
        </HomeWrapper>
    );
};

export default Home;
