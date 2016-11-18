import rootReducer from '../reducers';
import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';

export default function configureStore(initialState) {
  return createStore(
    rootReducer,
    initialState,
    compose(applyMiddleware(thunk))
  );
}
