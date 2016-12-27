import * as actions from '../constants/actions';
import { Client } from 'hprose-js';

// List

function traderListRequest() {
  return { type: actions.TRADER_LIST_REQUEST };
}

function traderListSuccess(algorithmId, list) {
  return { type: actions.TRADER_LIST_SUCCESS, algorithmId, list };
}

function traderListFailure(message) {
  return { type: actions.TRADER_LIST_FAILURE, message };
}

export function TraderList(algorithmId) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(traderListRequest());
    if (!cluster || !token) {
      dispatch(traderListFailure('No authorization'));
      return;
    }

    const client = Client.create(`${cluster}/api`, { Trader: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Trader.List(algorithmId, (resp) => {
      if (resp.success) {
        dispatch(traderListSuccess(algorithmId, resp.data));
      } else {
        dispatch(traderListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(traderListFailure('Server error'));
      console.log('【Hprose】Trader.List Error:', resp, err);
    });
  };
}

// Put

function traderPutRequest() {
  return { type: actions.TRADER_PUT_REQUEST };
}

function traderPutSuccess() {
  return { type: actions.TRADER_PUT_SUCCESS };
}

function traderPutFailure(message) {
  return { type: actions.TRADER_PUT_FAILURE, message };
}

export function TraderPut(req) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(traderPutRequest());
    if (!cluster || !token) {
      dispatch(traderPutFailure('No authorization'));
      return;
    }

    const client = Client.create(`${cluster}/api`, { Trader: ['Put'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Trader.Put(req, (resp) => {
      if (resp.success) {
        dispatch(traderPutSuccess());
        dispatch(TraderList(req.algorithmId));
      } else {
        dispatch(traderPutFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(traderPutFailure('Server error'));
      console.log('【Hprose】Trader.Put Error:', resp, err);
    });
  };
}

// Delete

function traderDeleteRequest() {
  return { type: actions.TRADER_DELETE_REQUEST };
}

function traderDeleteSuccess() {
  return { type: actions.TRADER_DELETE_SUCCESS };
}

function traderDeleteFailure(message) {
  return { type: actions.TRADER_DELETE_FAILURE, message };
}

export function TraderDelete(req) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(traderDeleteRequest());
    if (!cluster || !token) {
      dispatch(traderDeleteFailure('No authorization'));
      return;
    }

    const client = Client.create(`${cluster}/api`, { Trader: ['Delete'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Trader.Delete(req, (resp) => {
      if (resp.success) {
        dispatch(traderDeleteSuccess());
        dispatch(TraderList(req.algorithmId));
      } else {
        dispatch(traderDeleteFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(traderDeleteFailure('Server error'));
      console.log('【Hprose】Trader.Delete Error:', resp, err);
    });
  };
}

// Switch

function traderSwitchRequest() {
  return { type: actions.TRADER_SWITCH_REQUEST };
}

function traderSwitchSuccess() {
  return { type: actions.TRADER_SWITCH_SUCCESS };
}

function traderSwitchFailure(message) {
  return { type: actions.TRADER_SWITCH_FAILURE, message };
}

export function TraderSwitch(req) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(traderSwitchRequest());
    if (!cluster || !token) {
      dispatch(traderSwitchFailure('No authorization'));
      return;
    }

    const client = Client.create(`${cluster}/api`, { Trader: ['Switch'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Trader.Switch(req, (resp) => {
      if (resp.success) {
        dispatch(traderSwitchSuccess());
        dispatch(TraderList(req.algorithmId));
      } else {
        dispatch(traderSwitchFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(traderSwitchFailure('Server error'));
      console.log('【Hprose】Trader.Switch Error:', resp, err);
    });
  };
}

// Cache

export function TraderCache(cache) {
  return { type: actions.TRADER_CACHE, cache };
}
