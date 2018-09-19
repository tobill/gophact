export function uploadReduce(state = {}) {
    return state;
}


export function uploadAction(fileList) {
    return async(dispatch, getstate) => {

        
        console.log(fileList);
        for (var i = 0, numFiles = fileList.length; i < numFiles; i++) {
            const formData = new FormData();
            formData.append('file', fileList.item(i));
            const options = {
                method: 'POST',
                body: formData,
            };
            await fetch('./upload', options);
        }

    }
}