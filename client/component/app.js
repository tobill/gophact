import * as React from "react";
import {bindActionCreators,} from "redux";
import {connect,} from "react-redux";
import UploadForm from "./upload";
import {uploadAction} from "../store/upload"

class Application extends React.Component {


    constructor(props) {
        super(props);
        this.submitUpload = props.uploadAction
    }

    render() {
    
        return (
            <div>
                <div className="application">
                    <UploadForm onSubmit={this.submitUpload}/>
                </div>
            </div>
        );
    }
}

const mapStateToProps = (state) => ({});

const mapDispatchToProps = (dispatch) => {
   return {
        uploadAction: (fileList) => {
            dispatch(uploadAction(fileList));
        } 
   } 
} 

export const ApplicationContainer = connect(mapStateToProps, mapDispatchToProps)(Application);