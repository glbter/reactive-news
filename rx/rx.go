package rx

import "github.com/reactivex/rxgo/v2"

func FromStringChan(chs <-chan string) <-chan rxgo.Item {
	chi := make(chan rxgo.Item, 1)
	go func() {
		defer close(chi)

		for s := range chs {
			chi <- rxgo.Of(s)
		}
	}()

	return chi
}
