import { Action } from 'redux';

export interface Program {
    id: number,
    title: string,
    pfm: string,
    recTimestamp: Date,
    station: string,
    uri: string,
    info?: string,
}

export interface MainState {
    nowProgram?: Program,
    audioQueue: Array<Program>,
}

export enum MAIN_ACTIONS {
    ADD_FRONT = 'ADD_FRONT',
    ADD_QUEUE = 'ADD_QUEUE',
    SKIP = 'SKIP',
}

export interface AddFrontAction extends Action {
    type: MAIN_ACTIONS.ADD_FRONT;
    payload: { program: Program };
}

export interface AddQueueAction extends Action {
    type: MAIN_ACTIONS.ADD_QUEUE;
    payload: { program: Program };
}

export interface SkipAction extends Action {
    type: MAIN_ACTIONS.SKIP;
}

export type MainActions = AddFrontAction & AddQueueAction & SkipAction;
