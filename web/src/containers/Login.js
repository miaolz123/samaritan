import { ResetError } from '../actions';
import { UserLogin } from '../actions/user';
import React from 'react';
import { connect } from 'react-redux';
import { Button, Form, Input, Icon, Tooltip, notification } from 'antd';

class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      windowHeight: window.innerHeight || 720,
      messageShowed: false,
    };

    this.handleOk = this.handleOk.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageShowed } = this.state;
    const { user } = nextProps;

    if (!messageShowed && user.message) {
      this.setState({ messageShowed: true });
      notification['error']({
        message: 'Error',
        description: String(user.message),
        onClose: () => {
          this.setState({ messageShowed: false });
          dispatch(ResetError());
        },
      });
    }
  }

  handleOk(e) {
    const { form, dispatch } = this.props;

    if (e) {
      e.preventDefault();
    }

    form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      return dispatch(UserLogin(values.cluster, values.username, values.password));
    });
  }

  render() {
    const { windowHeight } = this.state;
    const { user } = this.props;
    const { getFieldDecorator } = this.props.form;
    const formItemLayout = {
      wrapperCol: { offset: 9, span: 6 },
    };
    const cluster = localStorage.getItem('cluster') || 'http://127.0.0.1:9876';

    return (
      <div style={{ paddingTop: windowHeight > 600 ? (windowHeight - 500) / 2 : windowHeight > 400 ? (windowHeight - 350) / 2 : 25 }}>
        <h1 style={{
          margin: 24,
          fontSize: '30px',
          textAlign: 'center',
        }}>Samaritan</h1>
        <Form horizontal onSubmit={this.handleOk}>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Cluster Path">
              {getFieldDecorator('cluster', {
                rules: [{ required: true }],
                initialValue: cluster,
              })(
                <Input addonBefore={<Icon type="appstore-o" />} placeholder="http://127.0.0.1:9876" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Username">
              {getFieldDecorator('username', {
                rules: [{ required: true }],
              })(
                <Input addonBefore={<Icon type="user" />} placeholder="username" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Password">
              {getFieldDecorator('password', {
                rules: [{ required: true }],
              })(
                <Input addonBefore={<Icon type="lock" />} type="password" placeholder="password" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item wrapperCol={{ span: 15, offset: 9 }} style={{ marginTop: 24 }}>
            <Button type="primary" htmlType="submit" loading={user.loading}>Login</Button>
          </Form.Item>
        </Form>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
});

export default Form.create()(connect(mapStateToProps)(Login));
