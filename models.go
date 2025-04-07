package OpenRouterAPI

import (
	"encoding/json"
	"io"
	"net/http"
)

type ModelList struct {
	Data []Model `json:"data"`
}

type Model struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Created          float64           `json:"created"`
	Description      string            `json:"description"`
	ContextLength    float64           `json:"context_length"`
	Architecture     Architecture      `json:"architecture"`
	TopProvider      *TopProvider      `json:"top_provider,omitempty"`
	Pricing          Pricing           `json:"pricing"`
	PerRequestLimits *PerRequestLimits `json:"per_request_limits,omitempty"`
	Endpoints        []Endpoint        `json:"endpoints,omitempty"`
}

type Architecture struct {
	Modality         string   `json:"modality,omitempty"`
	Tokenizer        string   `json:"tokenizer"`
	InputModalities  []string `json:"input_modalities,omitempty"`
	OutputModalities []string `json:"output_modalities,omitempty"`
	InstructType     string   `json:"instruct_type,omitempty"`
}

type TopProvider struct {
	ContextLength       int  `json:"context_length"`
	MaxCompletionTokens int  `json:"max_completion_tokens"`
	IsModerated         bool `json:"is_moderated"`
}

type Pricing struct {
	Prompt            string  `json:"prompt"`
	Completion        string  `json:"completion"`
	Image             string  `json:"image"`
	Request           string  `json:"request"`
	InputCacheRead    string  `json:"input_cache_read,omitempty"`
	InputCacheWrite   string  `json:"input_cache_write,omitempty"`
	WebSearch         string  `json:"web_search,omitempty"`
	InternalReasoning string  `json:"internal_reasoning,omitempty"`
	Discount          float64 `json:"discount"`
}

type PerRequestLimits struct {
	Key string `json:"key"`
}

type Endpoints struct {
	Name                string   `json:"name"`
	ContextLength       int      `json:"context_length"`
	Pricing             Pricing  `json:"pricing"`
	ProviderName        string   `json:"provider_name"`
	Quantization        string   `json:"quantization"`
	MaxCompletionTokens *int     `json:"max_completion_tokens"`
	MaxPromptTokens     *int     `json:"max_prompt_tokens"`
	SupportedParameters []string `json:"supported_parameters"`
	Status              int      `json:"status"`
}

type Endpoint struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Created      float64      `json:"created"`
	Description  string       `json:"description"`
	Architecture Architecture `json:"architecture"`
	Endpoints    []Endpoints  `json:"endpoints"`
}

type EndpointList struct {
	Data Endpoint `json:"data"`
}

func GetModelsAvail() ModelList {
	var list ModelList
	resp, err := http.Get("https://openrouter.ai/api/v1/models")
	if err != nil {
		return list
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return list
	}
	return list
}

func GetModelEndpoints(author string, model string) (EndpointList, error) {
	var list EndpointList
	url := "https://openrouter.ai/api/v1/models/" + author + "/" + model + "/endpoints"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return list, err
	}

	req.Header.Set("User-Agent", "curl/7.79.1")

	resp, err := client.Do(req)
	if err != nil {
		return list, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return list, err
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}

	return list, nil
}
