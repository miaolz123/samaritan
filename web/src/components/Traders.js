import React from 'react';
import { Tag, Button, Table, Modal, Form, Input, Select, notification } from 'antd';
import axios from 'axios';

import config from '../config';

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
    };

    this.handleRefresh = this.handleRefresh.bind(this);
    this.fetchTraders = this.fetchTraders.bind(this);
    this.postTrader = this.postTrader.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoAddShow = this.handleInfoAddShow.bind(this);
    this.handleExchangeChange = this.handleExchangeChange.bind(this);
    this.handleExchangeClose = this.handleExchangeClose.bind(this);
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
    this.handleTraderAction = this.handleTraderAction.bind(this);
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
        if (response.data.success) {
          const { data } = response.data;

          this.setState({
            loading: false,
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
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        }
      });
  }

  postTrader(trader) {
    axios.post(`${config.api}/trader`, trader, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.setState({ infoModal: false });
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

  handleInfoShow(info) {
    axios.get(`${config.api}/strategy`, { headers: { Authorization: `Bearer ${this.state.token}` } })
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

    axios.get(`${config.api}/exchange`, { headers: { Authorization: `Bearer ${this.state.token}` } })
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

    if (info) {
      const selectedExchanges = info.Exchanges.map(e => e);

      this.setState({
        info,
        selectedExchanges,
        infoModal: true,
      });
    }
  }

  handleInfoAddShow() {
    this.setState({
      info: {
        ID: 0,
        Name: '',
        StrategyID: '',
      },
      infoModal: true,
    });
  }

  handleExchangeChange(value) {
    const { selectedExchanges, exchanges } = this.state;

    if (exchanges[value]) {
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
        StrategyID: parseInt(values.StrategyID, 0),
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

  handleTraderAction(action, id) {
    axios.post(`${config.api}/${action}`, {ID: id}, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          notification['success']({
            message: 'Success',
            description: `${action} the trader success`,
            duration: 2,
          });
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

  render() {
    const { info, tableData, strategies, exchanges, selectedExchanges } = this.state;
    const { getFieldProps } = this.props.form;
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
      title: 'Action',
      dataIndex: 'Status',
      render: (status, record) => status > 0
      ? <Button type="" size="small" onClick={this.handleTraderAction.bind(this, 'stop', record.ID)}>Stop</Button>
      : <Button type="" size="small" onClick={this.handleTraderAction.bind(this, 'run', record.ID)}>Run</Button>,
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
        <Table columns={columns}
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
              <Input {...getFieldProps('Name', {
                rules: [{ required: true }],
                initialValue: info.Name,
              })} />
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Strategy"
            >
              <Select {...getFieldProps('StrategyID', {
                rules: [{ required: true }],
                initialValue: String(info.StrategyID),
              })}>
                {strategies.map(s => <Option key={String(s.ID)} value={String(s.ID)}>{s.Name}</Option>)}
              </Select>
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Exchanges"
            >
              <Select onSelect={this.handleExchangeChange}
                {...getFieldProps('Exchange', {
                  rules: [{ required: true }],
                  initialValue: selectedExchanges ? '0' : '',
                })}>
                {exchanges.map((e, i) => <Option key={String(i)} value={String(i)}>{e.Name}</Option>)}
              </Select>
              <div style={{ marginTop: 8 }}>
                {selectedExchanges.map((e, i) => <Tag closable
                  color={i > 0 ? '' : 'blue'}
                  key={String(i)}
                  style={{ marginRight: 5 }}
                  onClose={this.handleExchangeClose.bind(this, i)}
                >{e.Name}</Tag>)}
              </div>
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Traders);
