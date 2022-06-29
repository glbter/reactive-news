package main

import (
	"context"
	"fmt"
	"reactiveNews"
	"reactiveNews/notification"
	notifConsole "reactiveNews/notification/console"
	"reactiveNews/rx"
	infTickers "reactiveNews/tickers/indefinite"
	tickers "reactiveNews/tickers/independent"
	"reactiveNews/yahoo"
	yahooCl "reactiveNews/yahoo/stub"
	"time"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	yahooClient := yahooCl.Stub{}
	//yahooClient := yahooHttp.NewService(http.Client{}, "https://yfapi.net", "qcsa6I5EEu8nb5LEloxxn8M8XnZskVTl3yqWbnRL")

	ctx := context.Background()

	notifService := notifConsole.Console{}

	infTicker := infTickers.NewService(tickers.NewService(), time.Second*2)

	for _, t := range []string{"AAPL", "TSLA", "META", "NFLX", "GOOG", "AMZN", "F", "NSANY"} {
		infTicker.Add(ctx, t)
	}

	ch, _ := infTicker.GetAllChan(ctx)

	obs := rxgo.
		FromChannel(rx.FromStringChan(ch)).
		WindowWithTimeOrCount(rxgo.WithDuration(time.Second*2), 10).
		Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			return i.(rxgo.Observable).ToSlice(10)
		}).
		Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			sli := i.([]interface{})
			sls := make([]string, 0, len(sli))
			for _, s := range sli {
				sls = append(sls, s.(string))
			}

			return yahooClient.GetRealTimeQuoteData(sls)
		}, rxgo.WithObservationStrategy(rxgo.Lazy)).
		FlatMap(func(item rxgo.Item) rxgo.Observable {
			res, err := item.V.(yahoo.QuoteDataResponse).MarketData()
			if err != nil {
				fmt.Println(err)
			}

			mapped := make([]interface{}, 0, len(res))
			for _, r := range res {
				mapped = append(mapped, r)
			}

			return rxgo.Just(mapped...)()
		})

	time.Sleep(time.Second * 3)

	SendMarketRise(GetMarketRise(NotificationStream(obs)), notifService)

	SendRealTimeData(RealTimeData(obs), notifService)

	time.Sleep(time.Second * 10)
}

func NotificationStream(marketDataObs rxgo.Observable) rxgo.Observable {
	return marketDataObs.Filter(func(i interface{}) bool {
		data := i.(reactiveNews.MarketData)

		if data.Bid-data.FiftyDayAverage > data.FiftyDayAverage*0.001 {
			return true
		}

		return false
	})
}

func GetMarketRise(marketDataObs rxgo.Observable) rxgo.Observable {
	return marketDataObs.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		md := i.(reactiveNews.MarketData)
		return reactiveNews.MarketRise{
			Ask:             md.Ask,
			Bid:             md.Bid,
			Currency:        md.Currency,
			DisplayName:     md.DisplayName,
			FiftyDayAverage: md.FiftyDayAverage,
		}, nil
	}, rxgo.WithBufferedChannel(5))
}

func SendMarketRise(marketRiseObs rxgo.Observable, client notification.Client) {
	marketRiseObs.
		Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			mp := i.(reactiveNews.MarketRise)
			return nil, client.SendPriceChange(mp)
		}).
		DoOnError(func(err error) {
			fmt.Println("got market change error: ", err)
		})
}

func RealTimeData(marketDataObs rxgo.Observable) rxgo.Observable {
	return marketDataObs.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		return i.(reactiveNews.MarketData).MarketPrice(), nil
	})
}

func SendRealTimeData(marketPriceObs rxgo.Observable, client notification.Client) {
	marketPriceObs.
		Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			mp := i.(reactiveNews.MarketPrice)
			return nil, client.SendRealTimeData(mp)
		}).
		DoOnError(func(err error) {
			fmt.Println("got real time error ", err)
		})
}
