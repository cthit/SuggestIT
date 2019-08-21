import React from 'react';
import { DigitNavLink, DigitText, DigitIconButton } from '@cthit/react-digit-components';
import './SuggestionItem.css';
import { translateTimestamp } from '../methods';
import { Clear } from '@material-ui/icons';
import { Grid } from '@material-ui/core';
import { updateSuggestions, deleteSuggestion } from '../../../services/data.service';

export const SuggestionItem = ({suggestion,ts,showtext})=>{

    let text;
    if(showtext)
        text = <DigitText.Subtitle className="suggestionText" text={suggestion.text}/>;

    return(
    <div className ="card">
      <div className="innerCard">
        <Grid container >
          <Grid item xs={10}>
            <DigitNavLink text={suggestion.title} link={"/suggestion/" + suggestion.id}/>
          </Grid>
          <Grid item xs={1}>
            <DigitIconButton icon={Clear} onClick= {() => _deleteSuggestion(suggestion.id)}/>
          </Grid>
        </Grid>
        <section>
          <DigitText.Subtitle2 text={suggestion.author} className = "authorLabel"/>
          <DigitText.Subtitle2 text={translateTimestamp(ts)} className = "timeStampLabel"/>
          {text}
        </section>
      </div>
    </div>
    )};

export const _deleteSuggestion = (uuid) =>
  {
    deleteSuggestion(uuid).then(res => 
      updateSuggestions()
    )
  }

export default SuggestionItem;