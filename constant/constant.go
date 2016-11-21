package constant

const (
	// log type
	ERROR  = -1
	INFO   = 0
	PROFIT = 1
	BUY    = 2
	SELL   = 3
	CANCEL = 4

	// delete log time type
	LastTime = "0"
	Day      = "1"
	Week     = "2"
	Month    = "3"

	// exchange type
	OkCoinCn     = "okcoin.cn"
	Huobi        = "huobi"
	Poloniex     = "poloniex"
	Btcc         = "btcc"
	Chbtc        = "chbtc"
	OkcoinFuture = "okcoin.future"
	OandaV20     = "oanda.v20"

	// trade type
	TradeTypeBuy        = "BUY"
	TradeTypeSell       = "SELL"
	TradeTypeLong       = "LONG"
	TradeTypeShort      = "SHORT"
	TradeTypeLongClose  = "LONG_CLOSE"
	TradeTypeShortClose = "SHORT_CLOSE"

	// stock type (will useless)
	BTC = "BTC"
	LTC = "LTC"
)

var (
	// CONSTS : Javascript Global Constants
	CONSTS        = []string{"BTC", "LTC", "M", "M5", "M15", "M30", "H", "D", "W"}
	ExchangeTypes = []string{"okcoin.cn", "huobi", "poloniex", "btcc", "chbtc", "okcoin.future", "oanda.v20"}
)
