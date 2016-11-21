import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const USER_INIT = {
  loading: false,
  status: 0,
  data: null,
  total: 0,
  list: [],
  cluster: '',
  token: '',
  message: '',
};

function user(state = USER_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: '',
      });
    case actions.USER_LOGIN_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.USER_LOGIN_SUCCESS:
      localStorage.setItem('cluster', action.cluster);
      localStorage.setItem('token', action.token);
      return assign({}, state, {
        loading: false,
        status: 1,
        cluster: action.cluster,
        token: action.token,
      });
    case actions.USER_LOGIN_FAILURE:
      localStorage.removeItem('cluster');
      localStorage.removeItem('token');
      return assign({}, state, {
        loading: false,
        status: -1,
        message: action.message,
      });
    case actions.USER_GET_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.USER_GET_SUCCESS:
      return assign({}, state, {
        loading: false,
        status: 1,
        data: action.data,
      });
    case actions.USER_GET_FAILURE:
      localStorage.removeItem('cluster');
      localStorage.removeItem('token');
      return assign({}, state, {
        loading: false,
        status: -1,
        message: action.message,
      });
    case actions.USER_LIST_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.USER_LIST_SUCCESS:
      return assign({}, state, {
        loading: false,
        total: action.total,
        list: action.list,
      });
    case actions.USER_LIST_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.USER_PUT_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.USER_PUT_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.USER_PUT_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.USER_DELETE_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.USER_DELETE_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.USER_DELETE_FAILURE:
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.LOGOUT:
      localStorage.removeItem('cluster');
      localStorage.removeItem('token');
      return USER_INIT;
    default:
      return state;
  }
}

export default user;
