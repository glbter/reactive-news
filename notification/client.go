package notification

import "reactiveNews"

type Client interface {
	SendRealTimeData(data reactiveNews.MarketPrice) error
	SendPriceChange(change reactiveNews.MarketRise) error
}
