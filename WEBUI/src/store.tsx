import {
    Action,
    createStore,
    Store,
    applyMiddleware,
} from '../node_modules/redux';

export interface IRootState {
    loading: boolean,
    nowplaying?: HTMLAudioElement,
    audioQueue: Array<HTMLAudioElement>,
}
