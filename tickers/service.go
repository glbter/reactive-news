package tickers

import "context"

type Service interface {
	Add(ctx context.Context, ticker string) error
	Remove(ctx context.Context, ticker string) error
	GetAllChan(ctx context.Context) (<-chan string, error)
}
