package OpenRouterAPI

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Balance struct {
	Data Data `json:"data"`
}
type Data struct {
	TotalCredits float64 `json:"total_credits"`
	TotalUsage   float64 `json:"total_usage"`
}

func GetCredits(key string) float64 {
	var balance Balance
	var c http.Client
	u, err := url.Parse("https://openrouter.ai/api/v1/credits")
	if err != nil {
		return -3
	}
	auth := "Bearer " + key
	r := http.Request{
		Header: http.Header{
			"Authorization": []string{auth},
		},
		Method: "GET",
		URL:    u,
	}
	resp, err := c.Do(&r)
	if err != nil {
		return -4
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -5
	}

	err = json.Unmarshal(body, &balance)
	if err != nil {
		return -6
	}

	return balance.Data.TotalUsage
}
