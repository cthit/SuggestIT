import React, {Component} from 'react';
import {DigitNavLink, DigitText} from '@cthit/react-digit-components';
import './SuggestionBoard.css';
import {suggestions} from '../SuggestionStore';

const SuggestionItem = ({props,ts})=>(
    <div className ="card">
        <div className="innerCard">
            <DigitNavLink text={props.title} link={"/suggestion/" + props.id}/>
            <DigitText.Subtitle2 className="grayText" text={props.author}/>
            <DigitText.Subtitle2 className="grayText" text={ts}/>
        </div>
    </div>
);

class SuggestionBoard extends Component{

    constructor(props){
        super(props);
        this.state = {
            suggestions: props.suggestions
        };

        suggestions.subscribe(()=>{
            this.setState({
                suggestions: suggestions.getState()
            })
        });
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