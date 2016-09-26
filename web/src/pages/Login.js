import React from 'react';
import { Button, Form, Input } from 'antd';
import axios from 'axios';

import config from '../config';

class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: false,
    };

    this.handleOk = this.handleOk.bind(this);
  }

  handleOk() {
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

  render() {
    const { getFieldProps } = this.props.form;

    return (
      <div>
        <Form horizontal>
          <Form.Item
            {...Form.ItemLayout}
            label="Username"
          >
            <Input
              disabled={info.ID > 0}
              {...getFieldProps('Name', {
                rules: [{ required: true }],
                initialValue: info.Name,
              })} />
          </Form.Item>
          <Form.Item
            {...Form.ItemLayout}
            label="Level"
          >
            <InputNumber
            max={tableData.length > 0 ? tableData[0].Level : 99}
            {...getFieldProps('Level', {
              rules: [{ required: true }],
              initialValue: info.Level,
            })} />
          </Form.Item>
          <Form.Item
            {...Form.ItemLayout}
            label="Password"
          >
            <Input type="Password" {...passwdProps} />
          </Form.Item>
          <Form.Item
            {...Form.ItemLayout}
            label="Password Confirm"
          >
            <Input type="Password" {...repasswdProps} />
          </Form.Item>
        </Form>
      </div>
    );
  }
}

export default Form.create()(Login);
