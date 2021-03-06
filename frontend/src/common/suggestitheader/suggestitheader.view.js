import React, { useContext, useEffect } from "react";
import {
    DigitHeader,
    DigitButton,
    DigitLayout,
    useDigitCustomDialog,
} from "@cthit/react-digit-components";
import About from "./elements/about";
import { logOut, loginRedirect, checkLogin } from "services/data.service";
import UserContext from "../context/user-context";

const SuggestITHeader = ({ renderMain }) => {
    const [user, setUser] = useContext(UserContext);
    const [dialogOpen, ,] = useDigitCustomDialog();

    useEffect(() => {
        checkLogin(setUser);
    }, [setUser]);

    return (
        <DigitHeader
            renderMain={renderMain}
            title="SuggestIT"
            headerRowProps={{ flex: "1", justifyContent: "space-between" }}
            renderHeader={() => (
                <>
                    <DigitLayout.Row>
                        <DigitButton
                            text="About"
                            onClick={values =>
                                dialogOpen({
                                    title: "About suggestIT",
                                    renderMain: () => <About />,
                                    renderButtons: (confirm, cancel) => (
                                        <DigitButton
                                            text="Close"
                                            onClick={cancel}
                                        />
                                    ),
                                    onConfirm: (confirm, reject) => confirm(),
                                })
                            }
                        />
                        {!user ? (
                            <>
                                <DigitButton
                                    text="Login"
                                    outlined
                                    onClick={() => loginRedirect()}
                                />
                            </>
                        ) : (
                            <DigitButton
                                outlined
                                text="logout"
                                onClick={() => {
                                    logOut(setUser);
                                }}
                            />
                        )}
                    </DigitLayout.Row>
                </>
            )}
        />
    );
};

export default SuggestITHeader;
