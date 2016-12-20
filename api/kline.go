package api

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/samaritan/constant"
)

var sosobtcSymbolMap = map[string]map[string]string{
	constant.Btcc: {
		"BTC/CNY": "btcchinabtccny",
		"LTC/CNY": "btcchinaltccny",
	},
}

var sosobtcPeriodMap = map[string]int{
	"M":   60,
	"M3":  180,
	"M5":  300,
	"M10": 600,
	"M15": 900,
	"M30": 1800,
	"H":   3600,
	"H2":  7200,
	"H4":  14400,
	"H6":  21600,
	"H12": 43200,
	"D":   86400,
	"D3":  259200,
	"W":   604800,
}

func getSosobtcRecords(recordsOld []Record, exchangeType, stockType, period string, size int) (records []Record, err error) {
	periodInt := sosobtcPeriodMap[period]
	if periodInt == 0 {
		err = fmt.Errorf("unrecognized period: %v", period)
		return
	}
	symbol := sosobtcSymbolMap[exchangeType][stockType]
	if symbol == "" {
		err = fmt.Errorf("unrecognized stockType: %v", stockType)
		return
	}
	resp, err := get(fmt.Sprintf("http://k.sosobtc.com/data/period?symbol=%v&step=%v", symbol, periodInt))
	if err != nil {
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		return
	}
	timeLast := int64(0)
	if len(recordsOld) > 0 {
		timeLast = recordsOld[len(recordsOld)-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		recordTime := recordJSON.GetIndex(0).MustInt64()
		if recordTime > timeLast {
			recordsNew = append([]Record{{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			recordsOld[len(recordsOld)-1] = Record{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}
		} else {
			break
		}
	}
	records = append(recordsOld, recordsNew...)
	if len(records) > size {
		records = records[len(records)-size:]
	}
	return
}
