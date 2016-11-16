import './styles/app.less';
import './styles/app.css';

import React from 'react';
import { render } from 'react-dom';
import { LocaleProvider, Menu, Icon, Modal } from 'antd';
import enUS from 'antd/lib/locale-provider/en_US';
import axios from 'axios';

import config from './config';
import Home from './pages/Home';
import Users from './pages/Users';
import Exchanges from './pages/Exchanges';
import Strategies from './pages/Strategies';
import Traders from './pages/Traders';
import Login from './pages/Login';

class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
      loading: true,
      collapse: true,
      current: 'traders',
      loginShow: false,
    };

    this.handleClick = this.handleClick.bind(this);
    this.renderMain = this.renderMain.bind(this);
    this.reLogin = this.reLogin.bind(this);
    this.onCollapseChange = this.onCollapseChange.bind(this);
  }

  componentWillMount() {
    const client = hprose.Client.create('http://127.0.0.1:9888', ['hello']);
    client.hello('hahahah', (resp) => {
      console.log(383838, resp);
    });
    axios.post(`${config.api}/token`, null, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        this.setState({ loading: false });
      }, (response) => {
        this.setState({ loginShow: true });
      });
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
        this.setState({ loginShow: true });
        Modal.confirm({
          title: 'Confirm',
          content: 'Are you sure to logout ?',
          onOk: () => {
            localStorage.removeItem('token');
            window.location.href = window.location.href;
          },
          onCancel: () => {
            window.location.href = window.location.href;
          },
        });
        return '';
      default:
        return <div style={{ height: '100%' }}>ERROR!</div>;
    }
  }

  reLogin() {
    this.setState({ loginShow: true });
  }

  onCollapseChange() {
    this.setState({
      collapse: !this.state.collapse,
    });
  }

  render() {
    const { loading, collapse, loginShow } = this.state;

    if (loginShow) {
      return (<LocaleProvider locale={enUS}>
          <Login />
        </LocaleProvider>);
    } else if (loading) {
      return null;
    }

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
                <Icon type="poweroff" /><span className="nav-text">Logout</span>
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
        </div>
      </LocaleProvider>
    );
  }
}

render(<App />, document.getElementById('react-app'));
