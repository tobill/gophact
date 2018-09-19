import {createStore, applyMiddleware, Middleware,} from "redux";
import {composeWithDevTools,} from "redux-devtools-extension";
import thunk from "redux-thunk";
import {reducer,} from "./root-reducer";
import {actionCreators,} from "./root-actions";


export const store = createStore(reducer, composeWithDevTools(applyMiddleware(thunk)));




