package OpenRouterAPI

import (
	"strings"
	"testing"
)

func TestGetModelsAvail(t *testing.T) {
	var list ModelList
	t.Run("Model Endpoint", func(t *testing.T) {
		list = GetModelsAvail()
		if list.Data == nil {
			t.Fatal("No models returned")

		}
		t.Run("Model Endpoint", func(t *testing.T) {
			model := strings.Split(list.Data[3].ID, "/")
			mlist, err := GetModelEndpoints(model[0], model[1])
			if err != nil {
				t.Fatal(err, "No endpoints returned")
			}
			t.Log(mlist)
		})
	})

}
