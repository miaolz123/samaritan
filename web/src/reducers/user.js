import * as actions from '../constants/actions';
import merge from 'lodash/merge';

const USER_INIT = { loading: false, data: null, token: '', message: '' };

function user(state = USER_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return merge({}, state, {
        loading: false,
        message: '',
      });
    case actions.USER_LOGIN_REQUEST:
      return merge({}, state, {
        loading: true,
      });
    case actions.USER_LOGIN_SUCCESS:
      return merge({}, state, {
        loading: false,
        token: action.token,
      });
    case actions.USER_LOGIN_FAILURE:
      return merge({}, state, {
        loading: false,
        message: action.message,
      });
    default:
      return state;
  }
}

export default user;
