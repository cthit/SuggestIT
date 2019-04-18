import React, { Component } from 'react';
import {DigitTextField, DigitTextArea, DigitButton} from '@cthit/react-digit-components'

import './Prompt.css';

class Prompt extends Component{

    constructor(props){
        super(props);
        this.state = {
            title: "",
            description: ""
        }
    }

    render() {
        return(
        <div>
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
                upperLabel="Förklaring"
                rows={5}
                rowsMax={10}
            />
            <DigitButton text = "Skicka"></DigitButton>
        </div>
        );
    }
}

export default Prompt;