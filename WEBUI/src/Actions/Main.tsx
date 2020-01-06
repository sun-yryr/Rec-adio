import { Action } from 'redux';
import {
    Program,
    AddFrontAction,
    MAIN_ACTIONS,
    AddQueueAction,
    SkipAction,
} from '../Types/Main';

class MainActionCreator {
    public addFront = (prog: Program): AddFrontAction => ({
        type: MAIN_ACTIONS.ADD_FRONT,
        payload: { program: prog },
    });

    public addQueue = (prog: Program): AddQueueAction => ({
        type: MAIN_ACTIONS.ADD_QUEUE,
        payload: { program: prog },
    });

    public skip = (): SkipAction => ({
        type: MAIN_ACTIONS.SKIP,
    });
}
export const mainActionCreator = new MainActionCreator();
