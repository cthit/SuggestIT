import { suggestions } from "../redux/SuggestionStore";
import axios from "axios";
import Cookies from "universal-cookie";

const cookies = new Cookies();
const baseUrl = "http://suggestit.chalmers.it/api";
const authCookieName = "PRIT_AUTH_KEY";

export const updateSuggestions = () =>
    axios
        .get(baseUrl, {
            headers: {
                Authorization: cookies.get("PRIT_AUTH_KEY"),
            },
        })
        .then(res => {
            suggestions.dispatch({
                type: "set",
                suggestion: res.data,
            });
        })
        .catch(err => {
            //User might not be logged in
        });

export const getSuggestion = uuid =>
    axios
        .get(`${baseUrl}/${uuid}`, {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        })
        .catch(error => {
            //User might not be logged in
        });

export const addSuggestion = _suggestion => axios.post(baseUrl, _suggestion);

export const deleteSuggestion = uuid =>
    axios
        .delete(`${baseUrl}/${uuid}`, {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        })
        .catch(error => {
            //User might not be logged in
        });

export const deleteSuggestions = suggestions =>
    axios.put(
        `${baseUrl}/delete`,
        {
            ids: suggestions.map(e => e.id),
        },
        {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        }
    );

export const login = password =>
    axios
        .put(`${baseUrl}/authenticate`, { password: password })
        .then(res => cookies.set(authCookieName, res.data.key));

export const checkLogin = () =>
    axios.get(`${baseUrl}/authenticate`, {
        headers: {
            Authorization: cookies.get(authCookieName),
        },
    });
