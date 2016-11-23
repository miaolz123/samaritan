import { ResetError } from '../actions';
import { AlgorithmList, AlgorithmCache, AlgorithmDelete } from '../actions/algorithm';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Button, Table, Modal, notification } from 'antd';
import map from 'lodash/map';

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
    };

    this.reload = this.reload.bind(this);
    this.onSelectChange = this.onSelectChange.bind(this);
    this.handleTableChange = this.handleTableChange.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleEdit = this.handleEdit.bind(this);
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

  handleDelete() {
    Modal.confirm({
      title: 'Are you sure to delete ?',
      onOk: () => {
        const { dispatch, algorithm } = this.props;
        const { selectedRowKeys, pagination } = this.state;

        if (selectedRowKeys.length > 0) {
          dispatch(AlgorithmDelete(
            map(selectedRowKeys, (i) => algorithm.list[i].id),
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

  render() {
    const { selectedRowKeys, pagination } = this.state;
    const { algorithm } = this.props;
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
    }];
    const rowSelection = {
      selectedRowKeys,
      onChange: this.onSelectChange,
    };

    return (
      <div>
        <div className="table-operations">
          <Button.Group>
            <Button type="primary" onClick={this.reload}>Reload</Button>
            <Button onClick={this.handleEdit}>Add</Button>
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
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
  algorithm: state.algorithm,
});

export default connect(mapStateToProps)(Algorithm);
