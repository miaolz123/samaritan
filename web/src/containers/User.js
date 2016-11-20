import { UserList } from '../actions/user';
import React from 'react';
import { connect } from 'react-redux';
import { Tag, Button, Table, Modal, Form, Input, InputNumber } from 'antd';

const FormItem = Form.Item;

class User extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
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

  componentWillReceiveProps(nextProps) {
    const { pagination } = this.state;
    const { user } = nextProps;

    if (user.total > 0) {
      pagination.total = user.total;
      this.setState({
        pagination,
        tableData: user.list,
      });
    }
  }

  componentWillMount() {
    this.fetchUsers();
  }

  handleRefresh() {
    this.fetchUsers();
  }

  fetchUsers() {
    const { pagination } = this.state;
    const { dispatch } = this.props;

    dispatch(UserList(pagination.pageSize, pagination.current));
  }

  postUser(user) {
  }

  deleteUser(id) {
  }

  handleTableChange(pagination, filters, sorter) {
    this.setState({ pagination: pagination });
    this.fetchUsers();
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
      dataIndex: 'name',
      render: (text, record) => <a onClick={this.handleInfoShow.bind(this, record)}>{text}</a>,
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

const mapStateToProps = (state) => ({
  user: state.user,
});

export default Form.create()(connect(mapStateToProps)(User));
