import { bindActionCreators, } from "redux";
export { uploadAction } from "./upload"
import { UploadListActions } from "./upload"
export const actionCreators = {
    uploadListFileAdd: UploadListActions.uploadListFileAdd, 
    uploadListFileRemove: UploadListActions.uploadListFileRemove
};
