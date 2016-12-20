import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const LOG_INIT = {
  loading: false,
  total: 0,
  list: [],
  message: '',
};

function log(state = LOG_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: '',
      });
    case actions.LOG_LIST_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.LOG_LIST_SUCCESS:
      return assign({}, state, {
        loading: false,
        total: action.total,
        list: action.list,
      });
    case actions.LOG_LIST_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    default:
      return state;
  }
}

export default log;
