import user from './user';
import { combineReducers } from 'redux';
import { routerReducer as routing } from 'react-router-redux';

const rootReducer = combineReducers({
  routing,
  user,
});

export default rootReducer;
