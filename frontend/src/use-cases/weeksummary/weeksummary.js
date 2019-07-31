import React, {Component} from 'react';
import axios from 'axios';
import SuggestionItem from  '../../components/SuggestionItem/SuggestionItem';

class WeekSummary extends Component{
    constructor(props){
        super(props);
        this.state = {
            suggestions: []
        }
        
        this.getData()
    }
    
    getData(){
        
         axios.get("http://localhost:5000/")
        .then(res=>{
            this.setState({
               suggestions: res.data 
            });
        })
        .catch(error=>{
            console.log("RIP, som error accured. =(");
            console.log(error);
        });
    }

    


    render(){
        return(
            
        <div className="grid">
                {this.state.suggestions.filter(e => this.lessThanAWeek(e.timestamp)).map(obj=>
                <SuggestionItem
                    key={obj.id} 
                    suggestion={obj} 
                    ts={obj.timestamp}
                    showtext
                />)}
        </div>
        )
    }

    lessThanAWeek(ts){
        let now = new Date();
        let diff = now.getTime() - new Date(ts).getTime() + now.getTimezoneOffset()*60000;
        return diff <= 604800000;
    }
}



export default WeekSummary;