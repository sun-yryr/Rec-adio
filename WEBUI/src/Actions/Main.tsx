import {
    Program,
    AddFrontAction,
    MAIN_ACTIONS,
    AddQueueAction,
    SkipAction,
    SetPassAction,
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

    public setPassword = (pass: string): SetPassAction => ({
        type: MAIN_ACTIONS.SET_PASS,
        payload: { password: pass },
    });
}
export const mainActionCreator = new MainActionCreator();
