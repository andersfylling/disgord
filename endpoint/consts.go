package endpoint

// TODO use environment variables to set a desired API version
//      https://discordapp.com/api/v{version_number}

const (
	// APIComEncoding data encoding when communicating with the discord API
	APIComEncoding string = "json"

	// APIVersion desired API version to use
	APIVersion int = 6
)
