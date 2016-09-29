package constant

const (
	// logs type
	ERROR  = -1
	INFO   = 0
	PROFIT = 1
	BUY    = 2
	SELL   = 3
	CANCEL = 4
	// exchange type
	OkCoinCn = "okcoin.cn"
	Huobi    = "huobi"
	// stock type
	BTC = "BTC"
	LTC = "LTC"
	// delete logsTime type
	LastTime = "0"
	Day      = "1"
	Week     = "2"
	Month    = "3"
	// task name
	GetAccount = "GetAccount"
	Buy        = "Buy"
	Sell       = "Sell"
)

var (
	// CONSTS : Javascript Global Constants
	CONSTS = []string{"BTC", "LTC", "M", "M5", "M15", "M30", "H", "D", "W"}
)
