import { Reducer } from 'redux';
import { ApiState, ApiActions, API_ACTIONS } from '../Types/Api';

const initState: ApiState = {
    onFetch: false,
    data: [],
};

export const apiReducer: Reducer<ApiState, ApiActions> = (
    state: ApiState = initState,
    action: ApiActions,
) => {
    switch (action.type) {
        case API_ACTIONS.START: {
            return {
                ...state,
                onFetch: true,
            };
        }
        case API_ACTIONS.FAILURE: {
            return {
                ...state,
                onFetch: false,
                error: action.payload.message,
            };
        }
        case API_ACTIONS.SUCCESS: {
            return {
                ...state,
                onFetch: false,
                data: action.payload,
                error: undefined,
            };
        }
        default: {
            return state;
        }
    }
};
