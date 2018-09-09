// Package disgord GoLang module for interacting with the Discord API
package disgord

const (
	// JSONEncoding const for JSON encoding type
	JSONEncoding = "JSON"

	// APIVersion desired API version to use
	APIVersion = 6 // February 5, 2018
	// DefaultAPIVersion the default Discord API version
	DefaultAPIVersion = 6

	// GitHubURL complete url for this project
	GitHubURL = "https://github.com/andersfylling/disgord"

	// Version project version
	Version = "v0.3.1"
)

// LibraryInfo returns name + version
func LibraryInfo() string {
	return "Disgord " + Version
}
