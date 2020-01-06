import { ApiState, ApiActions } from './Api';

// export {
//     StartFetchAction,
//     FailureFetchAction,
//     SuccessFetchAction,
//     API_ACTIONS,
// } from './Api';

export interface RootState {
    Api: ApiState,
}

export type RootActions = ApiActions;
