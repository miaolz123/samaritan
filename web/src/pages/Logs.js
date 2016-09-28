import React from 'react';
import { Tag, Button, Modal, Table, Select, notification } from 'antd';
import keys from 'lodash.keys';
import axios from 'axios';

import config from '../config';

class Logs extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
      innerWidth: window.innerWidth || 1280,
      windowHeight: window.innerHeight || 720,
      loading: false,
      pagination: {
        pageSize: 20,
        current: 1,
        total: 0,
      },
      filters: {},
      tableData: [],
    };

    this.handleRefresh = this.handleRefresh.bind(this);
    this.fetchLogs = this.fetchLogs.bind(this);
    this.deleteLogs = this.deleteLogs.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
  }

  componentWillMount() {
    const { trader } = this.props;

    if (!trader) {
      Modal.error({
        title: 'Error',
        content: 'No trader found !',
        onOk: () => {
          window.location.href = window.location.href;
        },
      });
    }
    this.fetchLogs();
  }

  handleRefresh() {
    this.fetchLogs();
  }

  fetchLogs(pagination, filters) {
    const { trader } = this.props;

    if (!pagination) {
      pagination = this.state.pagination;
    }
    if (!filters) {
      filters = this.state.filters;
    }
    this.setState({ loading: true });
    axios.post(`${config.api}/logs`, { trader, pagination, filters }, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        this.setState({ loading: false });
        if (response.data.success) {
          const thisPagination = this.state.pagination;
          const { data, total } = response.data;

          thisPagination.total = total;
          this.setState({
            pagination: thisPagination,
            tableData: data,
          });
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
      }, (response) => {
        this.setState({ loading: false });
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        } else {
          notification['error']({
            message: 'Error',
            description: String(response),
            duration: null,
          });
        }
      });
  }

  deleteLogs(value) {
    const { trader } = this.props;

    axios.delete(`${config.api}/logs?id=${trader.ID}&type=${value}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.handleRefresh();
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
      }, (response) => {
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        } else {
          notification['error']({
            message: 'Error',
            description: String(response),
            duration: null,
          });
        }
      });
  }

  handleTableChange(pagination, filters) {
    this.setState({
      pagination,
      filters,
    });
    this.fetchLogs(pagination, filters);
  }

  render() {
    const { tableData, innerWidth, windowHeight } = this.state;
    const exchangeTypes = config.exchangeTypes.map(t => ({ text: t, value: `'${t}'` }));
    const logTypes = keys(config.logTypes).map(k => ({ text: config.logTypes[k], value: k }));
    const columns = [{
      title: 'Time',
      dataIndex: 'Time',
      fixed: 'left',
      width: 120,
    }, {
      title: 'Exchange',
      dataIndex: 'ExchangeType',
      fixed: 'left',
      filters: exchangeTypes,
      width: 100,
      render: text => text && <Tag color={text === 'global' ? '' : 'blue'}>{text}</Tag>,
    }, {
      title: 'Type',
      dataIndex: 'Type',
      filters: logTypes,
      width: 100,
      render: text => <Tag
        color={text < 0 ? 'red' : text < 1 ? '' : text < 2 ? 'yellow' : 'blue'}
      >{config.logTypes[text]}</Tag>,
    }, {
      title: 'Price',
      dataIndex: 'Price',
      width: 100,
      render: text => text === 0.0 ? '' : text.toFixed(3),
    }, {
      title: 'Amount',
      dataIndex: 'Amount',
      width: 100,
      render: text => text === 0.0 ? '' : text.toFixed(3),
    }, {
      title: 'Message',
      dataIndex: 'Message',
    }];

    return (
      <div>
        <div style={{ marginBottom: 16, textAlign: 'right' }}>
          <Button style={{ marginRight: 5 }} type="primary" onClick={this.props.goBack}>Go Back</Button>
          <Button style={{ marginRight: 5 }} onClick={this.handleRefresh}>Refresh</Button>
          <Select
            placeholder="Delete"
            onSelect={this.deleteLogs}
            dropdownMatchSelectWidth={false}
            style={{ marginRight: 5, width: 80, textAlign: 'left' }}
          >
            <Select.Option value="-1" disabled>Earlier Than</Select.Option>
            <Select.Option value="0">Last Run</Select.Option>
            <Select.Option value="1">One Day</Select.Option>
            <Select.Option value="2">One Week</Select.Option>
            <Select.Option value="3">One Month</Select.Option>
          </Select>
          <Tag>Total: {this.state.pagination.total}</Tag>
        </div>
        <Table
          size="middle"
          scroll={{x: innerWidth > 1000 ? innerWidth - 250 : 750, y: windowHeight > 500 ? windowHeight - 230 : 500}}
          columns={columns}
          dataSource={tableData}
          pagination={this.state.pagination}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }
}

export default Logs;
