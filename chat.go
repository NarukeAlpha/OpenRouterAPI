package OpenRouterAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func orint(i int) *int           { return &i }
func orfloat(i float64) *float64 { return &i }
func orbool(b bool) *bool        { return &b }

type Role string

const (
	RoleSystem    Role = "system"
	RoleDeveloper Role = "developer"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type ChatResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model             string             `json:"model"`
	Messages          []Message          `json:"messages"`
	Stream            *bool              `json:"stream,omitempty"`
	MaxTokens         *int               `json:"max_tokens,omitempty"`
	Temperature       *float64           `json:"temperature,omitempty"`
	Seed              *int               `json:"seed,omitempty"`
	TopP              *float64           `json:"top_p,omitempty"`
	TopK              *int               `json:"top_k,omitempty"`
	FrequencyPenalty  *float64           `json:"frequency_penalty,omitempty"`
	PresencePenalty   *float64           `json:"presence_penalty,omitempty"`
	RepetitionPenalty *float64           `json:"repetition_penalty,omitempty"`
	LogitBias         map[string]float64 `json:"logit_bias,omitempty"`
	TopLogprobs       *int               `json:"top_logprobs,omitempty"`
	MinP              *float64           `json:"min_p,omitempty"`
	TopA              *float64           `json:"top_a,omitempty"`
	Transforms        []string           `json:"transforms,omitempty"`
	Models            []string           `json:"models,omitempty"`
	Provider          *ProviderPrefs     `json:"provider,omitempty"`
	Reasoning         *ReasoningConfig   `json:"reasoning,omitempty"`
}

type ProviderPrefs struct {
	Sort string `json:"sort,omitempty"`
}

type ReasoningConfig struct {
	Effort    string `json:"effort,omitempty"`
	MaxTokens *int   `json:"max_tokens,omitempty"`
	Exclude   *bool  `json:"exclude,omitempty"`
}

func Generate(key string, data ChatRequest) (ChatResponse, error) {
	var response ChatResponse

	b, err := json.Marshal(data)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to marshal chat request"))
	}
	u, err := url.Parse("https://openrouter.ai/api/v1/chat/completions")
	if err != nil {
		return response, errors.Join(err, errors.New("failed to parse url"))
	}

	r, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return response, errors.Join(err, errors.New("failed to create request"))
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", key)

	c := &http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to make request for chat gen"))
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to decode response for chat gen"))
	}

	return response, nil

}
