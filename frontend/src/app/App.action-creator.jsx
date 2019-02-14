import { LOAD_ALL_SUGGESTIONS_SUCCESSFULLY } from "./App.actions.jsx";
import axios from "axios";

export function loadAllSuggestions() {
  return dispatch => {
    axios
      .get("http://localhost:5000")
      .then(response => {
        dispatch(loadAllSuggestionsSuccesfully(response.data));
      })
      .catch(error => {
        console.log("Error communicating with the server: " + error);
      });
  };
}

function loadAllSuggestionsSuccesfully(data) {
  return {
    type: LOAD_ALL_SUGGESTIONS_SUCCESSFULLY,
    payload: {
      ...data
    },
    error: false
  };
}
