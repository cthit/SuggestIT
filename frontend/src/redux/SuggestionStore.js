import {createStore} from "redux";

const reducer = function(state, action){
    if(action.suggestion != null)
        return state.concat(action.suggestion).sort(function(b,a){
            return new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime();
        });
    
    return state;
}

export const suggestions = createStore(reducer,[]);

suggestions.subscribe(()=>{
    console.log("store changed")
})
