import { ApiState, ApiActions } from './Api';
import { MainState } from './Main';

// export {
//     StartFetchAction,
//     FailureFetchAction,
//     SuccessFetchAction,
//     API_ACTIONS,
// } from './Api';

export interface RootState {
    Api: ApiState,
    Main: MainState,
}

export type RootActions = ApiActions;
