package OpenRouterAPI

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Amount struct {
	TotalCredits float32 `json:"total_credits"`
	TotalUsage   float32 `json:"total_usage"`
}

type Balance struct {
	Data []Amount `json:"data"`
}

func GetCredits(key string) float32 {
	var balance Balance
	var c http.Client
	u, err := url.Parse("https://openrouter.ai/api/v1/credits")
	if err != nil {
		return -3
	}
	r := http.Request{
		Header: http.Header{
			"Authorization": []string{key},
		},
		Method: "GET",
		URL:    u,
	}
	resp, err := c.Do(&r)
	if err != nil {
		return -4
	}
	err = json.NewDecoder(resp.Body).Decode(&balance)
	if err != nil {
		return -5
	}

	return balance.Data[0].TotalCredits
}
