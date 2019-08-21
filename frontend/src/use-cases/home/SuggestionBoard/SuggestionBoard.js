import React, {Component} from 'react';
import './SuggestionBoard.css';
import {suggestions} from '../../../redux/SuggestionStore';
import SuggestionItem from '../../common/SuggestionItem/SuggestionItem';
import { updateSuggestions, deleteSuggestions } from '../../../services/data.service';
import { DigitButton, DigitText } from '@cthit/react-digit-components';
import { Dialog, DialogContent, DialogActions } from '@material-ui/core';

class SuggestionBoard extends Component{
    constructor(props){
      super(props);
      this.state = {
        suggestions: [],
        confirmOpen: false,
        clearButton: <div></div>
      }
      suggestions.subscribe(()=>{
        this.setState({
          suggestions: suggestions.getState(),
          clearButton: suggestions.getState().length > 0 ? this.ClearButton() : <div></div>
        })
      })
      this.getData()
    }
    
    getData(){
      updateSuggestions()
    }

    render(){
        return(
        <div>
          <div className="grid">
            {this.state.clearButton}
          </div>
          <Dialog open={this.state.confirmOpen} aria-labelledby="form-dialog-title">
            <DigitText.Title className="dialog-title" text="Är du säker?"/>
            <DialogContent>
              <DigitText.Subtitle text=
                {`Är du säker på att du vill ta bort ${this.state.suggestions.length} förslag?`}
                />
            </DialogContent>
            <DialogActions className="button-collection">
              <DigitButton onClick = {() =>this.setState({confirmOpen: false})}   text="Nej" primary/>
              <DigitButton onClick = {this.clearSuggestions} text="Ja" primary raised/>
            </DialogActions>
          </Dialog>
          <div className="grid">
            {this.state.suggestions.map(obj=>
            <SuggestionItem
              key={obj.id} 
              suggestion={obj} 
              ts={obj.timestamp}
              />
            )}
          </div>
        </div>
        )
    }

    ClearButton = () =>
      <DigitButton 
        className="clear-button" 
        text="Ränsa" 
        primary 
        raised
        onClick={() => this.openConfirm()}
        />

    clearSuggestions = () => {
      this.setState({
        confirmOpen: false
      });
      deleteSuggestions(this.state.suggestions)
      .then(res=> 
        updateSuggestions()  
      );
    }

    openConfirm = () => {
      console.log('Open confirm');
      this.setState({
      confirmOpen: true
    })
    }
}



export default SuggestionBoard;