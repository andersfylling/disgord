package discordws

import "net/http"

type Config struct {
	Token        string
	HTTPClient   *http.Client
	DAPIVersion  int
	DAPIEncoding string
	Debug        bool
}
