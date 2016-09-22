package constant

const (
	// model.Log.Type : ["error", "info", "profit", "buy", "sell", "cancel"]
	ERROR  = -1
	INFO   = 0
	PROFIT = 1
	BUY    = 2
	SELL   = 3
	CANCEL = 4
)

var (
	CONSTS = []string{"BTC", "LTC", "M", "M5", "M15", "M30", "H", "D", "W"}
)
