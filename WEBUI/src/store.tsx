import {
    Action,
    createStore,
    Store,
    applyMiddleware,
} from 'redux';

export interface IRootState {
    loading: boolean,
    nowplaying?: HTMLAudioElement,
    audioQueue: Array<HTMLAudioElement>,
}
