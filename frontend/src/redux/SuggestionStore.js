import { createStore } from "redux";
import {
  ADD_SUGGESTION, 
  SET_SUGGESTIONS
} from "./suggestionstore.actions";

const reducer = (state, action) => {
    if (action.suggestion != null)
        switch (action.type) {
            case SET_SUGGESTIONS:
                return [...action.suggestion].sort((b, a) => 
                        new Date(a.timestamp).getTime() -
                        new Date(b.timestamp).getTime()
                    );
            case ADD_SUGGESTION:
                return state.concat(action.suggestion).sort((b, a) =>
                        new Date(a.timestamp).getTime() -
                        new Date(b.timestamp).getTime()
                );
            default:
                return state;
        }

    return state;
};

export const suggestions = createStore(reducer, []);

suggestions.subscribe(() => {
    console.log("store changed");
});
