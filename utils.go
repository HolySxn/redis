package main

func isValidPort(port int) bool {
	if port < 1024 || port > 49151 {
		return false
	}

	return true
}
