import React, { Component } from 'react';
import { findDOMNode } from 'react-dom';
import { connect } from 'react-redux';
import { Row, Col, Tooltip, Input, Button } from 'antd';
import CodeMirror from 'codemirror';
require('codemirror/mode/javascript/javascript');
require('codemirror/lib/codemirror.css');
require('codemirror/theme/solarized.css');

class AlgorithmEdit extends Component {
  constructor(props) {
    super(props);

    this.state = {
      script: 'var i = 1;',
    };

    this.updateCode = this.updateCode.bind(this);
  }

  componentDidMount() {
    const textareaNode = findDOMNode(this.refs.script);

    this.codeMirror = CodeMirror.fromTextArea(textareaNode, {
      lineNumbers: true,
    });
    this.codeMirror.on('change', this.updateCode);
  }

  updateCode(doc, change) {
    if (change.origin !== 'setValue') {
      console.log(292929, doc.getValue());
    }
  }

  render() {
    return (
      <div className="container">
        <Row type="flex" justify="space-between">
          <Col span={18}>
            <Tooltip placement="bottomLeft" title="Algorithm Name">
              <Input placeholder="Algorithm Name" />
            </Tooltip>
          </Col>
          <Col span={6} style={{textAlign: 'right'}}>
            <Button.Group>
              <Button type="ghost">Cancel</Button>
              <Button type="primary">Submit</Button>
            </Button.Group>
          </Col>
        </Row>
        <Row style={{height: 500, marginTop: 18}}>
          <textarea ref="script" defaultValue={this.state.script} />
        </Row>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
  algorithm: state.algorithm,
});

export default connect(mapStateToProps)(AlgorithmEdit);
