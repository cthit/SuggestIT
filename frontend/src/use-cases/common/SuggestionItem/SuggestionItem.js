import React, { Component } from 'react';
import { DigitNavLink, DigitText, DigitIconButton } from '@cthit/react-digit-components';
import './SuggestionItem.css';
import { translateTimestamp } from '../methods';
import { Clear, ExpandMore, ExpandLess } from '@material-ui/icons';
import { Grid } from '@material-ui/core';
import { updateSuggestions, deleteSuggestion } from '../../../services/data.service';

export class SuggestionItem extends Component{

  constructor({suggestion, ts, ...props}){
    super(props);
    this.state = {
      showtext: false,
      text: <div></div>,
      suggestion: suggestion,
      ts: ts
    }
  }

  toggleText = () =>{

    this.setState({
      showtext: !this.state.showtext,
      text: this.state.showtext ? <div></div> : <DigitText.Subtitle className="suggestionText" text={this.state.suggestion.text}/>
    })
  }

    render = () =>
    <div className ="card">
      <div className="innerCard">
        <Grid container >
          <Grid item xs={10}>
            <DigitNavLink text={this.state.suggestion.title} link={"/suggestion/" + this.state.suggestion.id}/>
          </Grid>
          <Grid item xs={1}>
            <DigitIconButton icon={Clear} onClick= {() => _deleteSuggestion(this.state.suggestion.id)}/>
          </Grid>
        </Grid>
        <section>
          <Grid container>
            <Grid item xs={10}>
              <DigitText.Subtitle2 text={this.state.suggestion.author} className = "authorLabel"/>
              <DigitText.Subtitle2 text={translateTimestamp(this.state.ts)} className = "timeStampLabel"/>
            </Grid>
            <Grid item xs={1}>
              <DigitIconButton icon={this.state.showtext ? ExpandLess : ExpandMore} onClick={this.toggleText}/>
            </Grid>
          </Grid>
          {this.state.text}
        </section>
      </div>
    </div>
    
  };

export const _deleteSuggestion = (uuid) =>
  {
    deleteSuggestion(uuid).then(res => 
      updateSuggestions()
    )
  }

export default SuggestionItem;