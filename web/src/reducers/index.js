import user from './user';
import exchange from './exchange';
import algorithm from './algorithm';
import { combineReducers } from 'redux';
import { routerReducer as routing } from 'react-router-redux';

const rootReducer = combineReducers({
  routing,
  user,
  exchange,
  algorithm,
});

export default rootReducer;
