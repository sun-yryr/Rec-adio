import { Action } from 'redux';
import { Program } from './Main';

export interface ApiState {
    onFetch: boolean,
    error?: string,
    data: Array<Program>
}

export interface StartFetchAction extends Action {
    type: 'START_FETCH'
}

export interface FailureFetchAction extends Action {
    type: 'FAILURE_FETCH';
    payload: { message: string };
}
