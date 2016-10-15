package api

import "github.com/miaolz123/samaritan/constant"

func simulateTrade(tradeType string, amount float64, ticker Ticker) (total float64) {
	switch tradeType {
	case constant.TradeTypeBuy:
		if len(ticker.Asks) < 1 {
			return ticker.Sell * amount
		}
		dealAmount := 0.0
		for _, orderBook := range ticker.Asks {
			if dealAmount+orderBook.Amount >= amount {
				total += (amount - dealAmount) * orderBook.Price
				dealAmount = amount
				break
			}
			total += orderBook.Amount * orderBook.Price
			dealAmount += orderBook.Amount
		}
		if dealAmount < amount {
			total += (amount - dealAmount) * ticker.Asks[len(ticker.Asks)-1].Price
		}
	case constant.TradeTypeSell:
		if len(ticker.Bids) < 1 {
			return ticker.Buy * amount
		}
		dealAmount := 0.0
		for _, orderBook := range ticker.Bids {
			if dealAmount+orderBook.Amount >= amount {
				total += (amount - dealAmount) * orderBook.Price
				dealAmount = amount
				break
			}
			total += orderBook.Amount * orderBook.Price
			dealAmount += orderBook.Amount
		}
		if dealAmount < amount {
			total += (amount - dealAmount) * ticker.Bids[len(ticker.Bids)-1].Price
		}
	}
	return
}

func simulateBuy(amount float64, ticker Ticker) (total float64) {
	return
}

func simulateSell(amount float64, ticker Ticker) (total float64) {
	return
}
