import React, { Component } from 'react';
import {DigitText} from '@cthit/react-digit-components'
import axios from 'axios';
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
        <div>
            <DigitText.Heading6 text={this.state.title}/>
            <DigitText.Subtitle2 text={this.state.author}/>
            <DigitText.Subtitle2 text={this.state.timestamp}/>
            <DigitText.Text text={this.state.text}/>
        </div>
            );
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