import { ResetError } from '../actions';
import { UserLogin } from '../actions/user';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Button, Form, Input, Icon, Tooltip, notification } from 'antd';

class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      windowHeight: window.innerHeight || 720,
      messageErrorKey: '',
    };

    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageErrorKey } = this.state;
    const { user } = nextProps;

    if (!messageErrorKey && user.message) {
      this.setState({
        messageErrorKey: 'userLoginError',
      });
      notification['error']({
        key: 'userLoginError',
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

    if (!user.loading && user.token) {
      if (this.state.messageErrorKey) {
        this.setState({ messageErrorKey: '' });
        notification.close(this.state.messageErrorKey);
      }
      browserHistory.push('/');
    }
  }

  componentWillUnmount() {
    const { dispatch } = this.props;

    dispatch(ResetError());
  }

  handleSubmit(e) {
    const { form, dispatch } = this.props;

    if (e) {
      e.preventDefault();
    }

    form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      dispatch(UserLogin(values.cluster, values.username, values.password));
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
        <Form horizontal onSubmit={this.handleSubmit}>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Cluster URL">
              {getFieldDecorator('cluster', {
                rules: [{ type: 'url', required: true }],
                initialValue: cluster,
              })(
                <Input addonBefore={<Icon type="link" />} placeholder="http://127.0.0.1:9876" />
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
          <Form.Item wrapperCol={{ span: 6, offset: 9 }} style={{ marginTop: 24 }}>
            <Button type="primary" htmlType="submit" className="login-form-button">Login</Button>
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
