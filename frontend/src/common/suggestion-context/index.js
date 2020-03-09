import React, { useState, createContext } from 'react';

export const SuggestionsContext = createContext(null);

const SuggestionsProvider = ({ children }) => {
    const [ suggestions, setSuggestions ] = useState([]);
    return (
        <SuggestionsContext.Provider value={[ suggestions, setSuggestions ]}>{children}</SuggestionsContext.Provider>
    );
};

export default SuggestionsProvider;
