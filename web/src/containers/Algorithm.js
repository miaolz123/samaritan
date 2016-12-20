import { ResetError } from '../actions';
import { AlgorithmList, AlgorithmCache, AlgorithmDelete } from '../actions/algorithm';
import { ExchangeList } from '../actions/exchange';
import { TraderList, TraderPut, TraderDelete, TraderSwitch, TraderCache } from '../actions/trader';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Badge, Button, Dropdown, Form, Input, Menu, Modal, Select, Table, Tag, Tooltip, notification } from 'antd';

const FormItem = Form.Item;
const Option = Select.Option;

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
      traderModelShow: false,
      traderInfo: {
        exchanges: [],
      },
    };

    this.reload = this.reload.bind(this);
    this.onSelectChange = this.onSelectChange.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleTableExpand = this.handleTableExpand.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleEdit = this.handleEdit.bind(this);
    this.handleTraderEdit = this.handleTraderEdit.bind(this);
    this.handleTraderDelete = this.handleTraderDelete.bind(this);
    this.handleTraderSwitch = this.handleTraderSwitch.bind(this);
    this.handleTraderLog = this.handleTraderLog.bind(this);
    this.handleExchangeChange = this.handleExchangeChange.bind(this);
    this.handleExchangeClose = this.handleExchangeClose.bind(this);
    this.handleTraderModelOk = this.handleTraderModelOk.bind(this);
    this.handleTraderModelCancel = this.handleTraderModelCancel.bind(this);
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
    this.order = 'id';
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { pagination } = this.state;
    const { dispatch } = this.props;
    dispatch(AlgorithmList(pagination.pageSize, pagination.current, this.order));
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

  handleTableExpand(expanded, algorithm) {
    if (expanded) {
      const { dispatch } = this.props;

      dispatch(TraderList(algorithm.id));
    }
  }

  handleDelete() {
    Modal.confirm({
      title: 'Are you sure to delete ?',
      onOk: () => {
        const { dispatch } = this.props;
        const { selectedRowKeys, pagination } = this.state;

        if (selectedRowKeys.length > 0) {
          dispatch(AlgorithmDelete(
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

  handleEdit(info) {
    const { dispatch } = this.props;

    if (!info.id) {
      info = {
        id: 0,
        name: 'New Algorithm Name',
        description: '',
        script: '',
      };
    }
    dispatch(AlgorithmCache(info));
    browserHistory.push('/algorithm/edit');
  }

  handleTraderEdit(info, algorithm) {
    const { dispatch } = this.props;

    if (!info) {
      info = {
        id: 0,
        algorithmId: algorithm.id,
        name: `New Trader @ ${new Date().toLocaleDateString()}`,
        exchanges: [],
      };
    }

    this.setState({
      traderModelShow: true,
      traderInfo: info,
    });

    dispatch(ExchangeList(-1, 1, 'id'));
  }

  handleTraderDelete(req) {
    Modal.confirm({
      title: 'Are you sure to delete ?',
      onOk: () => {
        const { dispatch } = this.props;

        dispatch(TraderDelete(req));
      },
      iconType: 'exclamation-circle',
    });
  }

  handleTraderSwitch(req) {
    const { dispatch } = this.props;

    dispatch(TraderSwitch(req));
  }

  handleTraderLog(info) {
    const { dispatch } = this.props;

    dispatch(TraderCache(info));
    browserHistory.push('/algorithm/log');
  }

  handleExchangeChange(value) {
    const { exchange } = this.props;
    const { traderInfo } = this.state;

    if (exchange.list[value] && exchange.list[value].id > 0) {
      traderInfo.exchanges.push(exchange.list[value]);
      this.setState({ traderInfo });
    }
  }

  handleExchangeClose(i, event) {
    const { traderInfo } = this.state;

    if (i < traderInfo.exchanges.length) {
      traderInfo.exchanges.splice(i, 1);
      this.setState({ traderInfo });
    }
    event.preventDefault();
  }

  handleTraderModelOk() {
    this.props.form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const { traderInfo } = this.state;
      const { dispatch } = this.props;
      const info = {
        id: traderInfo.id,
        algorithmId: traderInfo.algorithmId,
        name: values.name,
        exchanges: traderInfo.exchanges,
      };

      dispatch(TraderPut(info));

      this.setState({
        traderModelShow: false,
        traderInfo: {
          exchanges: [],
        },
      });
    });
  }

  handleTraderModelCancel() {
    this.setState({
      traderModelShow: false,
      traderInfo: {
        exchanges: [],
      },
    });
    this.props.form.resetFields();
  }

  render() {
    const { getFieldDecorator } = this.props.form;
    const { selectedRowKeys, pagination, traderModelShow, traderInfo } = this.state;
    const { exchange, algorithm, trader } = this.props;
    const columns = [{
      title: 'Name',
      dataIndex: 'name',
      sorter: true,
      render: (v, r) => <a onClick={this.handleEdit.bind(this, r)}>{v}</a>,
    }, {
      title: 'Description',
      dataIndex: 'description',
      render: (v) => v.substr(0, 36),
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
    }, {
      title: 'Action',
      key: 'action',
      render: (v, r) => (
        <Button onClick={this.handleTraderEdit.bind(this, null, r)} type="ghost">Deploy</Button>
      ),
    }];
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onSelectChange,
    };
    const expcolumns = [{
      title: 'Name',
      dataIndex: 'name',
      render: (v, r) => <a onClick={this.handleTraderEdit.bind(this, r, null)}>{v}</a>,
    }, {
      title: 'Status',
      dataIndex: 'status',
      render: (v) => (v > 0 ? <Badge status="processing" text="RUN" /> : <Badge status="success" text="HALT" />),
    }, {
      title: 'CreatedAt',
      dataIndex: 'createdAt',
      render: (v) => v.toLocaleDateString(),
    }, {
      title: 'UpdatedAt',
      dataIndex: 'updatedAt',
      render: (v) => v.toLocaleDateString(),
    }, {
      title: 'Action',
      key: 'action',
      render: (v, r) => (
        <Dropdown.Button type="ghost" onClick={this.handleTraderSwitch.bind(this, r)} overlay={
          <Menu>
            <Menu.Item key="log">
              <a type="ghost" onClick={this.handleTraderLog.bind(this, r)}>View Log</a>
            </Menu.Item>
            <Menu.Item key="delete">
              <a type="ghost" onClick={this.handleTraderDelete.bind(this, r)}>Delete It</a>
            </Menu.Item>
          </Menu>
        }>{r.status > 0 ? 'Stop' : 'Run'}</Dropdown.Button>
      ),
    }];
    const expandedRowRender = (r) => {
      const data = trader.map[r.id];

      if (data && data.length > 0) {
        return (
          <Table className="womende" rowKey="id"
            size="middle"
            pagination={false}
            columns={expcolumns}
            loading={trader.loading}
            dataSource={trader.map[r.id]}
          />
        );
      }

      if (!trader.loading) {
        return (
          <p>
            No Trader under this algorithm, <a onClick={this.handleTraderEdit.bind(this, null, r)}>deploy</a> one ?
          </p>
        );
      }
    };
    const formItemLayout = {
      labelCol: { span: 7 },
      wrapperCol: { span: 12 },
    };

    return (
      <div>
        <div className="table-operations">
          <Button type="primary" onClick={this.reload}>Reload</Button>
          <Button type="ghost" onClick={this.handleEdit}>Add</Button>
          <Button disabled={selectedRowKeys.length <= 0} onClick={this.handleDelete}>Delete</Button>
        </div>
        <Table rowKey="id"
          columns={columns}
          expandedRowRender={expandedRowRender}
          dataSource={algorithm.list}
          rowSelection={rowSelection}
          pagination={pagination}
          loading={algorithm.loading}
          onChange={this.handleTableChange}
          onExpand={this.handleTableExpand}
        />
        <Modal closable
          maskClosable={false}
          width="50%"
          title={traderInfo.id > 0 ? `Trader - ${traderInfo.name}` : 'New Trader'}
          visible={traderModelShow}
          onOk={this.handleTraderModelOk}
          onCancel={this.handleTraderModelCancel}
        >
          <Form horizontal>
            <FormItem
              {...formItemLayout}
              label="Name"
            >
              {getFieldDecorator('name', {
                rules: [{ required: true }],
                initialValue: traderInfo.name,
              })(
                <Input />
              )}
            </FormItem>
            <FormItem
              {...formItemLayout}
              label="Exchanges"
            >
              <Select
                onSelect={this.handleExchangeChange}
                notFoundContent="Please add an exchange at first"
              >
                {exchange.list.map((e, i) => <Option key={String(i)} value={String(i)}>{e.name}</Option>)}
              </Select>
              {traderInfo.exchanges.length > 0 ? <div style={{ marginTop: 8 }}>
                {traderInfo.exchanges.map((e, i) => <Tooltip
                  key={String(i)}
                  title={`${i > 0 ? '' : 'E / Exchange / '}Es[${i}] / Exchanges[${i}]`}>
                  <Tag closable
                    color={i > 0 ? '' : '#108ee9'}
                    style={{ marginRight: 5 }}
                    onClose={this.handleExchangeClose.bind(this, i)}>
                    {e.name}
                  </Tag>
                </Tooltip>)}
              </div> : ''}
            </FormItem>
          </Form>
        </Modal>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
  exchange: state.exchange,
  algorithm: state.algorithm,
  trader: state.trader,
});

export default Form.create()(connect(mapStateToProps)(Algorithm));
