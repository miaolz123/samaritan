import { ResetError } from '../actions';
import { UserList, UserPut } from '../actions/user';
import React from 'react';
import { connect } from 'react-redux';
import { Button, Table, Modal, Form, Input, InputNumber, notification } from 'antd';

const FormItem = Form.Item;

class Users extends React.Component {
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
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { pagination } = this.state;
    const { dispatch } = this.props;

    dispatch(UserList(pagination.pageSize, pagination.current));
  }

  onSelectChange(selectedRowKeys) {
    this.setState({ selectedRowKeys });
  }

  handleTableChange(pagination, filters, sorter) {
    this.setState({ pagination: pagination });
    this.reload();
  }

  handleDelete() {
    const { list } = this.props.user;
    const { selectedRowKeys } = this.state;
    const ids = selectedRowKeys.each((i) => {
      console.log(list[i].id);
      return 111;
    });
    console.log(878787, ids);
  }

  handleInfoShow(info) {
    if (!info.name) {
      const { user } = this.props;
      info = {
        id: 0,
        name: '',
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
      const { info } = this.state;
      const req = {
        id: info.id,
        name: values.name,
        level: values.level,
      };

      dispatch(UserPut(req, values.password));
      this.setState({ infoModalShow: false });
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
      dataIndex: 'name',
      render: (v, r) => <a onClick={this.handleInfoShow.bind(this, r)}>{String(v)}</a>,
    }, {
      title: 'Level',
      dataIndex: 'level',
    }, {
      title: 'CreatedAt',
      dataIndex: 'createdAt',
      render: (v) => v.toLocaleString(),
    }, {
      title: 'UpdatedAt',
      dataIndex: 'updatedAt',
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
          <Button.Group>
            <Button type="primary" onClick={this.reload}>Reload</Button>
            <Button onClick={this.handleInfoShow}>Add</Button>
            <Button disabled={selectedRowKeys.length <= 0} onClick={this.handleDelete}>Delete</Button>
          </Button.Group>
        </div>
        <Table columns={columns}
          dataSource={user.list}
          rowSelection={rowSelection}
          pagination={pagination}
          loading={user.loading}
          onChange={this.handleTableChange}
        />
        <Modal closable
          maskClosable={false}
          width="50%"
          title={`User - ${info.name}` || 'New User'}
          visible={infoModalShow}
          footer=""
          onCancel={this.handleInfoCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Username"
            >
              {getFieldDecorator('name', {
                rules: [{ required: true }],
                initialValue: info.name,
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
                <InputNumber disabled={user.data && info.name === user.data.name} min={0} max={user.data ? user.data.level : 0} />
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
              label="Password Confirm"
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

export default Form.create()(connect(mapStateToProps)(Users));
