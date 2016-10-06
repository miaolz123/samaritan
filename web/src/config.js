module.exports = {
  api: location.href.replace(/(\w+)\s*\//, '$1'),
  // api: 'http://127.0.0.1:9876',
  exchangeTypes: ['okcoin.cn', 'huobi', 'poloniex', 'global'],
  logTypes: {
    '-1': 'ERROR',
    '0': 'INFO',
    '1': 'PROFIT',
    '2': 'BUY',
    '3': 'SELL',
    '4': 'CANCEL',
  },
};
