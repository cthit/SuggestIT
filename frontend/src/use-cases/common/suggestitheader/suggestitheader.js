import React, { Component } from "react";
import {DigitHeader} from '@cthit/react-digit-components';

class SuggestITHeader extends Component {
    constructor(props){
        super(props);
        this.state = {
            renderMain: props.renderMain
        }
    }

  render() {
    return (
      <DigitHeader renderMain = {this.state.renderMain} title = "SuggestIT"/>
    );
  }
}

export default SuggestITHeader;
