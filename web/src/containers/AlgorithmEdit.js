import { ResetError } from '../actions';
import { AlgorithmPut } from '../actions/algorithm';
import React, { Component } from 'react';
import { findDOMNode } from 'react-dom';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Row, Col, Tooltip, Input, Button, notification } from 'antd';
import CodeMirror from 'codemirror';
require('codemirror/lib/codemirror.css');
require('codemirror/mode/javascript/javascript.js');
require('codemirror/theme/eclipse.css');
require('codemirror/addon/edit/matchbrackets.js');
require('codemirror/addon/fold/foldcode.js');
require('codemirror/addon/fold/foldgutter.js');
require('codemirror/addon/fold/foldgutter.css');
require('codemirror/addon/fold/brace-fold.js');
require('codemirror/addon/fold/comment-fold.js');
require('codemirror/addon/lint/lint.js');
require('codemirror/addon/lint/lint.css');
require('codemirror/addon/lint/javascript-lint.js');
require('codemirror/addon/selection/active-line.js');
require('codemirror/addon/scroll/simplescrollbars.js');
require('codemirror/addon/scroll/simplescrollbars.css');

class AlgorithmEdit extends Component {
  constructor(props) {
    super(props);

    this.state = {
      innerHeight: window.innerHeight > 500 ? window.innerHeight : 500,
      messageErrorKey: '',
      name: '',
      description: '',
      script: '',
    };

    this.handleNameChange = this.handleNameChange.bind(this);
    this.handleDescriptionChange = this.handleDescriptionChange.bind(this);
    this.handleScriptChange = this.handleScriptChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { messageErrorKey } = this.state;
    const { algorithm } = nextProps;

    if (!algorithm.cache.name) {
      browserHistory.push('/algorithm');
    }

    if (!messageErrorKey && algorithm.message) {
      this.setState({
        messageErrorKey: 'algorithmEditError',
      });
      notification['error']({
        key: 'algorithmEditError',
        message: 'Error',
        description: String(algorithm.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        },
      });
    }
  }

  componentWillMount() {
    const { name } = this.state;
    const { algorithm } = this.props;

    if (!name) {
      this.setState({
        name: algorithm.cache.name,
        description: algorithm.cache.description,
        script: algorithm.cache.script,
      });
    }
  }

  componentDidMount() {
    const textareaNode = findDOMNode(this.refs.script);

    this.codeMirror = CodeMirror.fromTextArea(textareaNode, {
      theme: 'eclipse',
      tabSize: 2,
      lineWrapping: true,
      lineNumbers: true,
      matchBrackets: true,
      foldGutter: true,
      lint: true,
      styleActiveLine: true,
      scrollbarStyle: 'simple',
      gutters: ['CodeMirror-lint-markers', 'CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
    });
    this.codeMirror.on('change', this.handleScriptChange);
  }

  componentWillUnmount() {
    notification.destroy();
  }

  handleNameChange(e) {
    this.setState({ name: e.target.value });
  }

  handleDescriptionChange(e) {
    this.setState({ description: e.target.value });
  }

  handleScriptChange(doc, change) {
    if (change.origin !== 'setValue') {
      this.setState({ script: doc.getValue() });
    }
  }

  handleSubmit() {
    const { dispatch, algorithm } = this.props;
    const { name, description, script } = this.state;
    const req = {
      id: algorithm.cache.id,
      name,
      description,
      script,
    };

    dispatch(AlgorithmPut(req));
  }

  handleCancel() {
    browserHistory.goBack();
  }

  render() {
    const { innerHeight, name, description, script } = this.state;

    return (
      <div className="container">
        <Row type="flex" justify="space-between">
          <Col span={18}>
            <Tooltip placement="bottomLeft" title="Algorithm Name">
              <Input
                placeholder="Algorithm Name"
                defaultValue={name}
                onChange={this.handleNameChange}
              />
            </Tooltip>
          </Col>
          <Col span={6} style={{textAlign: 'right'}}>
            <Button.Group>
              <Button
                type="primary"
                disabled={!name}
                onClick={this.handleSubmit}
              >Submit</Button>
              <Button
                type="ghost"
                onClick={this.handleCancel}
              >Cancel</Button>
            </Button.Group>
          </Col>
        </Row>
        <Row style={{marginTop: 18}}>
          <Tooltip placement="bottomLeft" title="Algorithm Description">
            <Input
              rows={1}
              type="textarea"
              placeholder="Algorithm Description"
              defaultValue={description}
              onChange={this.handleDescriptionChange}
            />
          </Tooltip>
        </Row>
        <Row style={{height: innerHeight - 190, marginTop: 18}}>
          <textarea ref="script" defaultValue={script} />
        </Row>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  algorithm: state.algorithm,
});

export default connect(mapStateToProps)(AlgorithmEdit);
