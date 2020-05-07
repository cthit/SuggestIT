import axios from "axios";

const origin = window.location.origin;

export const updateSuggestions = setter =>
    axios
        .get(`/api`)
        .then(res => {
            setter(res.data ? res.data : []);
        })
        .catch(err => {
            //User might not be logged in
        });

export const getSuggestion = uuid => axios.get(`/api/suggestion?Id=${uuid}`);

export const addSuggestion = _suggestion => axios.post(`/api/`, _suggestion);

export const deleteSuggestion = uuid =>
    axios.delete(`/api/delete?Id=${uuid}`).catch(error => {
        //User might not be logged in
    });

export const deleteSuggestions = suggestions =>
    axios.put(`/api/delete`, {
        ids: suggestions.map(e => e.id),
    });

export const loginRedirect = () => (window.location = "/api/login");

export const checkLogin = setUser =>
    axios
        .get(`/api/checkLogin`)
        .then(res => setUser(true))
        .catch(err => setUser(null));

export const getToken = code => axios.post(`/api/auth/withCode?code=${code}`);

export const logOut = () => axios.post("/api/logout");
