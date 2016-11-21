import user from './user';
import exchange from './exchange';
import { combineReducers } from 'redux';
import { routerReducer as routing } from 'react-router-redux';

const rootReducer = combineReducers({
  routing,
  user,
  exchange,
});

export default rootReducer;
