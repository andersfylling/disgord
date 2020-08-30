package httd

import (
	"context"
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
var regexpURLReactionPrefix = regexp.MustCompile(`\/channels\/` + RegexpSnowflakes + `\/messages\/\{id\}\/reactions\/`)
var regexpURLReactionEmoji = regexp.MustCompile(`\/channels\/[0-9]+\/messages\/\{id\}\/reactions\/` + RegexpEmoji + `\/?`)
var regexpURLReactionEmojiSegment = regexp.MustCompile(`\/reactions\/` + RegexpEmoji)

// Request is populated before executing a Discord request to correctly generate a http request
type Request struct {
	Ctx context.Context

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
	if r.Ctx == nil {
		r.Ctx = context.Background()
	}
	// too much magic
	// if c.Body != nil && c.ContentType == "" {
	// 	c.ContentType = ContentTypeJSON
	// }

	r.hashedEndpoint = r.HashEndpoint()
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

	// check for reaction endpoints, convert emoji identifier to {emoji}
	reactionPrefixMatch := regexpURLReactionPrefix.FindAllString(buffer, -1)
	if reactionPrefixMatch != nil {
		if regexpURLReactionEmoji.FindAllString(buffer, -1) != nil {
			reactionEmojis := regexpURLReactionEmojiSegment.FindAllString(buffer, -1)
			for i := range reactionEmojis {
				buffer = strings.ReplaceAll(buffer, reactionEmojis[i], "/reactions/{emoji}")
			}
		} else {
			// corner case for urls with emojis
			suffix := buffer[len(reactionPrefixMatch[0]):]
			until := len(suffix)
			for i, rune := range suffix {
				if rune == '/' {
					until = i
					break
				}
			}
			newSuffix := "{emoji}" + suffix[until:]
			buffer = buffer[:len(buffer)-len(suffix)] + newSuffix
		}
	}

	if strings.HasSuffix(buffer, "/") {
		buffer = buffer[:len(buffer)-1]
	}
	return r.Method.String() + ":" + buffer
}
