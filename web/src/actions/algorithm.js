import * as actions from '../constants/actions';

// List

function algorithmListRequest() {
  return { type: actions.ALGORITHM_LIST_REQUEST };
}

function algorithmListSuccess(total, list) {
  return { type: actions.ALGORITHM_LIST_SUCCESS, total, list };
}

function algorithmListFailure(message) {
  return { type: actions.ALGORITHM_LIST_FAILURE, message };
}

export function AlgorithmList(size, page, order) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(algorithmListRequest());
    if (!cluster || !token) {
      dispatch(algorithmListFailure('No authorization'));
      return;
    }

    const client = hprose.Client.create(cluster, { Algorithm: ['List'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Algorithm.List(size, page, order, (resp) => {
      if (resp.success) {
        dispatch(algorithmListSuccess(resp.data.total, resp.data.list));
      } else {
        dispatch(algorithmListFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(algorithmListFailure('Server error'));
      console.log('【Hprose】Algorithm.List Error:', resp, err);
    });
  };
}

// Put

function algorithmPutRequest() {
  return { type: actions.ALGORITHM_PUT_REQUEST };
}

function algorithmPutSuccess() {
  return { type: actions.ALGORITHM_PUT_SUCCESS };
}

function algorithmPutFailure(message) {
  return { type: actions.ALGORITHM_PUT_FAILURE, message };
}

export function AlgorithmPut(req) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(algorithmPutRequest());
    if (!cluster || !token) {
      dispatch(algorithmPutFailure('No authorization'));
      return;
    }

    const client = hprose.Client.create(cluster, { Algorithm: ['Put'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Algorithm.Put(req, (resp) => {
      if (resp.success) {
        dispatch(algorithmPutSuccess());
      } else {
        dispatch(algorithmPutFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(algorithmPutFailure('Server error'));
      console.log('【Hprose】Algorithm.Put Error:', resp, err);
    });
  };
}

// Delete

function algorithmDeleteRequest() {
  return { type: actions.ALGORITHM_DELETE_REQUEST };
}

function algorithmDeleteSuccess() {
  return { type: actions.ALGORITHM_DELETE_SUCCESS };
}

function algorithmDeleteFailure(message) {
  return { type: actions.ALGORITHM_DELETE_FAILURE, message };
}

export function AlgorithmDelete(ids, size, page, order) {
  return (dispatch, getState) => {
    const cluster = localStorage.getItem('cluster');
    const token = localStorage.getItem('token');

    dispatch(algorithmDeleteRequest());
    if (!cluster || !token) {
      dispatch(algorithmDeleteFailure('No authorization'));
      return;
    }

    const client = hprose.Client.create(cluster, { Algorithm: ['Delete'] });

    client.setHeader('Authorization', `Bearer ${token}`);
    client.Algorithm.Delete(ids, (resp) => {
      if (resp.success) {
        dispatch(AlgorithmList(size, page, order));
        dispatch(algorithmDeleteSuccess());
      } else {
        dispatch(algorithmDeleteFailure(resp.message));
      }
    }, (resp, err) => {
      dispatch(algorithmDeleteFailure('Server error'));
      console.log('【Hprose】Algorithm.Delete Error:', resp, err);
    });
  };
}

// Cache

export function AlgorithmCache(cache) {
  return { type: actions.ALGORITHM_CACHE, cache };
}
