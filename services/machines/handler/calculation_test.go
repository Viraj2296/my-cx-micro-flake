package handler

import (
	"encoding/json"
	"testing"
)

func getCycleCountWithZeroReceived(t *testing.T) {
	jsonBody1, _ := json.Marshal(map[string]interface{}{"cycle_count": 345, "in_timestamp": 1687509464492})
	jsonBody2, _ := json.Marshal(map[string]interface{}{"cycle_count": 346, "in_timestamp": 1687509464493})
	jsonBody3, _ := json.Marshal(map[string]interface{}{"cycle_count": 0, "in_timestamp": 1687509464494})
	jsonBody4, _ := json.Marshal(map[string]interface{}{"cycle_count": 1, "in_timestamp": 1687509464495})
	jsonBody5, _ := json.Marshal(map[string]interface{}{"cycle_count": 2, "in_timestamp": 1687509464496})
	messages := []Message{Message{Topic: "A", TS: 1687509464492, Body: jsonBody1}, Message{Topic: "A", TS: 1687509464493, Body: jsonBody2}, Message{Topic: "A", TS: 1687509464494, Body: jsonBody3}, Message{Topic: "A", TS: 1687509464495, Body: jsonBody4}, Message{Topic: "A", TS: 1687509464496, Body: jsonBody5}}

	result := messagePreprocessing(messages)

	expected := 5

	if len(result) != expected {
		t.Errorf("Add(2, 3) returned %d, expected %d", len(result), expected)
	} else {
		t.Logf("Success %d", len(result))
	}

}
func getCycleCountWithZeroNotReceived(t *testing.T) {
	jsonBody1, _ := json.Marshal(map[string]interface{}{"cycle_count": 345, "in_timestamp": 1687509464492})
	jsonBody2, _ := json.Marshal(map[string]interface{}{"cycle_count": 346, "in_timestamp": 1687509464493})
	jsonBody4, _ := json.Marshal(map[string]interface{}{"cycle_count": 1, "in_timestamp": 1687509464495})
	jsonBody5, _ := json.Marshal(map[string]interface{}{"cycle_count": 2, "in_timestamp": 1687509464496})
	messages := []Message{Message{Topic: "A", TS: 1687509464492, Body: jsonBody1}, Message{Topic: "A", TS: 1687509464493, Body: jsonBody2}, Message{Topic: "A", TS: 1687509464495, Body: jsonBody4}, Message{Topic: "A", TS: 1687509464496, Body: jsonBody5}}

	result := messagePreprocessing(messages)
	expected := 5

	if len(result) != expected {
		t.Errorf("Add(2, 3) returned %d, expected %d", len(result), expected)
	} else {
		t.Logf("Success %d", len(result))
	}
}
