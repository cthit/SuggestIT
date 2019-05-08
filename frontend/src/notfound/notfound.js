import React, { Component } from "react";
import SuggestITHeader from '../suggestitheader/suggestitheader';

class NotFound extends Component {
  render() {
    return (
      <SuggestITHeader renderMain = {()=>
      <p>Site not found</p>
      }/>
    );
  }
}

export default NotFound;
