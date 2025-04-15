package OpenRouterAPI

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	key := os.Getenv("OPENROUTER_API_KEY")
	var rolelist = []Role{
		RoleSystem,
		RoleDeveloper,
		RoleUser,
		RoleAssistant,
		RoleTool,
	}
	model := GetModelsAvail()
	var payload ChatRequest
	payload.Model = model.Data[2].ID
	payload.MaxTokens = ORInt(100)
	payload.Stream = ORBool(false)
	payload.Temperature = ORFloat(0.9)

	for _, x := range rolelist {
		t.Run("Testing Prompt Generation", func(t *testing.T) {
			var m Message
			m.Content = "This is a test, reply something random"
			m.Role = ORRole(x)
			payload.Messages = append(payload.Messages, m)
			resp, err := Generate(key, payload)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp)
		})
	}
}

func TestGenerateStream(t *testing.T) {
	key := os.Getenv("OPENROUTER_API_KEY")
	if key == "" {
		t.Skip("OPENROUTER_API_KEY not set")
	}

	model := GetModelsAvail()
	var payload ChatRequest
	payload.Model = model.Data[9].ID
	payload.MaxTokens = ORInt(1000) // Set max tokens to 5000
	payload.Temperature = ORFloat(0.9)

	// Add test message
	var m Message
	m.Content = "Stream this test response, please write quantum physics for about 1k tokens"
	m.Role = ORRole(RoleUser)
	payload.Messages = append(payload.Messages, m)

	t.Run("Streaming Test", func(t *testing.T) {
		resp, err := GenerateStream(key, payload, os.Stdout) // Stream to stdout
		if err != nil {
			t.Fatal(err)
		}
		t.Log(resp.Choices[0])
	})
}
