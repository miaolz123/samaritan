## Installation

You can install samaritan from **installation package** or **docker**.

### From installation package

1. Download the samaritan installation package on [this page](https://github.com/miaolz123/samaritan/releases)
2. Unzip the samaritan installation package
3. Enter the extracted samaritan installation directory
4. Run `samaritan`

Default, samaritan is running at `http://localhost:9876`.

### From Docker

```shell
docker run --name=samaritan -p 19876:9876 miaolz123/samaritan
```

Then, samaritan is running at `http://localhost:19876`.

# Algorithm Reference

## Protocols

### Global constant

| Name | Type | Description |
| ---- | ---- | ----------- |
| Global/G| Object | a object with some global methods |
| Exchange/E | Object | a object with some exchange methods |
| Exchanges/Es | Object List | an `Exchange/E` list |

### Order type

| Name | Type | Description |
| ---- | ---- | ----------- |
| 1 | Number | buy order |
| -1 | Number | sell order |
| 2 | Number | buy market order |
| -2 | Number | sell market order |

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
| Total | Number | total asset |
| Net | Number | net asset |
| Balance | Number | balance amount |
| FrozenBalance | Number | frozen balance amount |
| BTC | Number | BTC amount |
| FrozenBTC | Number | frozen BTC amount |
| LTC | Number | LTC amount |
| FrozenLTC | Number | frozen LTC amount |
| Stock | Number | main stock amount |
| FrozenStock | Number | frozen main stock amount |

### Order

| Name | Type | Description |
| ---- | ---- | ----------- |
| ID | String | unique id |
| Price | Number | price |
| Amount | Number | total amount |
| DealAmount | Number | deal amount |
| OrderType | Number | [type reference](#order-type) |
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
| Price | Number | - |
| Amount | Number | - |

### Ticker

| Name | Type | Description |
| ---- | ---- | ----------- |
| Bids | OrderBook List | bid `OrderBook` list |
| Buy | Number | the price of first `Bids` |
| Mid | Number | `(Buy + Sell) / 2` |
| Sell | Number | the price of first `Asks` |
| Asks | OrderBook List | ask `OrderBook` list |

## Global/G

`Global`/`G` is a object with some global methods.

### Sleep

> G.Sleep(Interval: *Any*) => *No Return*

```javascript
G.Sleep(5000);
// The program will sleep for 5 seconds
// If Interval <= 0, will automatic execute Es[i].AutoSleep()
```

### Log

> G.Log(Message: *Any*) => *No Return*

```javascript
G.Log("I'm running…");
// Send a message to web control
```

### LogProfit

> G.LogProfit(Profit: *Number*, Message: *Any*) => *No Return*

```javascript
G.LogProfit(12.345, "Round 1 end");
// Send a profit message to web control to show profit chart
```

### LogStatus

> G.LogStatus(Message: *Any*) => *No Return*

```javascript
G.LogStatus("Latest BTC Ticker: ", E.GetTicker("BTC"));
// Send a status message to web control to show it real-time
```

### AddTask

> G.AddTask(Function: *Function*, Arguments: *Any*) => *Boolean*

```javascript
// Work with G.ExecTasks()
```

### ExecTasks

> G.ExecTasks() => *List*

```javascript
G.AddTask(E.GetAccount);
G.AddTask(E.GetTicker, "BTC");
// Send a task to task list
var results = G.ExecTasks();
// Execute all tasks at the same time and return all results
var thisAccount = results[0];
var thisTicker = results[1];
```

## Exchange/E

`Exchange`/`E` is a object with some exchange methods.

### Log

> E.Log(Message: *Any*) => *No Return*

```javascript
E.Log("I'm running…");
// Send a message of this exchange to web control
```

### GetType

> E.GetType() => *String*

```javascript
var thisType = E.GetType();
// Get the type of this exchange
```

### GetName

> E.GetName() => *String*

```javascript
var thisName = E.GetName();
// Get the name of this exchange
```

### GetMainStock

> E.GetMainStock() => *String*

```javascript
var thisMainStock = E.GetMainStock();
// Get the main stock type of this exchange
```

### SetMainStock

> E.SetMainStock(StockType: *String*) => *String*

```javascript
var newMainStockType = E.SetMainStock("LTC");
// Set the main stock type of this exchange
```

### SetLimit

> E.SetLimit(times: *Number*) => *Number*

```javascript
var newLimit = E.SetLimit(6);
// Set the limit calls amount per second of this exchange
```

### SetLimit

> E.SetLimit(times: *Number*) => *Number*

```javascript
var newLimit = E.SetLimit(6);
// Set the limit calls amount per second of this exchange
// Work with E.AutoSleep()
```

### AutoSleep

> E.AutoSleep() => *No Return*

```javascript
E.AutoSleep();
// Auto sleep to achieve the limit calls amount per second of this exchange
```

### GetAccount

> E.GetAccount() => [`Account`](dataStruct.html#account)

```javascript
var thisAccount = E.GetAccount();
// Get the account info of this exchange
```

### Buy

> E.Buy(StockType: *String*, Price: *Number*, Amount: *Number*, Message: *Any*) => *String*/*Boolean*

```javascript
// if Price <= 0, it's a market order, and the Amount will be different
E.Buy("BTC", 600, 0.5, "I paid $300"); // normal order
E.Buy("BTC", 0, 300, "I also paid $300"); // market order
// Return ID of this order if succeed
// Return false if fail
```

### Sell

> E.Sell(StockType: *String*, Price: *Number*, Amount: *Number*, Message: *Any*) => *String*/*Boolean*

```javascript
// if Price <= 0, it's a market order
E.Sell("BTC", 600, 0.5); // normal order
E.Sell("BTC", 0, 0.5); // market order
// Return ID of this order if succeed
// Return false if fail
```

### GetOrder

> E.GetOrder(StockType: *String*, ID: *String*) => *Order*/*Boolean*

```javascript
var thisOrder = E.GetOrder("BTC", "XXXXXX");
// Return info of this order if succeed
// Return false if fail
```

### GetOrders

> E.GetOrders(StockType: *String*) => *Order List*

```javascript
var thisOrders = E.GetOrders("BTC");
// Return all the undone orders
```

### GetTrades

> E.GetTrades(StockType: *String*) => *Order List*

```javascript
var thisTrades = E.GetTrades("BTC");
// Return all the done orders
```

### CancelOrder

> E.CancelOrder(Order: *Order*) => *Boolean*

```javascript
var thisOrders = E.GetOrders("BTC");
for (var i = 0; i < thisOrders.length; i++) {
var isCanceled = E.CancelOrder(thisOrders[i]);
// Return the result
}
```

### GetTicker

> E.GetTicker(StockType: *String*, Size: *Any*) => *Ticker*

```javascript
var thisTicker = E.GetTicker("BTC");
// Get the latest ticker of this exchange
```

### GetRecords

> E.GetRecords(StockType: *String*, Period: [*String*](#records-period), Size: *Any*) => *Record List*

```javascript
var thisRecords = E.GetRecords("BTC", "M5");
// Get the latest records of this exchange
```
