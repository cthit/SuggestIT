import React, {Component} from 'react';
import './SuggestionBoard.css';
import {suggestions} from '../../../redux/SuggestionStore';
import SuggestionItem from '../../common/SuggestionItem/SuggestionItem';
import { updateSuggestions } from '../../../services/data.service';

class SuggestionBoard extends Component{
    constructor(props){
      super(props);
      this.state = {
        suggestions: []
      }
      suggestions.subscribe(()=>{
        this.setState({
          suggestions: suggestions.getState()
        })
      })
      this.getData()
    }
    
    getData(){
      updateSuggestions()
    }


    render(){
        return(
            
        <div className="grid">
          {this.state.suggestions.map(obj=>
          <SuggestionItem
            key={obj.id} 
            suggestion={obj} 
            ts={obj.timestamp}
            />
          )}
        </div>
        )
    }
}



export default SuggestionBoard;