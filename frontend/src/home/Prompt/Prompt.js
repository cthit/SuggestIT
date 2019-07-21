import React, { Component } from 'react';
import {DigitTextField, DigitTextArea, DigitButton, DigitSwitch, DigitLayout, DigitText} from '@cthit/react-digit-components'
import axios from 'axios';
import {suggestions} from '../../redux/SuggestionStore';
import './Prompt.css';

class Prompt extends Component{

    constructor(props){
        super(props);
        this.state = {
            title: "",
            description: "",
            author: "",
            anonymus_author: false,
            title_isempty: false,
            title_error_message: "Titeln är ej ifylld!",
            description_isempty: false,
            description_error_message: "Du måste lägga in en förklaring"
        }
    }

    sendNewSuggestion(){
        //OBS! Det står http istället för https
        axios.post('http://localhost:5000/', {
            title: this.state.title,
            text: this.state.description,
            author: ((this.state.author === "" || this.state.anonymus_author)? "Anonym": this.state.author)
        })
        .then(res=>{
            console.log("Response from new add" + JSON.stringify(res.data));
            suggestions.dispatch({
                type: "add",
                suggestion: res.data
            });

            this.setState({
                title: "",
                description: "",
                author: ""
            });
        })
        .catch(error=>{
            console.log("RIP, som error accured. =(");
            console.log(error);
        });
    }

    render() {
        return(
            <div className="prompt">
                <div className="innerPrompt">
                        <DigitText.Heading6 text="Nytt förslag"/>
                        <DigitTextField
                        error={this.state.title_isempty}
                        errorMessage = {this.state.title_error_message}
                        onChange={e => {
                            this.setState({
                                title: e.target.value
                            });
                        }}
                        value={this.state.title}
                        upperLabel="Rubrik"
                    />
                    <DigitTextArea
                        error={this.state.description_isempty}
                        errorMessage = {this.state.description_error_message}
                        onChange={e => {
                            this.setState({
                                description: e.target.value
                            });
                        }}
                        value={this.state.description}
                        upperLabel="Förslag"
                        rows={5}
                        rowsMax={10}
                    />
                    <DigitLayout.Row>
                        <DigitTextField
                            onChange={e => {
                                this.setState({
                                    author: e.target.value
                                });
                            }}
                            value={this.state.author}
                            disabled={this.state.anonymus_author}
                            upperLabel="Författare"
                        />

                        <DigitSwitch 
                            value={this.state.anonymus_author}
                            label="Anonym"
                            primary
                            onChange = {e => {
                                this.setState({
                                    anonymus_author: e.target.checked
                                })
                            }}
                        />
                        
                        
                    </DigitLayout.Row>
                    <DigitButton text = "Skicka" onClick={a =>{
                        this.setState({
                            title_isempty: this.state.title === "",
                            description_isempty: this.state.description === ""
                        });
                        
                        if(this.state.title === "" || this.state.description === ""){
                            return;
                        }

                        this.sendNewSuggestion();

                    }}></DigitButton>
                </div>
            </div>
        );
    }
}

export default Prompt;