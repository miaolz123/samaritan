import React from 'react';
import { Button, Form, Input, notification } from 'antd';
import axios from 'axios';

import config from '../config';

class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      windowHeight: window.innerHeight || 720,
      loading: false,
    };

    this.handleOk = this.handleOk.bind(this);
  }

  handleOk(e) {
    if (e) {
      e.preventDefault();
    }
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      this.setState({ loading: true });
      axios.post(`${config.api}/login`, values)
        .then((response) => {
          this.setState({ loading: false });
          if (response.data.success) {
            localStorage.setItem('token', response.data.data);
            window.location.href = window.location.href;
          } else {
            notification['error']({
              message: 'Error',
              description: String(response.data.msg),
              duration: null,
            });
          }
        }, (response) => {
          this.setState({ loading: false });
          notification['error']({
            message: 'Error',
            description: String(response),
            duration: null,
          });
        });
    });
  }

  render() {
    const { loading, windowHeight } = this.state;
    const { getFieldDecorator } = this.props.form;
    const formItemLayout = {
      labelCol: { span: 9 },
      wrapperCol: { span: 6 },
    };

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
            label="Username"
          >
            {getFieldDecorator('Name', {
              rules: [{ required: true }],
            })(
              <Input />
            )}
          </Form.Item>
          <Form.Item
            {...formItemLayout}
            label="Password"
          >
            {getFieldDecorator('Password', {
              rules: [{ required: true }],
            })(
              <Input type="password" />
            )}
          </Form.Item>
          <Form.Item wrapperCol={{ span: 15, offset: 9 }} style={{ marginTop: 24 }}>
            <Button type="primary" htmlType="submit" loading={loading}>Submit</Button>
          </Form.Item>
        </Form>
      </div>
    );
  }
}

export default Form.create()(Login);
