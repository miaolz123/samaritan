import * as actions from '../constants/actions';
import { Client } from 'hprose-html5/dist/hprose-html5';

// List

function logListRequest() {
  return { type: actions.LOG_LIST_REQUEST };
}

function logListSuccess(total, list) {
  return { type: actions.LOG_LIST_SUCCESS, total, list };
}

function logListFailure(message) {
  return { type: actions.LOG_LIST_FAILURE, message };
}

export function LogList(trader, pagination, filters) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(logListRequest());
    if (!cluster || !token) {
      dispatch(logListFailure('No authorization'));
      return;
    }

    const client = Client.create(cluster, { Log: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Log.List(trader, pagination, filters, (resp) => {
      if (resp.success) {
        dispatch(logListSuccess(resp.data.total, resp.data.list));
      } else {
        dispatch(logListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(logListFailure('Server error'));
      console.log('【Hprose】Log.List Error:', resp, err);
    });
  };
}
