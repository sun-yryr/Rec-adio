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
    loading: boolean,
    nowplaying?: HTMLAudioElement,
    audioQueue: Array<HTMLAudioElement>,
}
