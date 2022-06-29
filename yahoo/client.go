package yahoo

import (
	"errors"
	"fmt"
	"reactiveNews"
)

type Client interface {
	GetRealTimeQuoteData(tickers []string) (QuoteDataResponse, error)
}

type Res struct {
	Ask                          float32 `json:"ask"`
	AskSize                      float32 `json:"askSize"`
	Bid                          float32 `json:"bid"`
	BidSize                      float32 `json:"bidSize"`
	Currency                     string  `json:"currency"`
	DisplayName                  string  `json:"displayName"`
	FiftyDayAverage              float32 `json:"fiftyDayAverage"`
	FiftyDayAverageChange        float32 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent float32 `json:"fiftyDayAverageChangePercent"`
}

type QResp struct {
	Result []Res `json:"result"`
}

type QuoteDataResponse struct {
	QuoteResponse QResp `json:"quoteResponse"`
}

func (q QuoteDataResponse) MarketData() ([]reactiveNews.MarketData, error) {
	var (
		errs []error
		res  = make([]reactiveNews.MarketData, 0, len(q.QuoteResponse.Result))
	)

	for _, r := range q.QuoteResponse.Result {
		c, err := reactiveNews.ToCurrency(r.Currency)
		if err != nil {
			errs = append(errs, fmt.Errorf("get currency for %s: %w", r.DisplayName, err))
		}

		res = append(res, reactiveNews.MarketData{
			Ask:                          r.Ask,
			AskSize:                      r.AskSize,
			Bid:                          r.Bid,
			BidSize:                      r.BidSize,
			Currency:                     c,
			DisplayName:                  r.DisplayName,
			FiftyDayAverage:              r.FiftyDayAverage,
			FiftyDayAverageChange:        r.FiftyDayAverageChange,
			FiftyDayAverageChangePercent: r.FiftyDayAverageChangePercent,
		})
	}

	if len(errs) != 0 {
		err := errors.New("map from Yahoo response")
		for _, e := range errs {
			err = fmt.Errorf("%s, %s", e, err)
		}

		return res, err
	}

	return res, nil
}
