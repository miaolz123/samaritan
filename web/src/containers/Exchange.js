import { ResetError } from '../actions';
import { ExchangeList, ExchangePut, ExchangeDelete } from '../actions/exchange';
import React from 'react';
import { connect } from 'react-redux';
import { Button, Table, Modal, Form, Input, Select, notification } from 'antd';

const FormItem = Form.Item;
const Option = Select.Option;

class Exchange extends React.Component {
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
    const { exchange } = nextProps;

    if (!messageErrorKey && exchange.message) {
      this.setState({
        messageErrorKey: 'exchangeError',
      });
      notification['error']({
        key: 'exchangeError',
        message: 'Error',
        description: String(exchange.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        },
      });
    }
    pagination.total = exchange.total;
    this.setState({ pagination });
  }

  componentWillMount() {
    this.order = 'id';
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { pagination } = this.state;
    const { dispatch } = this.props;

    dispatch(ExchangeList(pagination.pageSize, pagination.current, this.order));
  }

  onSelectChange(selectedRowKeys) {
    this.setState({ selectedRowKeys });
  }

  handleTableChange(newPagination, filters, sorter) {
    const { pagination } = this.state;

    if (sorter.field) {
      this.order = `${sorter.field} ${sorter.order.replace('end', '')}`;
    } else {
      this.order = 'id';
    }
    pagination.current = newPagination.current;
    this.setState({ pagination });
    this.reload();
  }

  handleDelete() {
    Modal.confirm({
      title: 'Are you sure to delete ?',
      onOk: () => {
        const { dispatch } = this.props;
        const { selectedRowKeys, pagination } = this.state;

        if (selectedRowKeys.length > 0) {
          dispatch(ExchangeDelete(
            selectedRowKeys,
            pagination.pageSize,
            pagination.current,
            this.order
          ));
          this.setState({ selectedRowKeys: [] });
        }
      },
      iconType: 'exclamation-circle',
    });
  }

  handleInfoShow(info) {
    if (!info.id) {
      info = {
        id: 0,
        name: '',
        type: '',
        accessKey: '',
        secretKey: '',
      };
    }
    this.setState({ info, infoModalShow: true });
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
        type: values.type,
        accessKey: values.accessKey,
        secretKey: values.secretKey,
      };

      dispatch(ExchangePut(req, pagination.pageSize, pagination.current, this.order));
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
    const { exchange } = this.props;
    const { getFieldDecorator } = this.props.form;
    const columns = [{
      title: 'Name',
      dataIndex: 'name',
      sorter: true,
      render: (v, r) => <a onClick={this.handleInfoShow.bind(this, r)}>{String(v)}</a>,
    }, {
      title: 'Type',
      dataIndex: 'type',
      sorter: true,
    }, {
      title: 'CreatedAt',
      dataIndex: 'createdAt',
      sorter: true,
      render: (v) => v.toLocaleString(),
    }, {
      title: 'UpdatedAt',
      dataIndex: 'updatedAt',
      sorter: true,
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
          <Button type="primary" onClick={this.reload}>Reload</Button>
          <Button type="ghost" onClick={this.handleInfoShow}>Add</Button>
          <Button disabled={selectedRowKeys.length <= 0} onClick={this.handleDelete}>Delete</Button>
        </div>
        <Table rowKey="id"
          columns={columns}
          dataSource={exchange.list}
          rowSelection={rowSelection}
          pagination={pagination}
          loading={exchange.loading}
          onChange={this.handleTableChange}
        />
        <Modal closable
          maskClosable={false}
          width="50%"
          title={info.name ? `Exchange - ${info.name}` : 'New Exchange'}
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
              label="Type"
            >
              {getFieldDecorator('type', {
                rules: [{ required: true }],
                initialValue: info.type,
              })(
                <Select disabled={info.id > 0}>
                  {exchange.types.map((v, i) => <Option key={i} value={v}>{v}</Option>)}
                </Select>
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="AccessKey"
            >
              {getFieldDecorator('accessKey', {
                rules: [{ required: true }],
                initialValue: info.accessKey,
              })(
                <Input />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="SecretKey"
            >
              {getFieldDecorator('secretKey', {
                rules: [{ required: true }],
                initialValue: info.secretKey,
              })(
                <Input />
              )}
            </FormItem>
            <Form.Item wrapperCol={{ span: 12, offset: 7 }} style={{ marginTop: 24 }}>
              <Button type="primary" onClick={this.handleInfoSubmit} loading={exchange.loading}>Submit</Button>
            </Form.Item>
          </Form>
        </Modal>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
  exchange: state.exchange,
});

export default Form.create()(connect(mapStateToProps)(Exchange));
