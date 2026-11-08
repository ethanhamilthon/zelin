package provider

import "encoding/json"

func mergeDetails(parts []json.RawMessage) json.RawMessage {
	var merged []map[string]any
	slot := map[float64]int{}
	for _, part := range parts {
		var items []map[string]any
		if json.Unmarshal(part, &items) != nil {
			continue
		}
		for _, item := range items {
			index, ok := item["index"].(float64)
			at, seen := slot[index]
			if !ok || !seen {
				if ok {
					slot[index] = len(merged)
				}
				merged = append(merged, item)
				continue
			}
			for key, value := range item {
				text, isText := value.(string)
				prev, hadText := merged[at][key].(string)
				if isText && hadText && (key == "text" || key == "summary" || key == "data") {
					merged[at][key] = prev + text
				} else if value != nil {
					merged[at][key] = value
				}
			}
		}
	}
	if len(merged) == 0 {
		return nil
	}
	data, _ := json.Marshal(merged)
	return data
}
