import React from 'react';
import { Route, IndexRoute } from 'react-router';
import App from './containers/App';
import Dashboard from './containers/Dashboard';
import Login from './containers/Login';
import User from './containers/User';
import Exchange from './containers/Exchange';
import Algorithm from './containers/Algorithm';
import AlgorithmEdit from './containers/AlgorithmEdit';
import AlgorithmLog from './containers/AlgorithmLog';

export default (
  <Route>
    <Route path="/" component={App}>
      <IndexRoute component={Dashboard} />
      <Route path="/user" component={User} />
      <Route path="/exchange" component={Exchange} />
      <Route path="/algorithm" component={Algorithm} />
      <Route path="/algorithmEdit" component={AlgorithmEdit} />
      <Route path="/algorithmLog" component={AlgorithmLog} />
    </Route>
    <Route path="/login" component={Login} />
  </Route>
);
