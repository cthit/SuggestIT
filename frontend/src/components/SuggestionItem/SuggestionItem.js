import React from 'react';
import {DigitNavLink, DigitText} from '@cthit/react-digit-components';
import './SuggestionItem.css';

export const SuggestionItem = ({suggestion,ts,showtext})=>{

    let text;
    if(showtext)
        text = <DigitText.Subtitle className="suggestionText" text={suggestion.text}/>;

    return(
    <div className ="card">
        <div className="innerCard">
            <DigitNavLink text={suggestion.title} link={"/suggestion/" + suggestion.id}/>
            <section>
                <DigitText.Subtitle2 text={suggestion.author} className = "authorLabel"/>
                <DigitText.Subtitle2 text={ts} className = "timeStampLabel"/>
                {text}
            </section>
        </div>
    </div>
    )};

export default SuggestionItem;