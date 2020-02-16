import { suggestions } from "../redux/SuggestionStore";
import axios from "axios";
import Cookies from "universal-cookie";
import { SET_SUGGESTIONS } from "../redux/suggestionstore.actions";
import * as jwt from "jsonwebtoken";

const cookies = new Cookies();
const baseUrl = process.env.REACT_APP_BACKEND_URL;

const authCookieName = "AUTH_TOKEN";

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

export const checkLogin = () =>
    new Promise((resolve, reject) => {
        if (!cookies.get(authCookieName)) {
            reject(undefined);
            return;
        }

        axios
            .get(`${baseUrl}/authenticate`, {
                headers: {
                    Authorization: cookies.get(authCookieName),
                },
            })
            .then(val => resolve(val))
            .catch(err => {
                let token = jwt.decode(cookies.get(authCookieName));
                logOut();
                if (!token || !token.groups || token.groups.includes("prit")) {
                    reject(undefined);
                }
                console.log(
                    "You must be a member of P.R.I.T. to read suggestions"
                );
                reject("You must be a member of P.R.I.T. to read suggestions");
            });
    });

export const logOut = () => cookies.remove(authCookieName);
