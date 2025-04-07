package OpenRouterAPI

import (
	"os"
	"testing"
)

func TestGetCredits(t *testing.T) {
	key := os.Getenv("OPENROUTER_API_KEY")
	if key == "" {
		t.Skip("OPENROUTER_API_KEY not set")
	}
	c := GetCredits(key)
	if c < 0 {
		t.Error("GetCredits failed")
	}
	t.Log(c)
}
