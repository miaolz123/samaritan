import './styles/app.css';

import React from 'react';
import { render } from 'react-dom';
import { LocaleProvider, Menu, Icon, Modal, Form, Input, notification } from 'antd';
import enUS from 'antd/lib/locale-provider/en_US';
import axios from 'axios';

import config from './config';
import Home from './pages/Home';
import Users from './pages/Users';
import Exchanges from './pages/Exchanges';
import Strategies from './pages/Strategies';
import Traders from './pages/Traders';

const FormItem = Form.Item;

class Example extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      collapse: false,
      current: 'traders',
      loginModalShow: false,
    };

    this.handleClick = this.handleClick.bind(this);
    this.renderMain = this.renderMain.bind(this);
    this.handleLoginOk = this.handleLoginOk.bind(this);
    this.reLogin = this.reLogin.bind(this);
    this.onCollapseChange = this.onCollapseChange.bind(this);
  }

  handleClick(e) {
    this.setState({
      current: e.key,
    });
  }

  renderMain() {
    switch (this.state.current) {
      case 'home':
        return <Home style={{ height: '100%' }} reLogin={this.reLogin} />;
      case 'users':
        return <Users style={{ height: '100%' }} reLogin={this.reLogin} />;
      case 'exchanges':
        return <Exchanges style={{ height: '100%' }} reLogin={this.reLogin} />;
      case 'strategies':
        return <Strategies style={{ height: '100%' }} reLogin={this.reLogin} />;
      case 'traders':
        return <Traders style={{ height: '100%' }} reLogin={this.reLogin} />;
      case 'logout':
        localStorage.removeItem('token');
        window.location.href = window.location.href;
        return '';
      default:
        return <div style={{ height: '100%' }}>ERROR!</div>;
    }
  }

  handleLoginOk() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      axios.post(`${config.api}/login`, values)
        .then((response) => {
          if (response.data.success) {
            localStorage.setItem('token', response.data.data);
            this.setState({ loginModalShow: false });
            window.location.href = window.location.href;
          } else {
            notification['error']({
              message: 'Error',
              description: String(response.data.msg),
              duration: null,
            });
          }
        }, () => {});
    });
  }

  reLogin() {
    this.setState({ loginModalShow: true });
  }

  onCollapseChange() {
    this.setState({
      collapse: !this.state.collapse,
    });
  }

  render() {
    const { collapse } = this.state;
    const { getFieldProps } = this.props.form;
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };

    return (
      <LocaleProvider locale={enUS}>
        <div className={collapse ? 'ant-layout-aside ant-layout-aside-collapse' : 'ant-layout-aside'}>
          <aside className="ant-layout-sider">
            {collapse ? '' : <div className="ant-layout-logo"></div>}
            <Menu theme="dark"
              onClick={this.handleClick}
              defaultOpenKeys={['traders']}
              selectedKeys={[this.state.current]}
              mode="inline"
            >
              <Menu.Item key="traders">
                <Icon type="appstore-o" /><span className="nav-text">Traders</span>
              </Menu.Item>
              <Menu.Item key="strategies">
                <Icon type="copy" /><span className="nav-text">Strategies</span>
              </Menu.Item>
              <Menu.Item key="exchanges">
                <Icon type="solution" /><span className="nav-text">Exchanges</span>
              </Menu.Item>
              <Menu.Item key="users">
                <Icon type="team" /><span className="nav-text">Users</span>
              </Menu.Item>
              <Menu.Item key="logout">
                <Icon type="logout" /><span className="nav-text">Logout</span>
              </Menu.Item>
            </Menu>
            <div className="ant-aside-action" onClick={this.onCollapseChange}>
              {collapse ? <Icon type="right" /> : <Icon type="left" />}
            </div>
          </aside>
          <div className="ant-layout-main">
            <div className="ant-layout-container">
              <div className="ant-layout-content">
                {this.renderMain()}
              </div>
            </div>
            <div className="ant-layout-footer">
              <a href="https://github.com/miaolz123/samaritan">Samaritan</a> © 2016
            </div>
          </div>
          <Modal
            maskClosable={false}
            width="50%"
            title='Login'
            visible={this.state.loginModalShow}
            onOk={this.handleLoginOk}
          >
            <Form horizontal>
              <FormItem
                {...formItemLayout}
                label="Username"
              >
                <Input {...getFieldProps('name', {
                  rules: [{ required: true }],
                })} />
              </FormItem>
              <FormItem
                {...formItemLayout}
                label="Password"
              >
                <Input type="password" {...getFieldProps('password', {
                  rules: [{ required: true }],
                })} />
              </FormItem>
            </Form>
          </Modal>
        </div>
      </LocaleProvider>
    );
  }
}

const App = Form.create()(Example);

render(<App />, document.getElementById('react-app'));
