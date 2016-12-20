import { ResetError } from '../actions';
import { LogList } from '../actions/log';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Button, Table, Tag, notification } from 'antd';

class Log extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      messageErrorKey: '',
      pagination: {
        pageSize: 12,
        current: 1,
        total: 0,
      },
      filters: {},
    };

    this.reload = this.reload.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageErrorKey, pagination } = this.state;
    const { trader, log } = nextProps;

    if (!trader.cache.name) {
      browserHistory.push('/algorithm');
    }

    if (!messageErrorKey && log.message) {
      this.setState({
        messageErrorKey: 'logError',
      });
      notification['error']({
        key: 'logError',
        message: 'Error',
        description: String(log.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        },
      });
    }
    pagination.total = log.total;
    this.setState({ pagination });
  }

  componentWillMount() {
    this.filters = {};
    this.reload();
  }

  componentWillUnmount() {
    notification.destroy();
  }

  reload() {
    const { pagination } = this.state;
    const { trader, dispatch } = this.props;

    dispatch(LogList(trader.cache, pagination, this.filters));
  }

  handleTableChange(newPagination, filters) {
    const { pagination } = this.state;

    pagination.current = newPagination.current;
    this.filters = filters;
    this.setState({ pagination });
    this.reload();
  }

  handleCancel() {
    browserHistory.push('/algorithm');
  }

  render() {
    const { pagination } = this.state;
    const { log } = this.props;
    const colors = {
      'INFO': '#A9A9A9',
      'ERROR': '#F50F50',
      'PROFIT': '#4682B4',
      'CANCEL': '#5F9EA0',
    };
    const columns = [{
      width: 160,
      title: 'Time',
      dataIndex: 'time',
      render: (v) => v.toLocaleString(),
    }, {
      width: 100,
      title: 'Exchange',
      dataIndex: 'exchangeType',
      render: (v) => <Tag color={v === 'global' ? '' : '#00BFFF'}>{v}</Tag>,
    }, {
      width: 100,
      title: 'Type',
      dataIndex: 'type',
      render: (v) => <Tag color={colors[v] || '#00BFFF'}>{v}</Tag>,
    }, {
      title: 'Price',
      dataIndex: 'price',
      width: 100,
    }, {
      width: 100,
      title: 'Amount',
      dataIndex: 'amount',
    }, {
      title: 'Message',
      dataIndex: 'message',
    }];

    return (
      <div>
        <div className="table-operations">
          <Button type="primary" onClick={this.reload}>Reload</Button>
          <Button type="ghost" onClick={this.handleCancel}>Back</Button>
        </div>
        <Table rowKey="id"
          columns={columns}
          dataSource={log.list}
          pagination={pagination}
          loading={log.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
  trader: state.trader,
  log: state.log,
});

export default connect(mapStateToProps)(Log);
