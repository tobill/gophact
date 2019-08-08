import React from 'react'
import {UploadFiles } from "./uploadFiles"
import {
    Field,
    reduxForm
} from 'redux-form'

export class UploadForm extends React.Component {

    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleInputChange = this.handleInputChange.bind(this);
        //this.fileInput = React.createRef();
    }
    //const { handleSubmit } = props

    handleSubmit(event) {
        event.preventDefault();
        console.log(event);
        this.props.onSubmit();
    }

    handleInputChange(event) {
        event.preventDefault();
        this.props.onInputChanged(event.target.files);
    }

    render() {
        return (
            <div>
              <form onSubmit={this.handleSubmit}>
              <div className="custom-file">
                <input type="file" className="custom-file-input" id="customFile" onInput={this.handleInputChange}  multiple />
                <label className="custom-file-label" htmlFor="customFile">Choose file</label>
              </div>
              <UploadFiles/>
              <div>
                <button className="btn btn-primary btn-block" type="submit" onSubmit={this.handleSubmit} >Submit</button> 
              </div>
              </form>
            </div>
        )
    }

}

