import React, { useState, useEffect } from 'react';
import { DigitHeader, DigitButton, DigitLayout, useDigitCustomDialog } from '@cthit/react-digit-components';
import About from './elements/about';
import { checkLogin, logOut, loginRedirect } from '../../services/data.service';

const SuggestITHeader = ({ renderMain }) => {
    const [ isLoggedIn, setIsloggedIn ] = useState(false);
    const [ dialogOpen, , ] = useDigitCustomDialog();

    useEffect(() => {
        checkLogin().then((res) => setIsloggedIn(true));
        return () => {};
    }, []);

    return (
        <DigitHeader
            renderMain={renderMain}
            title="SuggestIT"
            renderHeader={() => (
                <div>
                    <DigitLayout.Row>
                        <DigitButton
                            text="About"
                            onClick={(values) =>
                                dialogOpen({
                                    title: 'About suggestIT',
                                    renderMain: () => <About />,
                                    renderButtons: (confirm, cancel) => <DigitButton text="Close" onClick={cancel} />,
                                    onConfirm: (confirm, reject) => confirm()
                                })}
                        />
                        {!isLoggedIn ? (
                            <div>
                                <DigitButton text="Login" outlined onClick={() => loginRedirect()} />
                            </div>
                        ) : (
                            <DigitButton
                                outlined
                                text="logout"
                                onClick={() => {
                                    logOut();
                                    window.location.reload(false);
                                }}
                            />
                        )}
                    </DigitLayout.Row>
                </div>
            )}
        />
    );
};

export default SuggestITHeader;
