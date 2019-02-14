import { combineReducers } from "redux";
import { LOAD_ALL_SUGGESTIONS_SUCCESSFULLY } from "./App.actions";

export const rootReducer = combineReducers({
  suggestIT
});

export function suggestIT(
  state = {
    suggestions: []
  },
  action
) {
  switch (action.type) {
    case LOAD_ALL_SUGGESTIONS_SUCCESSFULLY:
      return {
        ...state,
        ...action.payload
      };

    default:
      return state;
  }
}
