import { connect } from "react-redux";
import { withRouter } from "react-router-dom";
import App from "./App";
import { loadAllSuggestions } from "./App.action-creator";

const mapStateToProps = state => ({
    // Put data you want to use here!
});

const mapDispatchToProps = dispatch => ({
    // Put methods you want to use here!
    loadAllSuggestions: () => dispatch(loadAllSuggestions()),
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(withRouter(App));
