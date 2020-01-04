import React from 'react';

import logo from './logo.svg';
import './App.css';
import AudioProvider from './conponents/AudioProvider';

const client = new ApolloClient({
    uri: 'https://48p1r2roz4.sse.codesandbox.io',
});

const App: React.FC = () => (
    <ApolloProvider client={client}>
        <div>
            <h2>Apollo app</h2>
        </div>
    </ApolloProvider>
);

export default App;
