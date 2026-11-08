package app

func failure(reason string) string {
	if reason == "Stopped" {
		return "Stopped by user."
	}
	return "Error: " + reason
}
