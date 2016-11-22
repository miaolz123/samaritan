import { ResetError } from '../actions';
import { AlgorithmList, AlgorithmCache, AlgorithmPut, AlgorithmDelete } from '../actions/algorithm';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Button, Table, Modal, Form, Input, notification } from 'antd';
import map from 'lodash/map';

const FormItem = Form.Item;

class Algorithm extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      messageErrorKey: '',
      selectedRowKeys: [],
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0,
      },
      info: {},
      infoModalShow: false,
    };

    this.reload = this.reload.bind(this);
    this.onSelectChange = this.onSelectChange.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleInfoShow = this.handleInfoShow.bind(this);
    this.handleInfoSubmit = this.handleInfoSubmit.bind(this);
    this.handleInfoCancel = this.handleInfoCancel.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageErrorKey, pagination } = this.state;
    const { algorithm } = nextProps;

    if (!messageErrorKey && algorithm.message) {
      this.setState({
        messageErrorKey: 'algorithmError',
      });
      notification['error']({
        key: 'algorithmError',
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
    pagination.total = algorithm.total;
    this.setState({ pagination });
  }

  componentWillMount() {
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { pagination } = this.state;
    const { dispatch } = this.props;

    dispatch(AlgorithmList(pagination.pageSize, pagination.current));
  }

  onSelectChange(selectedRowKeys) {
    this.setState({ selectedRowKeys });
  }

  handleTableChange(newPagination, filters, sorter) {
    const { pagination } = this.state;

    pagination.current = newPagination.current;
    this.setState({ pagination });
    this.reload();
  }

  handleDelete() {
    const { dispatch, algorithm } = this.props;
    const { selectedRowKeys, pagination } = this.state;

    if (selectedRowKeys.length > 0) {
      dispatch(AlgorithmDelete(
        map(selectedRowKeys, (i) => algorithm.list[i].id),
        pagination.pageSize,
        pagination.current
      ));
      this.setState({ selectedRowKeys: [] });
    }
  }

  handleInfoShow(info) {
    const { dispatch } = this.props;

    if (!info.id) {
      info = {
        id: 0,
        name: '',
        description: '',
        script: '',
      };
    }
    dispatch(AlgorithmCache(info));
    browserHistory.push('/algorithm/edit');
    // this.setState({ info, infoModalShow: true });
  }

  handleInfoSubmit() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const { dispatch } = this.props;
      const { info, pagination } = this.state;
      const req = {
        id: info.id,
        name: values.name,
        description: values.description,
        script: values.script,
      };

      dispatch(AlgorithmPut(req, pagination.pageSize, pagination.current));
      this.setState({ infoModalShow: false });
      this.props.form.resetFields();
    });
  }

  handleInfoCancel() {
    this.setState({ infoModalShow: false });
    this.props.form.resetFields();
  }

  render() {
    const { selectedRowKeys, pagination, info, infoModalShow } = this.state;
    const { algorithm } = this.props;
    const { getFieldDecorator } = this.props.form;
    const columns = [{
      title: 'Name',
      dataIndex: 'name',
      render: (v, r) => <a onClick={this.handleInfoShow.bind(this, r)}>{String(v)}</a>,
    }, {
      title: 'Description',
      dataIndex: 'description',
      render: (v) => v.substr(0, 36),
    }, {
      title: 'CreatedAt',
      dataIndex: 'createdAt',
      render: (v) => v.toLocaleString(),
    }, {
      title: 'UpdatedAt',
      dataIndex: 'updatedAt',
      render: (v) => v.toLocaleString(),
    }];
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onSelectChange,
    };

    return (
      <div>
        <div className="table-operations">
          <Button.Group>
            <Button type="primary" onClick={this.reload}>Reload</Button>
            <Button onClick={this.handleInfoShow}>Add</Button>
            <Button disabled={selectedRowKeys.length <= 0} onClick={this.handleDelete}>Delete</Button>
          </Button.Group>
        </div>
        <Table columns={columns}
          dataSource={algorithm.list}
          rowSelection={rowSelection}
          pagination={pagination}
          loading={algorithm.loading}
          onChange={this.handleTableChange}
        />
        <Modal closable
          maskClosable={false}
          width="50%"
          title={info.name ? `Algorithm - ${info.name}` : 'New Algorithm'}
          visible={infoModalShow}
          footer=""
          onCancel={this.handleInfoCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Name"
            >
              {getFieldDecorator('name', {
                rules: [{ required: true }],
                initialValue: info.name,
              })(
                <Input />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Description"
            >
              {getFieldDecorator('description', {
                rules: [{ required: true }],
                initialValue: info.description,
              })(
                <Input />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Script"
            >
              {getFieldDecorator('script', {
                rules: [{ required: true }],
                initialValue: info.script,
              })(
                <Input />
              )}
            </FormItem>
            <Form.Item wrapperCol={{ span: 12, offset: 7 }} style={{ marginTop: 24 }}>
              <Button type="primary" onClick={this.handleInfoSubmit} loading={algorithm.loading}>Submit</Button>
            </Form.Item>
          </Form>
        </Modal>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
  algorithm: state.algorithm,
});

export default Form.create()(connect(mapStateToProps)(Algorithm));
