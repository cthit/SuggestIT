import React, { useEffect, useState, useContext } from "react";
import Cookies from "universal-cookie";
import { getToken, checkLogin } from "services/data.service";
import { DigitLoading } from "@cthit/react-digit-components";
import { Redirect } from "react-router-dom";
import UserContext from "common/context/user-context";

const cookies = new Cookies();

const Callback = ({ location }) => {
    const [redirect, setRedirect] = useState(false);
    const [, setUser] = useContext(UserContext);

    useEffect(() => {
        let params = new URLSearchParams(location.search);
        getToken(params.get("code"))
            .then(res => {
                cookies.set("AUTH_TOKEN", res.data.access_token);
                console.log(res.data);
                checkLogin(setUser, res.data.access_token);
                setRedirect(true);
            })
            .catch(err => {
                console.log(err);
                console.log("Failed to fetch token");
                setRedirect(true);
            });
    }, [location.search, setRedirect]);

    return (
        <>
            <DigitLoading size={80} />
            {redirect ? <Redirect to="/" /> : null}
        </>
    );
};

export default Callback;
