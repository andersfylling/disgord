package httd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type httpMethod string

var _ fmt.Stringer = (*httpMethod)(nil)

func (method httpMethod) String() string {
	return string(method)
}

const (
	MethodGet    httpMethod = http.MethodGet
	MethodDelete httpMethod = http.MethodDelete
	MethodPost   httpMethod = http.MethodPost
	MethodPatch  httpMethod = http.MethodPatch
	MethodPut    httpMethod = http.MethodPut
)

var regexpURLSnowflakes = regexp.MustCompile(RegexpURLSnowflakes)

// Request is populated before executing a Discord request to correctly generate a http request
type Request struct {
	Method      httpMethod
	Endpoint    string
	Body        interface{} // will automatically marshal to JSON if the ContentType is httd.ContentTypeJSON
	ContentType string

	// Reason is a X-Audit-Log-Reason header field that will show up on the audit log for this action.
	Reason string

	bodyReader     io.Reader
	hashedEndpoint string
}

func (r *Request) PopulateMissing() {
	if r.Method == "" {
		r.Method = MethodGet
	}
	// too much magic
	// if c.Body != nil && c.ContentType == "" {
	// 	c.ContentType = ContentTypeJSON
	// }

	r.hashedEndpoint = r.HashEndpoint()
}

func (r *Request) init() (err error) {
	r.PopulateMissing()
	if r.Body != nil && r.bodyReader == nil {
		switch b := r.Body.(type) { // Determine the type of the passed body so we can treat it differently
		case io.Reader:
			r.bodyReader = b
		default:
			// If the type is unknown, possibly Marshal it as JSON
			if r.ContentType != ContentTypeJSON {
				return errors.New("unknown request body types and only be used in conjunction with httd.ContentTypeJSON")
			}

			if r.bodyReader, err = convertStructToIOReader(r.Body); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Request) HashEndpoint() string {
	endpoint := strings.Split(r.Endpoint, "?")[0]
	matches := regexpURLSnowflakes.FindAllString(endpoint, -1)

	var isMajor bool
	for _, prefix := range []string{"/guilds", "/channels", "/webhooks"} {
		if strings.HasPrefix(endpoint, prefix) {
			isMajor = true
			break
		}
	}

	buffer := endpoint
	for i := range matches {
		if i == 0 && isMajor {
			continue
		}

		buffer = strings.ReplaceAll(buffer, matches[i], "/{id}/")
	}

	if strings.HasSuffix(buffer, "/") {
		buffer = buffer[:len(buffer)-1]
	}
	return r.Method.String() + ":" + buffer
}
