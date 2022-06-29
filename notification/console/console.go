package console

import (
	"fmt"
	"reactiveNews"
	"reactiveNews/notification"
)

var _ notification.Client = Console{}

type Console struct{}

func (c Console) SendRealTimeData(data reactiveNews.MarketPrice) error {
	fmt.Println(fmt.Sprintf("%#v", data))
	return nil
}

func (c Console) SendPriceChange(change reactiveNews.MarketRise) error {
	fmt.Println(fmt.Sprintf("%#v", change))
	return nil
}
