import React, { Component } from "react";
import SuggestITHeader from '../common/suggestitheader/suggestitheader';

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
