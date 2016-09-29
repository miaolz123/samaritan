import React from 'react';
import { Tag, Tooltip, Badge, Button, Table, Modal, Form, Input, Select, Popconfirm, notification } from 'antd';
import axios from 'axios';

import config from '../config';
import Logs from './Logs';

const FormItem = Form.Item;
const Option = Select.Option;

class Traders extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
      fetchTradersUrl: '/trader',
      loading: false,
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0,
      },
      tableData: [],
      info: {},
      infoModal: false,
      strategies: [],
      exchanges: [],
      selectedExchanges: [],
      showLogs: false,
    };

    this.handleRefresh = this.handleRefresh.bind(this);
    this.fetchTraders = this.fetchTraders.bind(this);
    this.postTrader = this.postTrader.bind(this);
    this.deleteTrader = this.deleteTrader.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.getStrategyAndExchange = this.getStrategyAndExchange.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoAddShow = this.handleInfoAddShow.bind(this);
    this.handleStrategyChange = this.handleStrategyChange.bind(this);
    this.handleExchangeChange = this.handleExchangeChange.bind(this);
    this.handleExchangeClose = this.handleExchangeClose.bind(this);
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
    this.handleTraderAction = this.handleTraderAction.bind(this);
    this.handleGoBack = this.handleGoBack.bind(this);
  }

  componentWillUpdate(nextProps, nextState) {
    const { showLogs } = this.state;

    if (showLogs && !nextState.showLogs) {
      this.handleRefresh();
    }
  }

  componentWillMount() {
    this.fetchTraders(config.api + this.state.fetchTradersUrl);
  }

  handleRefresh() {
    this.fetchTraders(config.api + this.state.fetchTradersUrl);
  }

  fetchTraders(url) {
    this.setState({ loading: true });

    axios.get(url, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        this.setState({ loading: false });
        if (response.data.success) {
          const { data } = response.data;

          this.setState({
            pagination: { total: data.length },
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

  postTrader(trader) {
    axios.post(`${config.api}/trader`, trader, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.setState({
            infoModal: false,
            selectedExchanges: [],
          });
          this.props.form.resetFields();
          this.fetchTraders(config.api + this.state.fetchTradersUrl);
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
        }
      });
  }

  deleteTrader(id) {
    axios.delete(`${config.api}/trader?id=${id}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.fetchTraders(config.api + this.state.fetchTradersUrl);
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
        }
      });
  }

  handleTableChange(pagination, filters, sorter) {
    let url = '/trader?';
    const sorterMap = {
      'CreatedAt': 'created_at',
      'UpdatedAt': 'updated_at',
    };

    if (sorter && sorter.field) {
      url += `order=${sorterMap[sorter.field]}`;
      if (sorter.order === 'descend') {
        url += ' DESC';
      }
    }

    this.setState({
      fetchTradersUrl: url,
      pagination: pagination,
    });
    this.fetchTraders(config.api + url);
  }

  getStrategyAndExchange(id) {
    axios.get(`${config.api}/strategy?id=${id}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          const { data } = response.data;

          if (data) {
            this.setState({ strategies: data });
          }
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
        }
      });

    axios.get(`${config.api}/exchange?id=${id}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          const { data } = response.data;

          if (data) {
            this.setState({ exchanges: data });
          }
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
        }
      });
  }

  handleInfoShow(info) {
    if (info) {
      this.getStrategyAndExchange(info.ID);
      this.setState({
        info,
        selectedExchanges: info.Exchanges.map(e => e),
        infoModal: true,
      });
    }
  }

  handleInfoAddShow() {
    this.getStrategyAndExchange(0);
    this.setState({
      info: {
        ID: 0,
        Name: '',
        StrategyID: '',
      },
      infoModal: true,
    });
  }

  handleStrategyChange(value) {
    const { getFieldValue, setFieldsValue } = this.props.form;
    const { strategies } = this.state;

    if (getFieldValue('Name') === '') {
      strategies.map(s => {
        if (String(s.ID) === value) {
          setFieldsValue({ Name: `${s.Name}@${new Date().toLocaleDateString()}` });
        }
      });
    }
  }

  handleExchangeChange(value) {
    const { selectedExchanges, exchanges } = this.state;

    if (exchanges[value] && exchanges[value].ID > 0) {
      selectedExchanges.push(exchanges[value]);
      this.setState({ selectedExchanges });
    }
  }

  handleExchangeClose(i, event) {
    const { selectedExchanges } = this.state;

    if (i < selectedExchanges.length) {
      selectedExchanges.splice(i, 1);
      this.setState({ selectedExchanges });
    }
    event.preventDefault();
  }

  handleInfoOk() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const trader = {
        ID: this.state.info.ID,
        Name: values.Name,
        StrategyID: parseInt(values.Strategy, 0),
        Exchanges: this.state.selectedExchanges,
      };

      this.postTrader(trader);
    });
  }

  handleInfoCancel() {
    this.setState({
      infoModal: false,
      strategies: [],
      exchanges: [],
      selectedExchanges: [],
    });
    this.props.form.resetFields();
  }

  handleTraderAction(action, info) {
    const descriptionMap = {
      run: 'Run the trader success',
      stop: 'Stop the trader success',
    };

    axios.post(`${config.api}/${action}`, info, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          notification['success']({
            message: 'Success',
            description: descriptionMap[action],
            duration: 2,
          });
          if (action === 'run') {
            this.setState({ showLogs: true, info });
          }
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
        this.handleRefresh();
      }, (response) => {
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        }
      });
  }

  handleGoBack() {
    this.setState({ showLogs: false });
  }

  render() {
    const { info, tableData, strategies, exchanges, selectedExchanges, showLogs } = this.state;

    if (showLogs) {
      return <Logs trader={info} goBack={this.handleGoBack} />;
    }

    const { getFieldDecorator } = this.props.form;
    const columns = [{
      title: 'Name',
      dataIndex: 'Name',
      render: (text, record) => <a onClick={this.handleInfoShow.bind(this, record)}>{text}</a>,
    }, {
      title: 'Strategy',
      dataIndex: 'Strategy',
      render: text => text.Name,
    }, {
      title: 'CreatedAt',
      dataIndex: 'CreatedAt',
      render: text => text.substr(0, 19),
      sorter: true,
    }, {
      title: 'UpdatedAt',
      dataIndex: 'UpdatedAt',
      render: text => text.substr(0, 19),
      sorter: true,
    }, {
      title: 'Status',
      dataIndex: 'Status',
      render: status => <Badge status={status > 0 ? 'processing' : 'error'} text={status > 0 ? 'RUNNING' : 'HALTAD'}/>,
    }, {
      title: 'Action',
      render: (text, record) => (<Button.Group>
        {record.Status > 0
        ? <Button
            icon="pause-circle-o"
            title="Stop"
            onClick={this.handleTraderAction.bind(this, 'stop', record)}
          />
        : <Popconfirm
            title="Are you sure to RUN it ?"
            onConfirm={this.handleTraderAction.bind(this, 'run', record)}
          >
            <Button icon="play-circle-o" title="Run" />
          </Popconfirm>}
        <Button
          icon="message"
          title="Logs"
          onClick={() => this.setState({ showLogs: true, info: record })}
        />
        <Popconfirm
          title="Are you sure to DELETE it ?"
          onConfirm={this.deleteTrader.bind(this, record.ID)}
        >
          <Button icon="delete" title="Delete" />
        </Popconfirm>
      </Button.Group>),
    }];
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };

    return (
      <div>
        <div style={{ marginBottom: 16, textAlign: 'right' }}>
          <Button style={{ marginRight: 5 }} type="primary" onClick={this.handleInfoAddShow}>Add</Button>
          <Button style={{ marginRight: 10 }} onClick={this.handleRefresh}>Refresh</Button>
          <Tag>Total: {this.state.pagination.total}</Tag>
        </div>
        <Table
          columns={columns}
          dataSource={tableData}
          pagination={this.state.pagination}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
        <Modal closable
          maskClosable={false}
          width="50%"
          title={info.Name || 'New Trader'}
          visible={this.state.infoModal}
          onOk={this.handleInfoOk}
          onCancel={this.handleInfoCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Name"
            >
              {getFieldDecorator('Name', {
                rules: [{ required: true }],
                initialValue: info.Name,
              })(
                <Input />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Strategy"
            >
              {getFieldDecorator('Strategy', {
                rules: [{ required: true }],
                initialValue: info.StrategyID > 0 ? String(info.StrategyID) : '',
              })(
                <Select
                  onSelect={this.handleStrategyChange}
                  notFoundContent="Please add a strategy at first"
                  >
                  {strategies.map(s => <Option key={String(s.ID)} value={String(s.ID)}>{s.Name}</Option>)}
                </Select>
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Exchanges"
            >
              <Select
                onSelect={this.handleExchangeChange}
                notFoundContent="Please add an exchange at first"
                >
                {exchanges.map((e, i) => <Option key={String(i)} value={String(i)}>{e.Name}</Option>)}
              </Select>
              {selectedExchanges.length > 0 ? <div style={{ marginTop: 8 }}>
                {selectedExchanges.map((e, i) => <Tooltip
                  key={String(i)}
                  title={`${i > 0 ? '' : 'Exchange / '}Exchanges[${i}]`}>
                  <Tag closable
                    color={i > 0 ? '' : 'blue'}
                    style={{ marginRight: 5 }}
                    onClose={this.handleExchangeClose.bind(this, i)}>
                    {e.Name}
                  </Tag>
                </Tooltip>)}
              </div> : ''}
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Traders);
