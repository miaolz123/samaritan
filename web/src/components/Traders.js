import React from 'react';
import { Tag, Tooltip, Button, Table, Modal, Form, Input, Select, notification } from 'antd';
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
    this.fetchSelectedExchanges = this.fetchSelectedExchanges.bind(this);
    this.postTrader = this.postTrader.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.getStrategyAndExchange = this.getStrategyAndExchange.bind(this);
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

  fetchSelectedExchanges(id) {
    axios.get(`${config.api}/trader/${id}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          const { data } = response.data;

          if (data) {
            this.setState({ selectedExchanges: response.data.data });
          } else {
            this.setState({ selectedExchanges: [] });
          }
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
      }, (response) => {
        this.setState({ selectedExchanges: [] });
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

  getStrategyAndExchange() {
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
  }

  handleInfoShow(info) {
    if (info) {
      this.fetchSelectedExchanges(info.ID);
      this.getStrategyAndExchange();
      this.setState({
        info,
        infoModal: true,
      });
    }
  }

  handleInfoAddShow() {
    this.getStrategyAndExchange();
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
    const { info, selectedExchanges, exchanges } = this.state;

    if (exchanges[value] && exchanges[value].ID > 0) {
      if (info.ID > 0) {
        axios.post(`${config.api}/trader/${info.ID}`, exchanges[value],
        { headers: { Authorization: `Bearer ${this.state.token}` } })
          .then((response) => {
            if (response.data.success) {
              this.fetchSelectedExchanges(info.ID);
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
      } else {
        selectedExchanges.push(exchanges[value]);
        this.setState({ selectedExchanges });
      }
    }
  }

  handleExchangeClose(i, event) {
    const { info, selectedExchanges } = this.state;

    if (i < selectedExchanges.length) {
      if (info.ID > 0) {
        axios.delete(`${config.api}/trader/${info.ID}?id=${selectedExchanges[i].ID}`,
        { headers: { Authorization: `Bearer ${this.state.token}` } })
          .then((response) => {
            if (response.data.success) {
              this.fetchSelectedExchanges(info.ID);
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
      } else {
        selectedExchanges.splice(i, 1);
        this.setState({ selectedExchanges });
      }
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
    const descriptionMap = {
      run: 'Run the trader success',
      stop: 'Stop the trader success',
    };

    axios.post(`${config.api}/${action}`, {ID: id}, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          notification['success']({
            message: 'Success',
            description: descriptionMap[action],
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
              <Select onSelect={this.handleExchangeChange}>
                {exchanges.map((e, i) => <Option key={String(i)} value={String(i)}>{e.Name}</Option>)}
              </Select>
              <div style={{ marginTop: 8 }}>
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
              </div>
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Traders);
