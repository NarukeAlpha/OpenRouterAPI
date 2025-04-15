package OpenRouterAPI

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ORInt(i int) *int           { return &i }
func ORFloat(i float64) *float64 { return &i }
func ORBool(b bool) *bool        { return &b }
func ORRole(r Role) string {
	s := string(r)
	return s
}

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
	auth := "Bearer " + key
	r.Header.Set("Authorization", auth)

	c := &http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to make request for chat gen"))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to decode response for chat gen"))
	}

	return response, nil

}

func GenerateStream(key string, data ChatRequest, w io.Writer) (ChatResponse, error) {
	var response ChatResponse
	data.Stream = ORBool(true)

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
	auth := "Bearer " + key
	r.Header.Set("Authorization", auth)

	c := &http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to make request for chat gen"))
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "stream_copy_*.tmp")
	if err != nil {
		return response, errors.New("failed to create temp file")
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	mw := io.MultiWriter(w, tmpFile)
	_, err = io.Copy(mw, resp.Body)
	if err != nil {
		return response, errors.Join(err, errors.New("failed to stream response"))
	}

	_, err = tmpFile.Seek(0, io.SeekStart)
	if err != nil {
		return response, errors.New("failed to rewind temp file")
	}

	s := bufio.NewScanner(tmpFile)
	var accumulatedContent string
	var responseID string
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "data: ") {
			dataPart := strings.TrimPrefix(line, "data: ")
			if dataPart == "[DONE]" {
				continue
			}
			var chunk struct {
				ID      string `json:"id"`
				Choices []struct {
					Delta struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			err := json.Unmarshal([]byte(dataPart), &chunk)
			if err != nil {
				continue // skip malformed chunks
			}
			if chunk.ID != "" {
				responseID = chunk.ID
			}
			if len(chunk.Choices) > 0 {
				if chunk.Choices[0].Delta.Role != "" {
					response.Choices = []Choice{{Message: Message{Role: chunk.Choices[0].Delta.Role}}}
				}
				accumulatedContent += chunk.Choices[0].Delta.Content
			}
		}
	}
	if err := s.Err(); err != nil {
		return response, errors.New("failed to scan temp file")
	}
	if accumulatedContent == "" {
		return response, errors.New("no content accumulated from streaming response")
	}
	response.ID = responseID
	if len(response.Choices) > 0 {
		response.Choices[0].Message.Content = accumulatedContent
	} else {
		// Fallback if role was somehow not captured (should be unlikely)
		response.Choices = []Choice{{Message: Message{Content: accumulatedContent}}}
	}

	return response, nil
}
