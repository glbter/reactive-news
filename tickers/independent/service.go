package independent

import (
	"context"
	"reactiveNews/tickers"
)

var _ tickers.Service = &Service{}

type Service struct {
	tickers map[string]struct{}
}

func NewService() *Service {
	//t := make([]string, 0)
	t := make(map[string]struct{})
	return &Service{
		tickers: t,
	}
}

func (s *Service) Add(_ context.Context, ticker string) error {
	s.tickers[ticker] = struct{}{}
	//s.tickers = append(s.tickers, ticker)
	return nil
}

func (s *Service) Remove(_ context.Context, ticker string) error {
	delete(s.tickers, ticker)
	return nil
}

func (s Service) GetAllChan(_ context.Context) (<-chan string, error) {
	ch := make(chan string, 1)
	go func() {
		defer close(ch)

		for t := range s.tickers {
			ch <- t
		}
	}()

	return ch, nil
}
