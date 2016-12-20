import * as actions from '../constants/actions';
import { Client } from 'hprose-html5/dist/hprose-html5';

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
    const client = Client.create(cluster, { User: ['Login'] });

    dispatch(userLoginRequest());
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
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(userGetRequest());
    if (!cluster || !token) {
      dispatch(userGetFailure('No authorization'));
      return;
    }

    const client = Client.create(cluster, { User: ['Get'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.User.Get(null, (resp) => {
      if (resp.success) {
        dispatch(userGetSuccess(resp.data));
      } else {
        dispatch(userGetFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userGetFailure('Server error'));
      console.log('【Hprose】User.Get Error:', resp, err);
    });
  };
}

// List

function userListRequest() {
  return { type: actions.USER_LIST_REQUEST };
}

function userListSuccess(total, list) {
  return { type: actions.USER_LIST_SUCCESS, total, list };
}

function userListFailure(message) {
  return { type: actions.USER_LIST_FAILURE, message };
}

export function UserList(size, page, order) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(userListRequest());
    if (!cluster || !token) {
      dispatch(userGetFailure('No authorization'));
      dispatch(userListFailure('No authorization'));
      return;
    }

    const client = Client.create(cluster, { User: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.User.List(size, page, order, (resp) => {
      if (resp.success) {
        dispatch(userListSuccess(resp.data.total, resp.data.list));
      } else {
        dispatch(userListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userListFailure('Server error'));
      console.log('【Hprose】User.List Error:', resp, err);
    });
  };
}

// Put

function userPutRequest() {
  return { type: actions.USER_PUT_REQUEST };
}

function userPutSuccess() {
  return { type: actions.USER_PUT_SUCCESS };
}

function userPutFailure(message) {
  return { type: actions.USER_PUT_FAILURE, message };
}

export function UserPut(req, password, size, page, order) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(userPutRequest());
    if (!cluster || !token) {
      dispatch(userGetFailure('No authorization'));
      dispatch(userPutFailure('No authorization'));
      return;
    }

    const client = Client.create(cluster, { User: ['Put'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.User.Put(req, password, (resp) => {
      if (resp.success) {
        dispatch(userPutSuccess());
        dispatch(UserList(size, page, order));
      } else {
        dispatch(userPutFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userPutFailure('Server error'));
      console.log('【Hprose】User.Put Error:', resp, err);
    });
  };
}

// Delete

function userDeleteRequest() {
  return { type: actions.USER_DELETE_REQUEST };
}

function userDeleteSuccess() {
  return { type: actions.USER_DELETE_SUCCESS };
}

function userDeleteFailure(message) {
  return { type: actions.USER_DELETE_FAILURE, message };
}

export function UserDelete(ids, size, page, order) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(userDeleteRequest());
    if (!cluster || !token) {
      dispatch(userGetFailure('No authorization'));
      dispatch(userDeleteFailure('No authorization'));
      return;
    }

    const client = Client.create(cluster, { User: ['Delete'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.User.Delete(ids, (resp) => {
      if (resp.success) {
        dispatch(userDeleteSuccess());
        dispatch(UserList(size, page, order));
      } else {
        dispatch(userDeleteFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(userDeleteFailure('Server error'));
      console.log('【Hprose】User.Delete Error:', resp, err);
    });
  };
}

// Logout

export function Logout() {
  return { type: actions.LOGOUT };
}
