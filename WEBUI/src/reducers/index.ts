import { combineReducers } from 'redux';
import { RootState } from '../Types';
import { apiReducer } from './Api';


export const rootReducer = combineReducers<RootState>({
    Api: apiReducer,
});
