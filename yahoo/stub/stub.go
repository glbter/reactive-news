package stub

import (
	"fmt"
	"reactiveNews/yahoo"
)

var _ yahoo.Client = &Stub{}

type Stub struct {
	counter int
}

func (s *Stub) GetRealTimeQuoteData(_ []string) (yahoo.QuoteDataResponse, error) {
	s.counter = s.counter + 1
	return yahoo.QuoteDataResponse{
		QuoteResponse: yahoo.QResp{
			Result: []yahoo.Res{
				{136.87, 9, 186.86, 8, "USD", fmt.Sprintf("Apple:%v", s.counter), 150.2476, -14.052597, -0.09352959},
				{701.38, 10, 701.01, 14, "USD", "Tesla", 801.7334, -109.063416, -0.13603452},
				//{157.06, 9, 157.15, 8, "USD",  "", 192.3466, -35.896606, -0.18662459},
				//{157.06, 9, 157.15, 8, "UAH",  "", 192.3466, -35.896606, -0.18662459},

				//				{177.62 8 177.32 8 "USD" "Netflix" 208.39 -32.054993 -0.15382212},
				//				{2246 9 2242.59 10 "USD" "Alphabet" 2324.2117 -93.48169 -0.040220816},
				//				{110.62 9 110.65 31 "USD" "Amazon.com" 122.53241 -12.352409 -0.10080932},
				//				{11.35 18 11.34 40 "USD"  13.688 -2.3979998 -0.17518993},
				//				{0, 0, 0, 0, "USD", 7.9728, 0.16220045, 0.020344228,},
			},
		},
	}, nil
}
