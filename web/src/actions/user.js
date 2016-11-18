import * as actions from '../constants/actions';

// Login

function userLoginRequest() {
  return { type: actions.USER_LOGIN_REQUEST };
}

function userLoginSuccess(token) {
  return { type: actions.USER_LOGIN_SUCCESS, token: token };
}

function userLoginFailure(message) {
  return { type: actions.USER_LOGIN_FAILURE, message: message };
}

export function UserLogin(cluster, username, password) {
  return (dispatch) => {
    dispatch(userLoginRequest());
    const client = hprose.Client.create(cluster, { User: ['Login'] });
    client.User.Login(username, password, (resp) => {
      if (resp.success) {
        localStorage.setItem('cluster', cluster);
        localStorage.setItem('token', resp.data);
        dispatch(userLoginSuccess(resp.data));
      } else {
        dispatch(userLoginFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userLoginFailure('Can'));
      console.log('【Hprose】UserLogin Error:', resp, err);
    });
  };
}
