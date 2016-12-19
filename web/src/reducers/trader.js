import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const TRADER_INIT = {
  loading: false,
  map: {},
  message: '',
};

function trader(state = TRADER_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: '',
      });
    case actions.TRADER_LIST_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.TRADER_LIST_SUCCESS:
      const { map } = state;

      map[action.algorithmId] = action.list;

      return assign({}, state, {
        loading: false,
        map,
      });
    case actions.TRADER_LIST_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.TRADER_PUT_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.TRADER_PUT_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.TRADER_PUT_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.TRADER_DELETE_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.TRADER_DELETE_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.TRADER_DELETE_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    default:
      return state;
  }
}

export default trader;
