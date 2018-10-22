import * as React from "react";
import {bindActionCreators,} from "redux";
import {connect,} from "react-redux";
import { UploadForm }from "./upload";
import { uploadAction  ,actionCreators } from "../store/root-actions"

class Application extends React.Component {

    constructor(props) {
        super(props);
        this.onInputChanged = props.uploadListFileAdd;
        this.onSubmit = props.uploadListStart;
    }

    render() {
    
        return (
            <div>
                <div className="application">
                    <UploadForm onSubmit={this.onSubmit} onInputChanged={this.onInputChanged}/>
                </div>
            </div>
        );
    }
}

const mapStateToProps = (state) => {
    return {
        upload: state.upload
    }
};

const mapDispatchToProps = (dispatch, state) => {
    return {
        uploadListFileAdd: (fileList) => {
            dispatch(actionCreators.uploadListFileAdd(fileList));
        }, 
        uploadListStart: () => {
            console.log("line")
            dispatch(uploadAction(state.upload));
        } 
 
   } 
} 

export const ApplicationContainer = connect(mapStateToProps, mapDispatchToProps)(Application);