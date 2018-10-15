package websocket

import (
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/httd"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	encodingJSON = "json"
)

type gatewayResponse struct {
	URL string `json:"url"`
}

// getGatewayRoute get the connection endpoint for the session
func getGatewayRoute(client *http.Client, version int) (url string, err error) {
	var resp *http.Response
	resp, err = client.Get(endpoint.Gateway(version))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	gatewayResponse := gatewayResponse{}
	err = httd.Unmarshal(body, &gatewayResponse)
	if err != nil {
		return
	}

	url = gatewayResponse.URL + "?v=" + strconv.Itoa(version) + "&encoding=" + encodingJSON
	return
}
