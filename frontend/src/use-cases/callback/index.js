import React, { useEffect, useState, useContext } from "react";
import { getToken, checkLogin } from "services/data.service";
import { DigitLoading, useDigitToast } from "@cthit/react-digit-components";
import { Redirect } from "react-router-dom";
import UserContext from "common/context/user-context";

const Callback = ({ location }) => {
    const [redirect, setRedirect] = useState(false);
    const [, setUser] = useContext(UserContext);
    const [toastOpen] = useDigitToast({ duration: 5000 });

    useEffect(() => {
        let params = new URLSearchParams(location.search);
        getToken(params.get("code"))
            .then(res => {
                checkLogin(setUser);
                console.log("Login success");
            })
            .catch(err => {
                console.log(err);
                toastOpen({
                    text: "You do not have the authority to view suggestions",
                });
            })
            .finally(() => setRedirect(true));
    }, [location.search, setRedirect, setUser]);

    return (
        <>
            <DigitLoading />
            {redirect ? <Redirect to="/" /> : null}
        </>
    );
};

export default Callback;
