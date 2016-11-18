import React, { Component } from 'react';
import { connect } from 'react-redux';

class Home extends Component {
  render() {
    return (
      <div className="container">
          home!
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
});

export default connect(mapStateToProps)(Home);
