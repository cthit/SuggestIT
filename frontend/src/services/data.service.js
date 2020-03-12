import axios from "axios";
import Cookies from "universal-cookie";
import * as jwt from "jsonwebtoken";

const cookies = new Cookies();
const authCookieName = "AUTH_TOKEN";

export const updateSuggestions = setter =>
    axios
        .get(`/api`, {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        })
        .then(res => {
            setter(res.data ? res.data : []);
        })
        .catch(err => {
            //User might not be logged in
        });

export const getSuggestion = uuid =>
    axios.get(`/api/suggestion?Id=${uuid}`, {
        headers: {
            Authorization: cookies.get(authCookieName),
        },
    });

export const addSuggestion = _suggestion => axios.post(`/api/`, _suggestion);

export const deleteSuggestion = uuid =>
    axios
        .delete(`/api/delete?Id=${uuid}`, {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        })
        .catch(error => {
            //User might not be logged in
        });

export const deleteSuggestions = suggestions =>
    axios.put(
        `/api/delete`,
        {
            ids: suggestions.map(e => e.id),
        },
        {
            headers: {
                Authorization: cookies.get(authCookieName),
            },
        }
    );

export const loginRedirect = () =>
    axios
        .get(`/api/clientid`)
        .then(res =>
            window.location.replace(
                `https://ldap-auth.chalmers.it/authenticate?client_id=${res.data.client_id}`
            )
        )
        .catch(err => console.log("Failed to get client_id"));

export const checkLogin = () =>
    new Promise((resolve, reject) => {
        if (!cookies.get(authCookieName)) {
            reject(undefined);
            return;
        }

        axios
            .get(`/api/authenticate`, {
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
