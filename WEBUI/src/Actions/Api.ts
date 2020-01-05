import { ActionCreator, Dispatch, Action } from 'redux';
import { ThunkAction } from 'redux-thunk';
import axios, { AxiosResponse } from 'axios';
import {
    RootActions,
    RootState,
    Program,
} from '../Types';
import {
    StartFetchAction,
    FailureFetchAction,
    SuccessFetchAction,
    API_ACTIONS,
} from '../Types/Api';

class ApiActionCreator {
    private startFetch = (): StartFetchAction => ({
        type: API_ACTIONS.START,
    });

    private failureFetch = (payload: string): FailureFetchAction => ({
        type: API_ACTIONS.FAILURE,
        payload: { message: payload },
    });

    private successFetch = (payload: Array<Program>): SuccessFetchAction => ({
        type: API_ACTIONS.SUCCESS,
        payload,
    });

    public getData = (query: string): ThunkAction<void, RootState, undefined, RootActions> => async (dispatch: Dispatch<Action>) => {
        dispatch(this.startFetch());
        const res = await axios.get('https://radio.sun-yryr.com/api/search', {
            params: {
                q: query,
            },
        }).then((response) => response.data).catch((e) => {
            dispatch(this.failureFetch(e));
        });
        dispatch(this.successFetch(res));
    }
}

export const apiActionCreator = new ApiActionCreator();
