import '../styles/app.less';
import '../styles/app.css';
import { UserGet, Logout } from '../actions/user';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { LocaleProvider, Menu, Icon } from 'antd';
import enUS from 'antd/lib/locale-provider/en_US';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      collapse: true,
      current: 'traders',
    };

    this.handleClick = this.handleClick.bind(this);
    this.onCollapseChange = this.onCollapseChange.bind(this);
  }

  componentWillMount() {
    const { dispatch } = this.props;

    dispatch(UserGet());
  }

  componentWillReceiveProps(nextProps) {
    const { user } = nextProps;

    if (!user.loading && user.message) {
      browserHistory.push('/login');
    }
  }

  handleClick(e) {
    this.setState({
      current: e.key,
    });

    if (e.key === 'logout') {
      const { dispatch } = this.props;

      dispatch(Logout());
      browserHistory.push('/login');
    }
  }

  onCollapseChange() {
    this.setState({
      collapse: !this.state.collapse,
    });
  }

  render() {
    const { collapse, current } = this.state;
    const { children } = this.props;

    return (
      <LocaleProvider locale={enUS}>
        <div className={collapse ? 'ant-layout-aside ant-layout-aside-collapse' : 'ant-layout-aside'}>
          <aside className="ant-layout-sider">
            {collapse ? '' : <div className="ant-layout-logo"></div>}
            <Menu theme="dark"
              onClick={this.handleClick}
              defaultOpenKeys={['traders']}
              selectedKeys={[current]}
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
                {children}
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

const mapStateToProps = (state) => ({
  user: state.user,
});

export default connect(mapStateToProps)(App);
