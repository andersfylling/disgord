// Disgord
//
package disgord

const (
	JSONEncoding = "JSON"

	// APIVersion desired API version to use
	APIVersion        = 6 // February 5, 2018
	DefaultAPIVersion = 6

	GitHubURL = "https://github.com/andersfylling/disgord"
	Version   = "v0.3.1"
)

func LibraryInfo() string {
	return "Disgord " + Version
}
