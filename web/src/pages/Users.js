import React from 'react';
import { Tag, Button, Table, Modal, Form, Input, InputNumber, Popconfirm, notification } from 'antd';
import axios from 'axios';

import config from '../config';

const FormItem = Form.Item;

class Users extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
      fetchUsersUrl: '/user',
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
    this.fetchUsers = this.fetchUsers.bind(this);
    this.postUser = this.postUser.bind(this);
    this.deleteUser = this.deleteUser.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoAddShow = this.handleInfoAddShow.bind(this);
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
  }

  componentWillMount() {
    this.fetchUsers(config.api + this.state.fetchUsersUrl);
  }

  handleRefresh() {
    this.fetchUsers(config.api + this.state.fetchUsersUrl);
  }

  fetchUsers(url) {
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

  postUser(user) {
    axios.post(`${config.api}/user`, user, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.setState({ infoModal: false });
          this.props.form.resetFields();
          this.fetchUsers(config.api + this.state.fetchUsersUrl);
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

  deleteUser(id) {
    axios.delete(`${config.api}/user?id=${id}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.fetchUsers(config.api + this.state.fetchUsersUrl);
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
    let url = '/user?';
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
      fetchUsersUrl: url,
      pagination: pagination,
    });
    this.fetchUsers(config.api + url);
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
        Level: 0,
      },
      infoModal: true,
    });
  }

  handleInfoOk() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const user = {
        ID: this.state.info.ID,
        Name: values.Name,
        Password: values.Password,
        Level: values.Level,
      };

      this.postUser(user);
    });
  }

  handleInfoCancel() {
    this.setState({ infoModal: false });
    this.props.form.resetFields();
  }

  render() {
    const { info, tableData } = this.state;
    const { getFieldDecorator, getFieldValue } = this.props.form;
    const columns = [{
      title: 'Username',
      dataIndex: 'Name',
      render: (text, record) => <a onClick={this.handleInfoShow.bind(this, record)}>{text}</a>,
    }, {
      title: 'Level',
      dataIndex: 'Level',
      sorter: true,
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
      key: 'action',
      render: (text, record) => <Popconfirm title="Are you sure to delete it ?" onConfirm={this.deleteUser.bind(this, record.ID)}>
        <Button
          icon="delete"
          title="Delete"
        />
      </Popconfirm>,
    }];
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };
    const checkPassword = (rule, value, callback) => {
      if (value && value !== getFieldValue('Password')) {
        callback('Confirm fail');
      } else {
        callback();
      }
    };
    const passwdProps = info.ID ? getFieldDecorator('Password', {
      rules: [{ required: false, whitespace: true }],
    }) : getFieldDecorator('Password', {
      rules: [{ required: true, whitespace: true }],
    });
    const repasswdProps = info.ID && !getFieldValue('Password') ? getFieldDecorator('rePassword', {
      rules: [
        { required: false, whitespace: true },
        { validator: checkPassword },
      ],
    }) : getFieldDecorator('rePassword', {
      rules: [
        { required: true, whitespace: true },
        { validator: checkPassword },
      ],
    });

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
          title={info.Name || 'New User'}
          visible={this.state.infoModal}
          onOk={this.handleInfoOk}
          onCancel={this.handleInfoCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Username"
            >
              {getFieldDecorator('Name', {
                rules: [{ required: true }],
                initialValue: info.Name,
              })(
                <Input disabled={info.ID > 0} />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Level"
            >
              {getFieldDecorator('Level', {
                rules: [{ required: true }],
                initialValue: info.Level,
              })(
                <InputNumber min={0} max={tableData.length > 0 ? tableData[0].Level : 99} />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Password"
            >
              {passwdProps(
                <Input type="Password" />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Password Confirm"
            >
              {repasswdProps(
                <Input type="Password" />
              )}
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Users);
