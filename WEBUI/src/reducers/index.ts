import { combineReducers } from 'redux';
import { RootState } from '../Types';
import { apiReducer } from './Api';
import { mainReducer } from './Main';


export const rootReducer = combineReducers<RootState>({
    Api: apiReducer,
    Main: mainReducer,
});
