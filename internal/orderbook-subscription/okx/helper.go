package okx

import "encoding/json"

func isValidOrderBookMessage(message []byte) bool {
	if len(message) == 0 {
		return false
	}

	var genericJSON map[string]interface{}
	err := json.Unmarshal(message, &genericJSON)
	if err != nil {
		return false
	}

	if arg, ok := genericJSON["arg"].(map[string]interface{}); ok {
		if channel, ok := arg["channel"].(string); ok && channel == "books5" {
			if _, ok := genericJSON["data"].([]interface{}); ok {
				return true
			}
		}
	}

	return false
}
