import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import Top from './components/Top';
import Search from './components/Search';
import ResultTable from './components/ResultTable';
import { rootReducer } from './reducers';

const store = createStore(
    rootReducer,
    applyMiddleware(thunk),
);

const App: React.FC = () => (
    <Provider store={store}>
        <BrowserRouter>
            <Route exact path="/" component={Top} />
            <Route exact path="/search" component={Search} />
            <Route exact path="/result" component={ResultTable} />
        </BrowserRouter>
    </Provider>
);

export default App;
