import * as React from "react";
import {connect,} from "react-redux";
import { actionCreators } from "../../store/root-actions"

class UploadFilesComponent extends React.Component {

    constructor(props) {
        super(props);
        this.handleRemove = this.handleRemove.bind(this);
    }

    handleRemove(event){
        this.props.uploadListFileRemove(event.target.id);
    }

    render() {
        const upload  = this.props.upload.fileList;
        console.log(upload);
        return  (
        <div className="upload-files">
            <ul className="list-group">
            {upload && upload.length > 0 && upload.map((uf, index) => {

                return <li className="list-group-item" key={index} >
                    <span>{uf.name}</span>
                    <a href="#" className="uplaod-file-item badge badge-danger float-right" 
                        onTouchEndCapture={this.handleRemove} onClick={this.handleRemove} id={index}>
                        Remove
                    </a>
                </li>
            }
            )}
            </ul>
        </div>)
    }

}

const mapStateToProps = (state) => {
    return {
        upload: state.upload
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        uploadListFileRemove: (id) => {
            dispatch(actionCreators.uploadListFileRemove(id));
        }, 
        } 
} 

export const UploadFiles = connect(mapStateToProps, mapDispatchToProps)(UploadFilesComponent);
