import { Dispatch, Action } from 'redux';
import { ThunkAction } from 'redux-thunk';
import axios from 'axios';
import NodeRSA from 'node-rsa';
import { RootActions, RootState } from '../Types';
import {
    StartFetchAction,
    FailureFetchAction,
    SuccessFetchAction,
    API_ACTIONS,
} from '../Types/Api';
import { Program } from '../Types/Main';

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

    public getData = (query: string, pass: string): ThunkAction<void, RootState, undefined, RootActions> => async (dispatch: Dispatch<Action>) => {
        dispatch(this.startFetch());
        const publicKey = `-----BEGIN PUBLIC KEY-----
        MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApsrNku3xhhe3b2n0SH/G
        HZJhZc6J8K4XtJWvxtuTYM/1N7EKIVOCz9nhghtdBGBBoLQ5hKHuPZ9eSxXcrnli
        +cjgUA0DQD4+CeO9WaDPIY1mlUR5xLqWhaFyub3FvdCd7l3ak7dCadBdpSDG30Bp
        7LjoDPDuUI15eRTP1rFqVL+3XheEsx9vC7yVRrK7lBd9jHspHkYylcGkr1IpLWwI
        4/h8gfiblx6Z4f/qAWW2p+u/yQBB+SMd62WEIbzZ6jFSzxSSMAb0LD7eQet7FA/B
        67eO3pa6gSuggIqGZ6X1G9cIsSGaBDzZE+GubtgjzxhkfgJvhohTjHfWq8JXy2lb
        aQIDAQAB
        -----END PUBLIC KEY-----`;
        const rsa = new NodeRSA(publicKey, 'public');
        const encryptText = rsa.encrypt(pass, 'base64');
        console.log(encryptText);
        const res: Array<any> = await axios.get('https://radio.sun-yryr.com/api/search', {
            params: {
                q: query,
                password: encryptText,
            },
        }).then((response) => response.data).catch((e) => {
            dispatch(this.failureFetch(e));
        });
        const r: Array<Program> = res.map((value) => ({
            ...value,
            recTimestamp: value['rec-timestamp'],
        }));
        dispatch(this.successFetch(r));
    }
}

export const apiActionCreator = new ApiActionCreator();
