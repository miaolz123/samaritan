import user from './user';
import exchange from './exchange';
import algorithm from './algorithm';
import trader from './trader';
import log from './log';
import { combineReducers } from 'redux';
import { routerReducer as routing } from 'react-router-redux';

const rootReducer = combineReducers({
  routing,
  user,
  exchange,
  algorithm,
  trader,
  log,
});

export default rootReducer;
