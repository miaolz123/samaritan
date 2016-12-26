# Samaritan

[![Travis](https://img.shields.io/travis/miaolz123/samaritan.svg)](https://travis-ci.org/miaolz123/samaritan) [![Go Report Card](https://goreportcard.com/badge/github.com/miaolz123/samaritan)](https://goreportcard.com/report/github.com/miaolz123/samaritan) [![Github All Releases](https://img.shields.io/github/downloads/miaolz123/samaritan/total.svg)](https://github.com/miaolz123/samaritan/releases) [![Docker Pulls](https://img.shields.io/docker/pulls/miaolz123/samaritan.svg)](https://hub.docker.com/r/miaolz123/samaritan/) [![license](https://img.shields.io/github/license/miaolz123/samaritan.svg)](https://github.com/miaolz123/samaritan/blob/master/LICENSE)

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

## Algorithm Reference

[Read Documentation](http://samaritan.stockdb.org/#/#algorithm-reference)

## Contributing

Contributions are not accepted in principle until the basic infrastructure is complete.

However, the [ISSUE](https://github.com/miaolz123/samaritan/issues) is welcome.

## License

Copyright (c) 2016 [miaolz123](https://github.com/miaolz123) by MIT
