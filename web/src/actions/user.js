import * as actions from '../constants/actions';

// Login

function userLoginRequest() {
  return { type: actions.USER_LOGIN_REQUEST };
}

function userLoginSuccess(token, cluster) {
  return { type: actions.USER_LOGIN_SUCCESS, token, cluster };
}

function userLoginFailure(message) {
  return { type: actions.USER_LOGIN_FAILURE, message };
}

export function UserLogin(cluster, username, password) {
  return (dispatch, getState) => {
    dispatch(userLoginRequest());
    const client = hprose.Client.create(cluster, { User: ['Login'] });
    client.User.Login(username, password, (resp) => {
      if (resp.success) {
        dispatch(userLoginSuccess(resp.data, cluster));
      } else {
        dispatch(userLoginFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userLoginFailure('Server error'));
      console.log('【Hprose】User.Login Error:', resp, err);
    });
  };
}

// Get

function userGetRequest() {
  return { type: actions.USER_GET_REQUEST };
}

function userGetSuccess(data) {
  return { type: actions.USER_GET_SUCCESS, data };
}

function userGetFailure(message) {
  return { type: actions.USER_GET_FAILURE, message };
}

export function UserGet() {
  return (dispatch, getState) => {
    dispatch(userGetRequest());
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    if (!cluster || !token) {
      dispatch(userGetFailure('No authorization'));
      return;
    }

    const client = hprose.Client.create(cluster, { User: ['Get'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.User.Get(null, (resp) => {
      console.log('【Hprose】User.Get OK!', resp);
      if (resp.success) {
        dispatch(userGetSuccess(resp.message.data));
      } else {
        dispatch(userGetFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userGetFailure('Server error'));
      console.log('【Hprose】User.Get Error:', resp, err);
    });
  };
}

// Logout

export function Logout() {
  return { type: actions.LOGOUT };
}
