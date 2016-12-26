# Samaritan service

## Installation

You can install samaritan from **installation package** or **docker**.

### From installation package

1. Download the samaritan installation package on [this page](https://github.com/miaolz123/samaritan/releases)
2. Unzip the samaritan installation package
3. Enter the extracted samaritan installation directory
4. Run `samaritan`

Then, samaritan is running at `http://localhost:9876`.

**Linux & Mac user quick start command**

```shell
wget https://github.com/miaolz123/samaritan/releases/download/v{{VERSION}}/samaritan_{{OS}}_{{ARCH}}.tar.gz && \
tar -xzvf samaritan_{{OS}}_{{ARCH}}.tar.gz.tar.gz && \
cd samaritan_{{OS}}_{{ARCH}} && \
./samaritan
```

Please replace *{{VERSION}}*, *{{OS}}*, *{{ARCH}}* first.

### From Docker

```shell
docker run --name=samaritan -p 19876:9876 miaolz123/samaritan
```

Then, samaritan is running at `http://localhost:19876`.

## Supported exchanges

| Exchange | Stock |
| -------- | ----- |
| okcoin.cn | `BTC/CNY`, `LTC/CNY` |
| huobi | `BTC/CNY`, `LTC/CNY` |
| poloniex | `ETH/BTC`, `XMR/BTC`, `BTC/USDT`, `LTC/BTC`, `ETC/BTC`, `XRP/BTC`, `ETH/USDT`, `ETC/ETH`, ... |
| btcc | `BTC/CNY`, `LTC/CNY`, `LTC/BTC` |
| chbtc | `BTC/CNY`, `LTC/CNY`, `ETH/CNY`, `ETC/CNY` |
| okcoin.future | `BTC.WEEK/USD`, `BTC.WEEK2/USD`, `BTC.MONTH3/USD`, `LTC.WEEK/USD`, ... |
| oanda.v20 | coming soon ...... |

# Algorithm Reference

## Protocols

### Global constant

| Name | Type | Description |
| ---- | ---- | ----------- |
| Global/G| Object | a object with some global methods |
| Exchange/E | Object | a object with some exchange methods |
| Exchanges/Es | Object List | an `Exchange/E` list |

### Trade type

| Name | Type | Description |
| ---- | ---- | ----------- |
| BUY | String | buy |
| SELL | String | sell |
| LONG | String | long contract |
| SHORT | String | short contract |
| LONG_CLOSE | String | close long contract |
| SHORT_CLOSE | String | close short contract |

### Records period

| Name | Type | Description |
| ---- | ---- | ----------- |
| M | String | 1 minute |
| M5 | String | 5 minutes |
| M15 | String | 15 minutes |
| M30 | String | 30 minutes |
| H | String | 1 hour |
| D | String | 1 day |
| W | String | 1 week |

## Data struct

### Account

| Name | Type | Description |
| ---- | ---- | ----------- |
| Balance | Number | balance amount |
| FrozenBalance | Number | frozen balance amount |
| BTC | Number | BTC amount |
| FrozenBTC | Number | frozen BTC amount |
| LTC | Number | LTC amount |
| FrozenLTC | Number | frozen LTC amount |
| ... | Number | ... amount |
| Frozen... | Number | frozen ... amount |
| Stock | Number | main stock amount |
| FrozenStock | Number | frozen main stock amount |

### Position

| Name | Type | Description |
| ---- | ---- | ----------- |
| Price | Number | price |
| Leverage | Number | leverage |
| Amount | Number | total position amount |
| FrozenAmount | Number | frozen position amount |
| Profit | Number | profit |
| ContractType | String | contract type |
| TradeType | String | trade type |
| StockType | String | stock type |

### Order

| Name | Type | Description |
| ---- | ---- | ----------- |
| ID | String | unique id |
| Price | Number | price |
| Amount | Number | total amount |
| DealAmount | Number | deal amount |
| Fee | Number | fee of this order |
| TradeType | Number | trade type |
| StockType | String | stock type |

### Record

| Name | Type | Description |
| ---- | ---- | ----------- |
| Time | Number | unix timestamp |
| Open | Number | open price |
| High | Number | high price |
| Low | Number | low price |
| Close | Number | close price |
| Volume | Number | trade volume |

### OrderBook

| Name | Type | Description |
| ---- | ---- | ----------- |
| Price | Number | price |
| Amount | Number | market depth amount |

### Ticker

| Name | Type | Description |
| ---- | ---- | ----------- |
| Bids | OrderBook List | bid market depth list |
| Buy | Number | the first bid price, `Bids[0].Price` |
| Mid | Number | `(Buy + Sell) / 2` |
| Sell | Number | the first ask price, `Asks[0].Price` |
| Asks | OrderBook List | ask market depth list |

## Global/G

`Global`/`G` is a object with some global methods.

### Sleep

> G.Sleep(Interval: *Any*) => *No Return*

```javascript
// the program will sleep for 5 seconds
// if Interval <= 0, will automatic execute AutoSleep() of all Exchanges
G.Sleep(5000);
```

### Log

> G.Log(Message: *Any*) => *No Return*

```javascript
// send a message to web control
G.Log("I'm running…");
```

### LogProfit

> G.LogProfit(Profit: *Number*, Message: *Any*) => *No Return*

```javascript
// send a profit message to web control to show profit chart
G.LogProfit(12.345, 'Round 1 end');
```

### LogStatus

> G.LogStatus(Message: *Any*) => *No Return*

```javascript
// send a status message to web control to show it real-time
G.LogStatus('Latest BTC Ticker: ', E.GetTicker('BTC/USD'));
```

### AddTask

> G.AddTask(Function: *Function*, Arguments: *Any*) => *Boolean*

```javascript
// work with G.ExecTasks()
```

### ExecTasks

> G.ExecTasks() => *List*

```javascript
// send same tasks to task list
G.AddTask(E.GetAccount);
G.AddTask(E.GetTicker, 'BTC/USD');

// execute all tasks at the same time and return all results
var results = G.ExecTasks();
var thisAccount = results[0];
var thisTicker = results[1];
```

## Exchange/E

`Exchange`/`E` is a object with some exchange methods.

### Log

> E.Log(Message: *Any*) => *No Return*

```javascript
// send a message of this exchange to web control
E.Log("I'm running…");
```

### GetType

> E.GetType() => *String*

```javascript
// get the type of this exchange
var thisType = E.GetType();
```

### GetName

> E.GetName() => *String*

```javascript
// get the name of this exchange
var thisName = E.GetName();
```

### GetMainStock

> E.GetMainStock() => *String*

```javascript
// get the main stock type of this exchange
var thisMainStock = E.GetMainStock();
```

### SetMainStock

> E.SetMainStock(StockType: *String*) => *String*

```javascript
// set the main stock type of this exchange
var newMainStockType = E.SetMainStock('LTC/USD');
```

### SetLimit

> E.SetLimit(times: *Number*) => *Number*

```javascript
// set the limit calls amount per second of this exchange
// work with E.AutoSleep()
var newLimit = E.SetLimit(6);
```

### AutoSleep

> E.AutoSleep() => *No Return*

```javascript
// auto sleep to achieve the limit calls amount per second of this exchange
E.AutoSleep();
```

### GetAccount

> E.GetAccount() => *Account*

```javascript
// get the account info of this exchange
var thisAccount = E.GetAccount();
```

### GetPositions

> E.GetPositions(StockType: *String*) => *Position List*

```javascript
// get the position list of this exchange
var thisPositions = E.GetPositions('BTC/USD');
```

### GetMinAmount

> E.GetMinAmount(StockType: *String*) => *Number*

```javascript
// get the min trade amount of this exchange
var thisMinAmount = E.GetMinAmount('BTC/USD');
```

### Trade

> E.Trade(TradeType: [*String*](#trade-type), StockType: *String*, Price: *Number*, Amount: *Number*, Message: *Any*) => *String*/*Boolean*

```javascript
// buy example
// if Price <= 0, it's a market order, and the Amount will be different
// return ID of this order if succeed
// return false if fail
E.Trade('BUY', 'BTC/USD', 600, 0.5, 'I paid $300'); // normal order
E.Trade('BUY', 'BTC/USD', 0, 300, 'I also paid $300'); // market order

// sell example
// if Price <= 0, it's a market order
// return ID of this order if succeed
// return false if fail
E.Trade('SELL', 'BTC/USD', 600, 0.5); // normal order
E.Trade('SELL', 'BTC/USD', 0, 0.5); // market order
```

### GetOrder

> E.GetOrder(StockType: *String*, ID: *String*) => *Order*/*Boolean*

```javascript
// return info of this order if succeed
// return false if fail
var thisOrder = E.GetOrder('BTC/USD', 'XXXXXX');
```

### GetOrders

> E.GetOrders(StockType: *String*) => *Order List*

```javascript
// return all the undone orders
var thisOrders = E.GetOrders('BTC/USD');
```

### GetTrades

> E.GetTrades(StockType: *String*) => *Order List*

```javascript
// return all the done orders
var thisTrades = E.GetTrades('BTC/USD');
```

### CancelOrder

> E.CancelOrder(Order: *Order*) => *Boolean*

```javascript
var thisOrders = E.GetOrders('BTC/USD');
for (var i = 0; i < thisOrders.length; i++) {
    // return the result
    var isCanceled = E.CancelOrder(thisOrders[i]);
}
```

### GetTicker

> E.GetTicker(StockType: *String*, Size: *Any*) => *Ticker*

```javascript
// get the latest ticker of this exchange
var thisTicker = E.GetTicker('BTC/USD');
```

### GetRecords

> E.GetRecords(StockType: *String*, Period: [*String*](#records-period), Size: *Any*) => *Record List*

```javascript
// get the latest records of this exchange
var thisRecords = E.GetRecords('BTC/USD', 'M5');
```
