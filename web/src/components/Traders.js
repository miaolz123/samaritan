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
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
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
      if (info.Exchanges) {
        info.ExchangeIDs = info.Exchanges.map(e => String(e.ID));
      } else {
        info.ExchangeIDs = [];
      }

      this.setState({
        info: info,
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
        ExchangeIDs: [],
      },
      infoModal: true,
    });
  }

  handleExchangeChange(value) {
    console.log(202202, value);
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

  render() {
    const { info, tableData, strategies, exchanges } = this.state;
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
                initialValue: info.StrategyID,
              })}>
                {strategies.map(s => <Option key={String(s.ID)} value={String(s.ID)}>{s.Name}</Option>)}
              </Select>
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Exchanges"
            >
              <Select onChange={this.handleExchangeChange}>
                {exchanges.map(e => <Option key={String(e.ID)} value={String(e.ID)}>{e.Name}</Option>)}
              </Select>
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Traders);
