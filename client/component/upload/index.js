import React from 'react'
import {
    Field,
    reduxForm
} from 'redux-form'

class UploadForm extends React.Component {

    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleInputChange = this.handleInputChange.bind(this);
        this.fileInput = React.createRef();
    }
    //const { handleSubmit } = props

    handleSubmit(event) {
        event.preventDefault();
        console.log("call onsubmit");
        this.props.onSubmit(this.fileInput.current.files);
    }


    handleInputChange(event) {
        event.preventDefault();
        console.log("call oninput");
    }

    render() {
        return (
            <div>
              <form onSubmit={this.handleSubmit}>
              <div className="custom-file">
                <input type="file" className="custom-file-input" id="customFile" onInput={this.handleInputChange} ref={this.fileInput} multiple />
                <label className="custom-file-label" htmlFor="customFile">Choose file</label>
              </div>
              <div>
                <button className="btn btn-primary" type="submit" onSubmit={this.handleSubmit} >Submit</button> 
              </div>
              </form>
            </div>
        )
    }

}

export default UploadForm