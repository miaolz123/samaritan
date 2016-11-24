import { ResetError } from '../actions';
import { UserList, UserPut, UserDelete } from '../actions/user';
import React from 'react';
import { connect } from 'react-redux';
import { Button, Table, Modal, Form, Input, InputNumber, notification } from 'antd';

const FormItem = Form.Item;

class User extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      messageErrorKey: '',
      selectedRowKeys: [],
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0,
      },
      info: {},
      infoModalShow: false,
    };

    this.reload = this.reload.bind(this);
    this.onSelectChange = this.onSelectChange.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoSubmit = this.handleInfoSubmit.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageErrorKey, pagination } = this.state;
    const { user } = nextProps;

    if (!messageErrorKey && user.message) {
      this.setState({
        messageErrorKey: 'userError',
      });
      notification['error']({
        key: 'userError',
        message: 'Error',
        description: String(user.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        },
      });
    }
    pagination.total = user.total;
    this.setState({ pagination });
  }

  componentWillMount() {
    this.order = 'id';
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { pagination } = this.state;
    const { dispatch } = this.props;

    dispatch(UserList(pagination.pageSize, pagination.current, this.order));
  }

  onSelectChange(selectedRowKeys) {
    this.setState({ selectedRowKeys });
  }

  handleTableChange(newPagination, filters, sorter) {
    const { pagination } = this.state;

    if (sorter.field) {
      this.order = `${sorter.field} ${sorter.order.replace('end', '')}`;
    } else {
      this.order = 'id';
    }
    pagination.current = newPagination.current;
    this.setState({ pagination });
    this.reload();
  }

  handleDelete() {
    Modal.confirm({
      title: 'Are you sure to delete ?',
      onOk: () => {
        const { dispatch } = this.props;
        const { selectedRowKeys, pagination } = this.state;

        if (selectedRowKeys.length > 0) {
          dispatch(UserDelete(
            selectedRowKeys,
            pagination.pageSize,
            pagination.current,
            this.order
          ));
          this.setState({ selectedRowKeys: [] });
        }
      },
      iconType: 'exclamation-circle',
    });
  }

  handleInfoShow(info) {
    if (!info.username) {
      const { user } = this.props;
      info = {
        id: 0,
        username: '',
        level: user.data ? user.data.level - 1 : 0,
      };
    }
    this.setState({ info, infoModalShow: true });
  }

  handleInfoSubmit() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const { dispatch } = this.props;
      const { info, pagination } = this.state;
      const req = {
        id: info.id,
        username: values.username,
        level: values.level,
      };

      dispatch(UserPut(req, values.password, pagination.pageSize, pagination.current, this.order));
      this.setState({ infoModalShow: false });
      this.props.form.resetFields();
    });
  }

  handleInfoCancel() {
    this.setState({ infoModalShow: false });
    this.props.form.resetFields();
  }

  render() {
    const { selectedRowKeys, pagination, info, infoModalShow } = this.state;
    const { user } = this.props;
    const { getFieldDecorator, getFieldValue } = this.props.form;
    const columns = [{
      title: 'Username',
      dataIndex: 'username',
      sorter: true,
      render: (v, r) => <a onClick={this.handleInfoShow.bind(this, r)}>{String(v)}</a>,
    }, {
      title: 'Level',
      dataIndex: 'level',
      sorter: true,
    }, {
      title: 'CreatedAt',
      dataIndex: 'createdAt',
      sorter: true,
      render: (v) => v.toLocaleString(),
    }, {
      title: 'UpdatedAt',
      dataIndex: 'updatedAt',
      sorter: true,
      render: (v) => v.toLocaleString(),
    }];
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };
    const checkPassword = (rule, value, callback) => {
      if (value && value !== getFieldValue('password')) {
        callback('Confirm fail');
      } else {
        callback();
      }
    };
    const passwdProps = info.id ? getFieldDecorator('password', {
      rules: [{ required: false, whitespace: true }],
    }) : getFieldDecorator('password', {
      rules: [{ required: true, whitespace: true }],
    });
    const repasswdProps = info.id && !getFieldValue('password') ? getFieldDecorator('rePassword', {
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
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onSelectChange,
    };

    return (
      <div>
        <div className="table-operations">
          <Button type="primary" onClick={this.reload}>Reload</Button>
          <Button type="ghost" onClick={this.handleInfoShow}>Add</Button>
          <Button disabled={selectedRowKeys.length <= 0} onClick={this.handleDelete}>Delete</Button>
        </div>
        <Table rowKey="id"
          columns={columns}
          dataSource={user.list}
          rowSelection={rowSelection}
          pagination={pagination}
          loading={user.loading}
          onChange={this.handleTableChange}
        />
        <Modal closable
          maskClosable={false}
          width="50%"
          title={info.username ? `User - ${info.username}` : 'New User'}
          visible={infoModalShow}
          footer=""
          onCancel={this.handleInfoCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Username"
            >
              {getFieldDecorator('username', {
                rules: [{ required: true }],
                initialValue: info.username,
              })(
                <Input disabled={info.id > 0} />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Level"
            >
              {getFieldDecorator('level', {
                rules: [{ required: true }],
                initialValue: info.level,
              })(
                <InputNumber disabled={user.data && info.username === user.data.username} min={0} max={user.data ? user.data.level : 0} />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Password"
            >
              {passwdProps(
                <Input type="password" />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Repeat"
            >
              {repasswdProps(
                <Input type="Password" />
              )}
            </FormItem>
            <Form.Item wrapperCol={{ span: 12, offset: 7 }} style={{ marginTop: 24 }}>
              <Button type="primary" onClick={this.handleInfoSubmit} loading={user.loading}>Submit</Button>
            </Form.Item>
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
