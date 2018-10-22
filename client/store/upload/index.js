export const UPLOAD_LIST_FILE_ADDED = "UPLOAD_LIST_FILE_ADDED"
export const UPLOAD_LIST_UPLOAD_STARTED = "UPLOAD_LIST_UPLOAD_STARTED"
export const UPLOAD_LIST_UPLOAD_SUCCESSFUL = "UPLOAD_LIST_UPLOAD_SUCCESSFUL"
export const UPLOAD_LIST_UPLOAD_FAILURE = "UPLOAD_LIST_UPLOAD_FAILURE"

export function uploadReduce(state = {fileList: []}, action) {
    switch (action.type) {
    case UPLOAD_LIST_FILE_ADDED:
        return {
            ...state,
            fileList:  [...state.fileList, ...action.payload]
        }
    case UPLOAD_LIST_UPLOAD_SUCCESSFUL:
        return {
            ...state,
            fileList:  []
        }
    case UPLOAD_LIST_UPLOAD_FAILURE:
        return {
            ...state
        }

    default:
        return state;
    }
 
    return state;
}

export const UploadListActions = {
    uploadListFileAdd: (payload) => {
        return {
            payload: payload,
            type: UPLOAD_LIST_FILE_ADDED
            };
    },
    uploadListUploadStart: () => ({  
        payload: null,
        type: UPLOAD_LIST_UPLOAD_STARTED
    }),
    uploadListUploadSuccessful: () => ({  
        payload: null,
        type: UPLOAD_LIST_UPLOAD_SUCCESSFUL
    }),
    uploadListUploadFailure: () => ({  
        payload: null,
        type: UPLOAD_LIST_UPLOAD_FAILURE
    })
}


export function uploadAction() {
    return async(dispatch, getstate) => {
        dispatch(UploadListActions.uploadListUploadStart())
        try {
            let fileList = getstate().upload.fileList;
            for (var i in fileList) {
                const formData = new FormData();
                formData.append('file', fileList[i]);
                const options = {
                    method: 'POST',
                    body: formData,
                }   ;
                await fetch('./api/file/upload', options);
            }
            dispatch(UploadListActions.uploadListUploadSuccessful())
        }
        catch (e) {
            console.log(e)
            dispatch(UploadListActions.uploadListUploadFailure())
        }
    }
}