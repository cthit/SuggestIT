import React, { Component } from 'react';
import {DigitTextField, DigitTextArea, DigitButton, DigitSwitch, DigitLayout, DigitText} from '@cthit/react-digit-components'
import axios from 'axios';
import './Prompt.css';

class Prompt extends Component{

    constructor(props){
        super(props);
        this.state = {
            title: "",
            description: "",
            author: "",
            anonymus_author: false
        }
    }

    render() {
        return(
            <div className="prompt">
                <div className="innerPrompt">
                        <DigitText.Heading6 text="Nytt förslag"/>
                        <DigitTextField
                        onChange={e => {
                            this.setState({
                                title: e.target.value
                            });
                        }}
                        value={this.state.title}
                        upperLabel="Rubrik"
                    />
                    <DigitTextArea
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
                        if(this.state.anonymus_author){
                            this.setState({
                                author: "Anonym"
                            })
                        }

                        //OBS! Det står http istället för https
                        axios.post('http://localhost:5000/', {
                            title: this.state.title,
                            text: this.state.description,
                            author: this.state.author
                        })
                        .then(res=>{
                            this.location.reload();
                        })
                        .catch(error=>{
                            console.log("RIP, som error accured. =(");
                            console.log(error);
                        });

                    }}></DigitButton>
                </div>
            </div>
        );
    }
}

export default Prompt;