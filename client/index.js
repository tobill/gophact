import * as React from "react";
import { render, } from "react-dom";
import { Provider, } from "react-redux";
import { store, } from "./store";
import { ApplicationContainer } from "./component/app"


render( 
    <div>
        <Provider store = { store } >
            <ApplicationContainer />
        </Provider>
    </div>
    ,
    document.getElementById("root")
);
