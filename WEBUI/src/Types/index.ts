import { ApiState, ApiActions } from './Api';

// export {
//     StartFetchAction,
//     FailureFetchAction,
//     SuccessFetchAction,
//     API_ACTIONS,
// } from './Api';

export {
    Program,
} from './Main';

export interface RootState {
    Api: ApiState,
}

export type RootActions = ApiActions;
