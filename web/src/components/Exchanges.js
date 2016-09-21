import React from 'react';
import { Tag, Button, Table, Modal, Form, Input, notification } from 'antd';
import axios from 'axios';

import config from '../config';

const FormItem = Form.Item;

class Exchanges extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
      fetchExchangesUrl: '/exchange',
      loading: false,
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0,
      },
      tableData: [],
      info: {},
      infoModal: false,
    };

    this.handleRefresh = this.handleRefresh.bind(this);
    this.fetchExchanges = this.fetchExchanges.bind(this);
    this.postExchange = this.postExchange.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoAddShow = this.handleInfoAddShow.bind(this);
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
    this.handleInfoResetPasswd = this.handleInfoResetPasswd.bind(this);
  }

  componentWillMount() {
    this.fetchExchanges(config.api + this.state.fetchExchangesUrl);
  }

  handleRefresh() {
    this.fetchExchanges(config.api + this.state.fetchExchangesUrl);
  }

  fetchExchanges(url) {
    this.setState({ loading: true });

    axios.get(url, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          const { data } = response.data;
          console.log(535353, data);

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

  postExchange(exchange) {
    axios.post(`${config.api}/exchange`, exchange, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.setState({ infoModal: false });
          this.props.form.resetFields();
          this.fetchExchanges(config.api + this.state.fetchExchangesUrl);
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
    let url = '/exchange?';
    const sorterMap = {
      'Level': 'level',
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
      fetchExchangesUrl: url,
      pagination: pagination,
    });
    this.fetchExchanges(config.api + url);
  }

  handleInfoShow(info) {
    if (info) {
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
        Type: '',
        AccessKey: '',
        SecretKey: '',
      },
      infoModal: true,
    });
  }

  handleInfoOk() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const exchange = {
        ID: this.state.info.ID,
        Name: values.Name,
        Type: values.Type,
        AccessKey: values.AccessKey,
        SecretKey: values.SecretKey,
      };

      this.postExchange(exchange);
    });
  }

  handleInfoCancel() {
    this.setState({ infoModal: false });
    this.props.form.resetFields();
  }

  handleInfoResetPasswd() {
    Modal.confirm({
      title: 'Confirm password reset ?',
      content: 'Click OK to reset the password and the password will be set to same as the exchangename !',
      iconType: 'question-circle-o',
      onOk: () => {
        const { info } = this.state;
        const exchange = {
          ID: info.ID,
          Name: info.Name,
          Password: info.Name,
          Level: info.Level,
        };

        this.postExchange(exchange);
      },
    });
  }

  render() {
    const { info, tableData } = this.state;
    const { getFieldProps } = this.props.form;
    const columns = [{
      title: 'Name',
      dataIndex: 'Name',
      render: (text, record) => <a onClick={this.handleInfoShow.bind(this, record)}>{text}</a>,
    }, {
      title: 'Type',
      dataIndex: 'Type',
    }, {
      title: 'AccessKey',
      dataIndex: 'AccessKey',
    }, {
      title: 'SecretKey',
      dataIndex: 'SecretKey',
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
          title={info.Name || 'New Exchange'}
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
              label="Type"
            >
              <Input {...getFieldProps('Type', {
                rules: [{ required: true }],
                initialValue: info.Type,
              })} />
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="AccessKey"
            >
              <Input {...getFieldProps('AccessKey', {
                rules: [{ required: true }],
                initialValue: info.AccessKey,
              })} />
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="SecretKey"
            >
              <Input {...getFieldProps('SecretKey', {
                rules: [{ required: true }],
                initialValue: info.SecretKey,
              })} />
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Exchanges);
