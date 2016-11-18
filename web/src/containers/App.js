import '../styles/app.less';
import '../styles/app.css';
import React, { Component } from 'react';

export default class App extends Component {
  render() {
    const { children } = this.props;
    return (
      <div>
        {children}
      </div>
    );
  }
}
