import axios from "axios";
import Cookies from "universal-cookie";

const cookies = new Cookies();
const authCookieName = "AUTH_TOKEN";
const origin = window.location.origin;
const gammaURL = "http://localhost:8081/api";

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
                `${gammaURL}/oauth/authorize?response_type=code&client_id=${res.data.client_id}&redirect_uri=${origin}/callback`
            )
        )
        .catch(err => console.log("Failed to get client_id"));

export const checkLogin = (setUser, token = cookies.get(authCookieName)) =>
    axios
        .get(`${gammaURL}/users/me`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        })
        .then(res => setUser(res.data))
        .catch(err => setUser(null));

export const getToken = code =>
    axios.post(`/api/auth/token?code=${code}&redirect_uri=${origin}/callback`);

export const logOut = () => cookies.remove(authCookieName);
