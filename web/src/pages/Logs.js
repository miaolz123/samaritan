import React from 'react';
import { Tag, Button, Modal, Table, notification } from 'antd';
import keys from 'lodash.keys';
import axios from 'axios';

import config from '../config';

class Logs extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
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
        if (response.data.success) {
          const thisPagination = this.state.pagination;
          const { data, total } = response.data;

          thisPagination.total = total;
          this.setState({
            loading: false,
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
    const { tableData, windowHeight } = this.state;
    const exchangeTypes = config.exchangeTypes.map(t => ({ text: t, value: `'${t}'` }));
    const logTypes = keys(config.logTypes).map(k => ({ text: config.logTypes[k], value: k }));
    exchangeTypes.push({ text: 'globle', value: "''" });
    const columns = [{
      title: 'Time',
      dataIndex: 'Time',
      width: '10%',
    }, {
      title: 'Exchange',
      dataIndex: 'ExchangeType',
      filters: exchangeTypes,
      width: '15%',
      render: text => text ? <Tag color="blue">{text}</Tag> : <Tag>globle</Tag>,
    }, {
      title: 'Type',
      dataIndex: 'Type',
      filters: logTypes,
      width: '10%',
      render: text => <Tag
        color={text < 0 ? 'red' : text < 1 ? '' : text < 2 ? 'yellow' : 'blue'}
      >{config.logTypes[text]}</Tag>,
    }, {
      title: 'Price',
      dataIndex: 'Price',
      width: '10%',
    }, {
      title: 'Amount',
      dataIndex: 'Amount',
      width: '10%',
    }, {
      title: 'Message',
      dataIndex: 'Message',
    }];

    return (
      <div>
        <div style={{ marginBottom: 16, textAlign: 'right' }}>
          <Button style={{ marginRight: 5 }} type="primary" onClick={this.props.goBack}>Go Back</Button>
          <Button style={{ marginRight: 10 }} onClick={this.handleRefresh}>Refresh</Button>
          <Tag>Total: {this.state.pagination.total}</Tag>
        </div>
        <Table
          size="middle"
          scroll={{y: windowHeight > 500 ? windowHeight - 230 : 500}}
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
