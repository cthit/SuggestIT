import React, { useEffect, useState, useContext } from "react";
import { getToken, checkLogin } from "services/data.service";
import { DigitLoading } from "@cthit/react-digit-components";
import { Redirect } from "react-router-dom";
import UserContext from "common/context/user-context";

const Callback = ({ location }) => {
    const [redirect, setRedirect] = useState(false);
    const [, setUser] = useContext(UserContext);

    useEffect(() => {
        let params = new URLSearchParams(location.search);
        getToken(params.get("code"))
            .then(res => {
                checkLogin(setUser);
                setRedirect(true);
            })
            .catch(err => {
                console.log(err);
                console.log("Failed to fetch token");
                setRedirect(true);
            });
    }, [location.search, setRedirect, setUser]);

    return (
        <>
            <DigitLoading />
            {redirect ? <Redirect to="/" /> : null}
        </>
    );
};

export default Callback;
