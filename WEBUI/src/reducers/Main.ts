import { Reducer } from 'redux';
import { MainState, MainActions, MAIN_ACTIONS } from '../Types/Main';

const initState: MainState = {
    audioQueue: [],
};

export const mainReducer: Reducer<MainState, MainActions> = (
    state: MainState = initState,
    action: MainActions,
) => {
    switch (action.type) {
        case MAIN_ACTIONS.ADD_FRONT: {
            return {
                ...state,
                nowProgram: action.payload.program,
            };
        }
        case MAIN_ACTIONS.ADD_QUEUE: {
            const { audioQueue, ...other } = state;
            const tmp = audioQueue.map((i) => i);
            tmp.push(action.payload.program);
            return {
                ...other,
                audioQueue: tmp,
            };
        }
        case MAIN_ACTIONS.SKIP: {
            const { audioQueue } = state;
            const next = audioQueue.shift();
            return {
                nowProgram: next,
                audioQueue,
            };
        }
        default: {
            return state;
        }
    }
};
