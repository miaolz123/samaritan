import React from 'react';
import { Tag, Button, Table, Modal, Form, Input, Popconfirm, notification } from 'antd';
import axios from 'axios';
import CodeMirror from 'react-code-mirror';
require('codemirror/mode/javascript/javascript');
require('codemirror/lib/codemirror.css');
require('codemirror/theme/solarized.css');

import config from '../config';

const FormItem = Form.Item;

class Strategies extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: localStorage.getItem('token'),
      windowHeight: window.innerHeight || 720,
      fetchStrategiesUrl: '/strategy',
      loading: false,
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0,
      },
      tableData: [],
      info: {},
      script: '',
      infoModal: false,
    };

    this.handleRefresh = this.handleRefresh.bind(this);
    this.fetchStrategies = this.fetchStrategies.bind(this);
    this.postStrategy = this.postStrategy.bind(this);
    this.deleteStrategy = this.deleteStrategy.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleScriptChange = this.handleScriptChange.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoAddShow = this.handleInfoAddShow.bind(this);
    this.handleInfoOk = this.handleInfoOk.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
  }

  componentWillMount() {
    this.fetchStrategies(config.api + this.state.fetchStrategiesUrl);
  }

  handleRefresh() {
    this.fetchStrategies(config.api + this.state.fetchStrategiesUrl);
  }

  fetchStrategies(url) {
    this.setState({ loading: true });

    axios.get(url, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          const { data } = response.data;

          this.setState({
            loading: false,
            pagination: { total: data.length },
            tableData: data,
          });
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
      }, (response) => {
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        }
      });
  }

  postStrategy(strategy) {
    axios.post(`${config.api}/strategy`, strategy, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.setState({ infoModal: false });
          this.props.form.resetFields();
          this.fetchStrategies(config.api + this.state.fetchStrategiesUrl);
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
      }, (response) => {
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        }
      });
  }

  deleteStrategy(id) {
    axios.delete(`${config.api}/strategy?id=${id}`, { headers: { Authorization: `Bearer ${this.state.token}` } })
      .then((response) => {
        if (response.data.success) {
          this.fetchStrategies(config.api + this.state.fetchStrategiesUrl);
        } else {
          notification['error']({
            message: 'Error',
            description: String(response.data.msg),
            duration: null,
          });
        }
      }, (response) => {
        if (String(response).indexOf('401') > 0) {
          this.setState({ token: '' });
          localStorage.removeItem('token');
          this.props.reLogin();
        }
      });
  }

  handleTableChange(pagination, filters, sorter) {
    let url = '/strategy?';
    const sorterMap = {
      'CreatedAt': 'created_at',
      'UpdatedAt': 'updated_at',
    };

    if (sorter && sorter.field) {
      url += `order=${sorterMap[sorter.field]}`;
      if (sorter.order === 'descend') {
        url += ' DESC';
      }
    }

    this.setState({
      fetchStrategiesUrl: url,
      pagination: pagination,
    });
    this.fetchStrategies(config.api + url);
  }

  handleScriptChange(e) {
    this.setState({ script: e.target.value });
  }

  handleInfoShow(info) {
    if (info) {
      this.setState({
        info: info,
        script: info.Script,
        infoModal: true,
      });
    }
  }

  handleInfoAddShow() {
    this.setState({
      info: {
        ID: 0,
        Name: '',
        Description: '',
        Script: '',
      },
      infoModal: true,
    });
  }

  handleInfoOk() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const { info, script } = this.state;
      const strategy = {
        ID: info.ID,
        Name: values.Name,
        Description: values.Description,
        Script: script,
      };

      this.postStrategy(strategy);
    });
  }

  handleInfoCancel() {
    this.setState({ infoModal: false });
    this.props.form.resetFields();
  }

  render() {
    const { windowHeight, info, tableData, script } = this.state;
    const { getFieldProps } = this.props.form;
    const columns = [{
      title: 'Name',
      dataIndex: 'Name',
      render: (text, record) => <a onClick={this.handleInfoShow.bind(this, record)}>{text}</a>,
    }, {
      title: 'Description',
      dataIndex: 'Description',
      render: text => text.substr(0, 36),
    }, {
      title: 'CreatedAt',
      dataIndex: 'CreatedAt',
      render: text => text.substr(0, 19),
      sorter: true,
    }, {
      title: 'UpdatedAt',
      dataIndex: 'UpdatedAt',
      render: text => text.substr(0, 19),
      sorter: true,
    }, {
      title: 'Action',
      key: 'action',
      render: (text, record) => <Popconfirm title="Are you sure to delete it ?" onConfirm={this.deleteStrategy.bind(this, record.ID)}>
        <Button
          icon="delete"
          title="Delete"
        />
      </Popconfirm>,
    }];
    const formItemLayout = {
      labelCol: { span: 3 },
      wrapperCol: { span: 18 },
    };

    return (
      <div>
        <div style={{ marginBottom: 16, textAlign: 'right' }}>
          <Button style={{ marginRight: 5 }} type="primary" onClick={this.handleInfoAddShow}>Add</Button>
          <Button style={{ marginRight: 10 }} onClick={this.handleRefresh}>Refresh</Button>
          <Tag>Total: {this.state.pagination.total}</Tag>
        </div>
        <Table columns={columns}
          dataSource={tableData}
          pagination={this.state.pagination}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
        <Modal closable
          width="85%"
          style={{top: windowHeight > 590 ? (windowHeight - 590) / 2 : 20}}
          maskClosable={false}
          title={info.Name || 'New Strategy'}
          visible={this.state.infoModal}
          onOk={this.handleInfoOk}
          onCancel={this.handleInfoCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Name"
            >
              <Input {...getFieldProps('Name', {
                rules: [{ required: true }],
                initialValue: info.Name,
              })} />
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Description"
            >
              <Input type="textarea"
              {...getFieldProps('Description', {
                initialValue: info.Description,
              })} />
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Script"
            >
              <CodeMirror
                style={{border: '1px solid #d9d9d9'}}
                mode="javascript"
                theme="solarized"
                value={script}
                onChange={this.handleScriptChange}
                lineNumbers={true}
              />
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

export default Form.create()(Strategies);
