import React from 'react';
import { Tag, Button, Table, Modal, Form, Input, InputNumber, notification } from 'antd';
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
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoAddShow = this.handleInfoAddShow.bind(this);
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
    this.handleInfoResetPasswd = this.handleInfoResetPasswd.bind(this);
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
        Level: values.Level,
      };

      if (!user.ID) {
        user.Password = values.Password;
      }
      this.postUser(user);
    });
  }

  handleInfoCancel() {
    this.setState({ infoModal: false });
    this.props.form.resetFields();
  }

  handleInfoResetPasswd() {
    Modal.confirm({
      title: 'Confirm password reset ?',
      content: 'Click OK to reset the password and the password will be set to same as the username !',
      iconType: 'question-circle-o',
      onOk: () => {
        const { info } = this.state;
        const user = {
          ID: info.ID,
          Name: info.Name,
          Password: info.Name,
          Level: info.Level,
        };

        this.postUser(user);
      },
    });
  }

  render() {
    const { info, tableData } = this.state;
    const { getFieldProps } = this.props.form;
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
    }];
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };
    const checkPassword = (rule, value, callback) => {
      const { getFieldValue } = this.props.form;
      if (value && value !== getFieldValue('Password')) {
        callback('Confirm fail');
      } else {
        callback();
      }
    };
    const passwdProps = info.ID ? getFieldProps('Password', {
      rules: [{ required: false }],
    }) : getFieldProps('Password', {
      rules: [{ required: true, whitespace: true }],
    });
    const repasswdProps = info.ID ? getFieldProps('rePassword', {
      rules: [{ required: false }],
    }) : getFieldProps('rePassword', {
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
              <Input {...getFieldProps('Name', {
                rules: [{ required: true }],
                initialValue: info.Name,
              })} />
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Level"
            >
              <InputNumber {...getFieldProps('Level', {
                initialValue: info.Level,
              })} />
            </FormItem>
            {info.ID
            ? <FormItem wrapperCol={{ span: 16, offset: 6 }} style={{ marginTop: 24 }}>
                <Button type="dashed" size="default" onClick={this.handleInfoResetPasswd}>Reset Password</Button>
              </FormItem>
            : <div><FormItem
                {...formItemLayout}
                label="Password"
              >
                <Input type="Password" {...passwdProps} />
              </FormItem>
              <FormItem
                {...formItemLayout}
                label="Confirm"
              >
                <Input type="Password" {...repasswdProps} />
              </FormItem></div>
            }
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Users);
