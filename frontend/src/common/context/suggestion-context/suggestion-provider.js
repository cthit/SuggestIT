import React, { useState } from "react";
import SuggestionsContext from "./suggestion-context";

const SuggestionsProvider = ({ children }) => {
    const [suggestions, setSuggestions] = useState([]);
    return (
        <SuggestionsContext.Provider value={[suggestions, setSuggestions]}>
            {children}
        </SuggestionsContext.Provider>
    );
};

export default SuggestionsProvider;
