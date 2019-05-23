import React, { Component } from 'react';
import {DigitText} from '@cthit/react-digit-components'
import axios from 'axios';
import SuggestITHeader from '../suggestitheader/suggestitheader';
import './Suggestion.css';

class Suggestion extends Component{

    constructor(props){
        super(props);
        
        this.state = {
            id: props.match.params.id,
            author: "",
            text: "",
            timestamp: "",
            title: ""
                
        };
        this.getSuggestionById();
   }

    render() {
        return(
        <SuggestITHeader renderMain = {()=>
            <div className="suggestionCard">
                <DigitText.Heading6 text={this.state.title}/>
                <DigitText.Subtitle2 className="grayText" text={"Inlagd av: " + this.state.author}/>
                <DigitText.Subtitle2 className="grayText" text={this.translateTimestamp(this.state.timestamp)}/>
                <DigitText.Text text={this.state.text}/>
            </div>}
        />
            );
   }

   translateTimestamp(ts){
    var pre = "Inlagd "
    var now = new Date();
    var diff = now.getTime() - new Date(ts).getTime() + now.getTimezoneOffset()*60*1000;
    diff /= 1000;
    if(diff <= 60){
        return pre + parseInt(diff, 10) + " sekunder sedan";
    }
    diff /=60;
    if(diff <= 60){
        return pre + parseInt(diff, 10) + " minuter sedan";
    }
    diff /=60;
    if(diff <= 24){
        return pre + parseInt(diff, 10) + " timmar sedan";
    }
    diff /=24;
    if(diff <= 365.25/12){
        return pre + parseInt(diff, 10) + " dagar sedan";
    }
    diff /= 365.25/12;
    if(diff <= 12){
        return pre + parseInt(diff, 10) + " månader sedan";
    }
    diff /= 12;
    return pre + parseInt(diff, 10) + " år sedan";
}

   getSuggestionById(){

        axios.get("http://localhost:5000/" + this.state.id)
        .then(res=>{
            this.setState({
                author: res.data.author,
                text: res.data.text,
                timestamp: res.data.timestamp,
                title: res.data.title
            });
        })
        .catch(error=>{
            console.log("RIP, som error accured. =(");
            console.log(error);
        });
   }
}



export default Suggestion;