import { combineReducers, } from "redux";
import { uploadReduce } from "./upload";
import { reducer as formReduce} from "redux-form"

export const reducer = combineReducers({
    upload: uploadReduce,
    form: formReduce,
});