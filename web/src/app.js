import './styles/app.css';

import React from 'react';
import { render } from 'react-dom';
import { LocaleProvider, Menu, Icon, Modal, Form, Input } from 'antd';
import enUS from 'antd/lib/locale-provider/en_US';
import axios from 'axios';

import config from './config';
import Home from './components/Home';
import Users from './components/Users';
import Exchanges from './components/Exchanges';

const SubMenu = Menu.SubMenu;
const FormItem = Form.Item;

class Example extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      current: 'home',
      loginModalShow: false,
    };

    this.handleClick = this.handleClick.bind(this);
    this.renderMain = this.renderMain.bind(this);
    this.handleLoginOk = this.handleLoginOk.bind(this);
    this.reLogin = this.reLogin.bind(this);
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
          localStorage.setItem('token', response.data.Token);
          this.setState({ loginModalShow: false });
          window.location.href = window.location.href;
        }, () => {});
    });
  }

  reLogin() {
    this.setState({ loginModalShow: true });
  }

  render() {
    const { getFieldProps } = this.props.form;
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };

    return (
      <LocaleProvider locale={enUS}>
        <div className="ant-layout-aside">
          <aside className="ant-layout-sider">
            <div
              className="ant-layout-logo"
              style={{
                background: "url('http://120.27.103.15:8922/images/logo.png')",
                backgroundSize: '64px 30px',
                backgroundRepeat: 'no-repeat',
                backgroundPosition: 'center',
              }}
            >
            </div>
            <Menu theme="dark"
              onClick={this.handleClick}
              defaultOpenKeys={['home']}
              selectedKeys={[this.state.current]}
              mode="inline"
            >
              <Menu.Item key="home"><span><Icon type="pie-chart" /><span>Home</span></span></Menu.Item>
              <SubMenu key="manage" title={<span><Icon type="appstore" /><span>Manage</span></span>}>
                <Menu.Item key="users">Users</Menu.Item>
                <Menu.Item key="exchanges">Exchanges</Menu.Item>
              </SubMenu>
            </Menu>
          </aside>
          <div className="ant-layout-main">
            <Menu theme="dark"
              onClick={this.handleClick}
              selectedKeys={[this.state.current]}
              mode="horizontal"
            >
              <Menu.Item key="logout" style={{ float: 'right' }}><Icon type="logout" />Logout</Menu.Item>
            </Menu>
            <div className="ant-layout-container">
              <div className="ant-layout-content">
                {this.renderMain()}
              </div>
            </div>
            <div className="ant-layout-footer">
            Samaritan © 2016
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
