package reactiveNews

import "fmt"

type MarketPrice struct {
	Ask         float32
	Bid         float32
	Currency    Currency
	DisplayName string
}

type MarketRise struct {
	Ask             float32
	Bid             float32
	Currency        Currency
	DisplayName     string
	FiftyDayAverage float32
}

type MarketData struct {
	Ask                          float32
	AskSize                      float32
	Bid                          float32
	BidSize                      float32
	Currency                     Currency
	DisplayName                  string
	FiftyDayAverage              float32
	FiftyDayAverageChange        float32
	FiftyDayAverageChangePercent float32
}

func (m MarketData) MarketPrice() MarketPrice {
	return MarketPrice{
		Ask:         m.Ask,
		Bid:         m.Bid,
		Currency:    m.Currency,
		DisplayName: m.DisplayName,
	}
}

type Currency string

func (c Currency) String() string {
	return string(c)
}

func ToCurrency(s string) (Currency, error) {
	switch s {
	case USD.String():
		return USD, nil
	default:
		return "", fmt.Errorf("unknown currency %q", s)
	}
}

const (
	USD Currency = "USD"
)
