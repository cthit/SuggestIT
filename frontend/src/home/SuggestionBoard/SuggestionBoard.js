import React, {Component} from 'react';
import {DigitNavLink, DigitText} from '@cthit/react-digit-components';
import axios from 'axios';
import './SuggestionBoard.css';
import {suggestions} from '../SuggestionStore';

const SuggestionItem = ({props,ts})=>(
    <div className ="card">
        <div className="innerCard">
            <DigitNavLink text={props.title} link={"/suggestion/" + props.id}/>
            <DigitText.Subtitle2 text={props.author}/>
            <DigitText.Subtitle2 text={ts}/>
        </div>
    </div>
);

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
        
         axios.get("http://localhost:5000/")
        .then(res=>{
            /*this.setState({
               suggestions: res.data 
            });*/
            suggestions.dispatch({
                type: "add",
                suggestion: res.data
            });
        })
        .catch(error=>{
            console.log("RIP, som error accured. =(");
            console.log(error);
        });
    }

    translateTimestamp(ts){
        var pre = "Inlagd "
        var now = new Date();
        var diff = now.getTime() - new Date(ts).getTime() + now.getTimezoneOffset()*60000;
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


    render(){
        return(
            
        <div className="grid">
                {this.state.suggestions.map(obj=>
                <SuggestionItem
                    key={obj.id} 
                    props={obj} 
                    ts={this.translateTimestamp(obj.timestamp)}
                />)}
        </div>
        )
    }
}



export default SuggestionBoard;