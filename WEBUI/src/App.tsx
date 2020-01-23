import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import styled from 'styled-components';
import Top from './components/Top';
import Search from './components/Search';
import ResultTable from './components/ResultTable';
import { rootReducer } from './reducers';
import AudioProvider from './components/AudioProvider';

const store = createStore(
    rootReducer,
    applyMiddleware(thunk),
);

const App: React.FC = () => (
    <Provider store={store}>
        <Root>
            <RouteRoot>
                <BrowserRouter>
                    <Route exact path="/" component={Top} />
                    <Route exact path="/search" component={Search} />
                    <Route exact path="/result" component={ResultTable} />
                </BrowserRouter>
            </RouteRoot>
            <AudioProvider />
        </Root>
    </Provider>
);

export default App;

const Root = styled.div`
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
`;

const RouteRoot = styled.div`
    height: 86%;
    overflow: scroll;
`;
