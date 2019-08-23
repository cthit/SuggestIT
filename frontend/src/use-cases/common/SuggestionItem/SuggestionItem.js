import React, { Component, useReducer, useState } from 'react';
import { DigitNavLink, DigitText, DigitIconButton } from '@cthit/react-digit-components';
import './SuggestionItem.css';
import { translateTimestamp } from '../methods';
import { Clear, ExpandMore, ExpandLess } from '@material-ui/icons';
import { Grid, Snackbar } from '@material-ui/core';
import { updateSuggestions, deleteSuggestion } from '../../../services/data.service';
import { UndoSnackbar } from '../UndoSnackbar/UndoSnackbar';

export const SuggestionItem = ({suggestion, ts}) => {
  const [text, dispatchToggle] = useReducer(toggleText, null);
  const [snackOpen, setSnackVisible] = useState(false);
  const [visible, setVisible] = useState(true);

  const toggleText = (state, action) =>{
    return state ? null : <DigitText.Subtitle className="suggestionText" text={suggestion.text}/>;
  }

  const hadleSnackClose = () =>{
    setSnackVisible(false);
    deleteSuggestion(suggestion.id).then(res => 
      updateSuggestions()
    )
  }

  const undo = () => {
    console.log("Undo!")
  }

  const _deleteSuggestion = (uuid) =>
  {
    setSnackVisible(true);
    
  }

  return (
    <div>
      {visible ? <div className ="card">
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
      <Grid container>
        <Grid item xs={10}>
          <DigitText.Subtitle2 text={suggestion.author} className = "authorLabel"/>
          <DigitText.Subtitle2 text={translateTimestamp(ts)} className = "timeStampLabel"/>
        </Grid>
        <Grid item xs={1}>
          <DigitIconButton icon={text ? ExpandLess : ExpandMore} onClick={() => dispatchToggle({type: 'Toggle'})}/>
        </Grid>
      </Grid>
      {text}
      </section>
    </div>
  </div> : <UndoSnackbar open={snackOpen} handleClose={hadleSnackClose} handleUndo={undo}/>
  }
    </div>
    )
    
  };



export default SuggestionItem;