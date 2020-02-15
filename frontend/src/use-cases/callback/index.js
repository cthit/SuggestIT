import { Component } from "react";
import Cookies from "universal-cookie";
const cookies = new Cookies();

class Callback extends Component {
    constructor(props) {
        super(props);
        let params = new URLSearchParams(props.location.search);
        cookies.set("AUTH_TOKEN", params.get("token"));
        window.location.replace("/");
    }
}

export default Callback;
