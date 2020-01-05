import { Action } from 'redux';
import { Program } from './Main';

export interface ApiState {
    onFetch: boolean,
    error?: string,
    data: Array<Program>
}

export enum API_ACTIONS {
    START = 'START_FETCH',
    FAILURE = 'FAILURE_FETCH',
    SUCCESS = 'SUCCESS_FETCH',
}

export interface StartFetchAction extends Action {
    type: API_ACTIONS.START
}

export interface FailureFetchAction extends Action {
    type: API_ACTIONS.FAILURE;
    payload: { message: string };
}

export interface SuccessFetchAction extends Action {
    type: API_ACTIONS.SUCCESS;
    payload: Array<Program>;
}

export type ApiActions = StartFetchAction & FailureFetchAction & SuccessFetchAction;
