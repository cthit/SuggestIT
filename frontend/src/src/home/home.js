import React, { Component } from "react";
import axios from 'axios';

import "./home.css";
import Prompt from "./Prompt/Prompt";
import SuggestionBoard from "./SuggestionBoard/SuggestionBoard";
import SuggestITHeader from '../suggestitheader/suggestitheader';
import {suggestions} from './SuggestionStore';



class Home extends Component {

  constructor(props){
    super(props);
    this.state = {
        suggestions: []
    }
    this.getData();
}

getData(){
  axios.get("http://localhost:5000/")
 .then(res=>{
     suggestions.dispatch({type: "ADD", suggestion: res.data});
 })
 .catch(error=>{
     console.log("RIP, som error accured. =(");
     console.log(error);
 });
}

  render() {
    return (
      <SuggestITHeader renderMain = {()=>
      <div>
        <br/>
        <Prompt/>
        <SuggestionBoard suggestions = {this.state.suggestions}/>
      </div>}
      />
    );
  }

  

}

export default Home;
