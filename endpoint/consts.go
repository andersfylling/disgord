package endpoint

// TODO use environment variables to set a desired API version
//      https://discordapp.com/api/v{version_number}

const (
	// APIComEncoding data encoding when communicating with the discord API
	APIComEncoding string = "json"

	// APIVersion desired API version to use
	APIVersion string = "6"

	// BaseURL The base URL for all API requests
	BaseURL string = "https://discordapp.com/api"

	// BaseURLWithVersion uses the Discord API version specified
	BaseURLWithVersion string = BaseURL + "/v" + APIVersion

	// Gateway returns a object containing a WSS url, and does not require auth.
	Gateway string = BaseURLWithVersion + "/gateway"

	// GatewayBot returns WSS url and shard counts. Requires authentication.
	GatewayBot string = Gateway + "/bot"
)
