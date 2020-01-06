import { ApiState } from './Api';

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
    now: {
        program?: Program,
    },
    audioQueue: Array<Program>,
}

export enum MAIN_ACTIONS {
    ADD_PLAY = 'ADD_PLAY',
    ADD_QUEUE = 'ADD_QUEUE',
}