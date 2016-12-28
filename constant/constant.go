package constant

// error constants
const (
	Banner = `   _____                            _ __            
  / ___/____ _____ ___  ____ ______(_/ /_____ _____ 
  \__ \/ __ ` + "`/ __ `__ \\/ __ `/ ___/ / __/ __ `" + `/ __ \
 ___/ / /_/ / / / / / / /_/ / /  / / /_/ /_/ / / / /
/____/\__,_/_/ /_/ /_/\__,_/_/  /_/\__/\__,_/_/ /_/`
	Version                    = "0.1.2"
	ErrAuthorizationError      = "Authorization Error"
	ErrInsufficientPermissions = "Insufficient Permissions"
)

// exchange types
const (
	OkCoinCn     = "okcoin.cn"
	Huobi        = "huobi"
	Poloniex     = "poloniex"
	Btcc         = "btcc"
	Chbtc        = "chbtc"
	OkcoinFuture = "okcoin.future"
	OandaV20     = "oanda.v20"
)

// log types
const (
	ERROR      = "ERROR"
	INFO       = "INFO"
	PROFIT     = "PROFIT"
	BUY        = "BUY"
	SELL       = "SELL"
	LONG       = "LONG"
	SHORT      = "SHORT"
	LONGCLOSE  = "LONG_CLOSE"
	SHORTCLOSE = "SHORT_CLOSE"
	CANCEL     = "CANCEL"
)

// delete log time types
const (
	LastTime = "0"
	Day      = "1"
	Week     = "2"
	Month    = "3"
)

// trade types
const (
	TradeTypeBuy        = "BUY"
	TradeTypeSell       = "SELL"
	TradeTypeLong       = "LONG"
	TradeTypeShort      = "SHORT"
	TradeTypeLongClose  = "LONG_CLOSE"
	TradeTypeShortClose = "SHORT_CLOSE"
)

// stock types (will useless)
const (
	BTC = "BTC"
	LTC = "LTC"
)

// some variables
var (
	Consts        = []string{"BTC", "LTC", "M", "M5", "M15", "M30", "H", "D", "W"}
	ExchangeTypes = []string{OkCoinCn, Huobi, Poloniex, Btcc, Chbtc, OkcoinFuture}
)
