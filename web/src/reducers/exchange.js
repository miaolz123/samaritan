import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const EXCHANGE_INIT = {
  loading: false,
  types: [],
  total: 0,
  list: [],
  message: '',
};

function exchange(state = EXCHANGE_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: '',
      });
    case actions.EXCHANGE_TYPES_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.EXCHANGE_TYPES_SUCCESS:
      return assign({}, state, {
        loading: false,
        types: action.types,
      });
    case actions.EXCHANGE_TYPES_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.EXCHANGE_LIST_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.EXCHANGE_LIST_SUCCESS:
      return assign({}, state, {
        loading: false,
        total: action.total,
        list: action.list,
      });
    case actions.EXCHANGE_LIST_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.EXCHANGE_PUT_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.EXCHANGE_PUT_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.EXCHANGE_PUT_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.EXCHANGE_DELETE_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.EXCHANGE_DELETE_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.EXCHANGE_DELETE_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    default:
      return state;
  }
}

export default exchange;
