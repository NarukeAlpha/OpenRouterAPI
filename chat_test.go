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
	payload.Model = model.Data[3].ID
	payload.MaxTokens = orint(100)
	payload.Stream = orbool(false)
	payload.Temperature = orfloat(0.9)

	for _, x := range rolelist {
		t.Run("Testing Prompt Generation", func(t *testing.T) {
			var m Message
			m.Content = "This is a test, reply something random"
			m.Role = orrole(x)
			payload.Messages = append(payload.Messages, m)
			resp, err := Generate(key, payload)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp)
		})
	}
}
