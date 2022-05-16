package util

// Function to shorten a string to specified length
func ShortenString(str string, length int) string {
	if len(str) <= length {
		return str
	}

	return str[:length]
}
