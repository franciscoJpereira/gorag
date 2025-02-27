package apiinterface_test

import (
	"encoding/json"
	apiinterface "ragAPI/pkg/apiInterface"
	"testing"
)

func TestJSONMarshalAndUnmarshal(t *testing.T) {
	message := apiinterface.ChatMessage{
		Role:    "SomeRole",
		Content: "Nice Message Content",
	}

	marshaled, err := json.Marshal(message)

	if err != nil {
		t.Fatalf("Failed marshaling message: %s\n", marshaled)
	}

	var sameMessage apiinterface.ChatMessage

	err = json.Unmarshal(marshaled, &sameMessage)
	if err != nil {
		t.Fatalf("Failed unmarshaling: %s\n", err)
	}

	if sameMessage.Content != message.Content || sameMessage.Role != message.Role {
		t.Fatalf("Different message values after unmarshaling: %v vs %v\n", sameMessage, message)
	}
}
