import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const ALGORITHM_INIT = {
  loading: false,
  total: 0,
  list: [],
  cache: {},
  message: '',
};

function algorithm(state = ALGORITHM_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: '',
      });
    case actions.ALGORITHM_LIST_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.ALGORITHM_LIST_SUCCESS:
      return assign({}, state, {
        loading: false,
        total: action.total,
        list: action.list,
      });
    case actions.ALGORITHM_LIST_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.ALGORITHM_PUT_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.ALGORITHM_PUT_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.ALGORITHM_PUT_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.ALGORITHM_DELETE_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.ALGORITHM_DELETE_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.ALGORITHM_DELETE_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.ALGORITHM_CACHE:
      return assign({}, state, {
        cache: action.cache,
      });
    default:
      return state;
  }
}

export default algorithm;
