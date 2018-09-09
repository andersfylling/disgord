package websocket

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/andersfylling/disgord/endpoint"
)

func getGatewayRoute(client *http.Client, version int, encoding string) (url string, err error) {
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

	gatewayResponse := &getGatewayResponse{}
	err = unmarshal([]byte(body), gatewayResponse) // redundant?
	if err != nil {
		return
	}

	url = gatewayResponse.URL + "?v=" + strconv.Itoa(version) + "&encoding=" + encoding
	return
}
