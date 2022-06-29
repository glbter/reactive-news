package indefinite

import (
	"context"
	"reactiveNews/tickers"
	"time"
)

var _ tickers.Service = &Service{}

type Service struct {
	ts       tickers.Service
	interval time.Duration
}

func NewService(service tickers.Service, interval time.Duration) Service {
	return Service{
		ts:       service,
		interval: interval,
	}
}

func (s *Service) Add(ctx context.Context, ticker string) error {
	return s.ts.Add(ctx, ticker)
}

func (s *Service) Remove(ctx context.Context, ticker string) error {
	return s.ts.Remove(ctx, ticker)
}

func (s *Service) GetAllChan(ctx context.Context) (<-chan string, error) {
	resCh := make(chan string, 1)
	iTicker := time.NewTicker(s.interval)
	go func() {
		defer close(resCh)
		defer iTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-iTicker.C:
				ch, err := s.ts.GetAllChan(ctx)
				if err != nil {
					return
				}

				for tick := range ch {
					resCh <- tick
				}
			}
		}
	}()

	return resCh, nil
}
