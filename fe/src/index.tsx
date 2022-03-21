import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import * as serviceWorker from './serviceWorker';
import { ApiClientContext } from './services/apiContext';
import { ApiClient } from './services/apiClient';
import axios from "axios";

const client = new ApiClient(axios.create(), `${process.env.REACT_APP_API_HOST}`);

ReactDOM.render((
    <ApiClientContext.Provider value={client}>
        <App/>
    </ApiClientContext.Provider>
), document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
