import { suggestions } from "../redux/SuggestionStore";
import axios from "axios";
import Cookies from "universal-cookie";
import { SET_SUGGESTIONS } from "../redux/suggestionstore.actions";

const cookies = new Cookies();
const baseUrl =
    process.env.NODE_ENV === "development"
        ? "http://localhost:5000/api"
        : "https://suggestit.chalmers.it/api";
const authCookieName = "PRIT_AUTH_KEY";

export const updateSuggestions = () =>
    axios
        .get(`${baseUrl}/`, {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        })
        .then(res => {
            suggestions.dispatch({
                type: SET_SUGGESTIONS,
                suggestion: res.data ? res.data : [],
            });
        })
        .catch(err => {
            //User might not be logged in
        });

export const getSuggestion = uuid =>
    axios
        .get(`${baseUrl}?Id=${uuid}`, {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        })
        .catch(error => {
            //User might not be logged in
        });

export const addSuggestion = _suggestion =>
    axios.post(`${baseUrl}/`, _suggestion);

export const deleteSuggestion = uuid =>
    axios
        .delete(`${baseUrl}/delete?Id=${uuid}`, {
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

export const login = (cid, password) =>
    new Promise((resolve, reject) =>
        axios
            .put(`${baseUrl}/authenticate`, {
                cid: cid,
                password: password,
            })
            .then(res => {
                cookies.set(authCookieName, res.data.token);
                resolve(true);
            })
            .catch(err => {
                reject("CID and password did not match");
            })
    );

export const checkLogin = () =>
    axios.get(`${baseUrl}/authenticate`, {
        headers: {
            Authorization: cookies.get(authCookieName),
        },
    });

export const logOut = () => cookies.remove(authCookieName);
