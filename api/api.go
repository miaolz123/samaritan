package api

// Option : exchange option
type Option struct {
	TraderID  uint
	Type      string // one of ["okcoin.cn", "huobi"]
	AccessKey string
	SecretKey string
	MainStock string
}
