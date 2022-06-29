package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reactiveNews/yahoo"
	"strings"
)

var _ yahoo.Client = Service{}

type Service struct {
	baseUrl    string
	httpClient http.Client
	apiKey     string
}

func NewService(client http.Client, url, apiKey string) Service {
	return Service{
		baseUrl:    url,
		httpClient: client,
		apiKey:     apiKey,
	}
}

func (s Service) GetRealTimeQuoteData(tickers []string) (yahoo.QuoteDataResponse, error) {
	url := fmt.Sprintf("%s/v6/finance/quote?symbols=%s", s.baseUrl, strings.Join(tickers, ","))

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return yahoo.QuoteDataResponse{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return yahoo.QuoteDataResponse{}, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if !isStatusSuccessful(resp.StatusCode) {
		return yahoo.QuoteDataResponse{}, fmt.Errorf("service responded with %v", resp.Status)
	}

	var yahooResp yahoo.QuoteDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		return yahoo.QuoteDataResponse{}, fmt.Errorf("decode response: %w", err)
	}

	return yahooResp, nil
}

func isStatusSuccessful(statusCode int) bool {
	return 200 <= statusCode && statusCode < 300
}
