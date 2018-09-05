package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/disgord/rest/httd"
	"github.com/json-iterator/go"
)

type URLParameters interface {
	GetQueryString() string
}

// unmarshalJSONIterator https://github.com/json-iterator/go
func unmarshalJSONIterator(data []byte, v interface{}) (err error) {
	err = jsoniter.Unmarshal(data, v)
	return
}

// unmarshalSTD standard GoLang implementation
func unmarshalSTD(data []byte, v interface{}) (err error) {
	err = json.Unmarshal(data, v)
	return
}

func unmarshal(data []byte, v interface{}) error {
	return unmarshalJSONIterator(data, v)
}

// testing
func createTestRequester() (*httd.Client, error) {
	reqConf := &httd.Config{
		APIVersion:         6,
		BotToken:           os.Getenv(constant.DisgordTestBot),
		UserAgentSourceURL: constant.GitHubURL,
		UserAgentVersion:   constant.Version,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		CancelRequestWhenRateLimited: false,
	}
	if reqConf.BotToken == "" {
		return nil, errors.New("missing bot token")
	}
	return httd.NewClient(reqConf), nil
}
