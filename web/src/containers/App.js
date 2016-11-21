import '../styles/app.less';
import '../styles/app.css';
import { UserGet, Logout } from '../actions/user';
import { ExchangeTypes } from '../actions/exchange';
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

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { user } = nextProps;

    if (user.status < 0) {
      dispatch(Logout());
      browserHistory.push('/login');
    }
  }

  componentWillMount() {
    const { dispatch } = this.props;

    dispatch(UserGet());
    dispatch(ExchangeTypes());
  }

  handleClick(e) {
    const { dispatch } = this.props;

    this.setState({
      current: e.key,
    });

    switch (e.key) {
      case 'trader':
        browserHistory.push('/');
        break;
      case 'user':
        browserHistory.push('/user');
        break;
      case 'exchange':
        browserHistory.push('/exchange');
        break;
      case 'logout':
        dispatch(Logout());
        browserHistory.push('/login');
        break;
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
              defaultOpenKeys={['trader']}
              selectedKeys={[current]}
              mode="inline"
            >
              <Menu.Item key="trader">
                <Icon type="appstore-o" /><span className="nav-text">Trader</span>
              </Menu.Item>
              <Menu.Item key="strategy">
                <Icon type="copy" /><span className="nav-text">Strategy</span>
              </Menu.Item>
              <Menu.Item key="exchange">
                <Icon type="solution" /><span className="nav-text">Exchange</span>
              </Menu.Item>
              <Menu.Item key="user">
                <Icon type="team" /><span className="nav-text">User</span>
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
