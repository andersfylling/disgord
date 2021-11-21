package httd

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var regexpURLSnowflakes = regexp.MustCompile(RegexpURLSnowflakes)
var regexpURLReactionPrefix = regexp.MustCompile(`\/channels\/` + RegexpSnowflakes + `\/messages\/\{id\}\/reactions\/`)
var regexpURLReactionEmoji = regexp.MustCompile(`\/channels\/[0-9]+\/messages\/\{id\}\/reactions\/` + RegexpEmoji + `\/?`)
var regexpURLReactionEmojiSegment = regexp.MustCompile(`\/reactions\/` + RegexpEmoji)
// https://discord.com/developers/docs/topics/threads#enumerating-threads
var regexpURLGuildThreadsActive = regexp.MustCompile(`\/channels\/` + RegexpSnowflakes + `\/threads\/active`)
var regexpURLCurrentUserThreadsArchivedPrivate = regexp.MustCompile(`\/channels\/` + RegexpSnowflakes + `\/users\/@me\/threads\/archived\/private`)
var regexpURLChannelThreadsArchivedPublic = regexp.MustCompile(`\/channels\/` + RegexpSnowflakes + `\/threads\/archived\/public`)
var regexpURLChannelThreadsArchivedPrivate = regexp.MustCompile(`\/channels\/` + RegexpSnowflakes + `\/threads\/archived\/private`)


// Request is populated before executing a Discord request to correctly generate a http request
type Request struct {
	Ctx context.Context

	Method      string
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
		r.Method = http.MethodGet
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
	for _, prefix := range []string{"/guilds", "/channels", "/webhooks", "/interactions"} {
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

	// check for any thread endpoints
	guildThreadMatch := regexpURLGuildThreadsActive.FindAllString(buffer, -1)
	currentUserThreadsMatch := regexpURLCurrentUserThreadsArchivedPrivate.FindAllString(buffer, -1)
	chanThreadsPublicMatch := regexpURLChannelThreadsArchivedPublic.FindAllString(buffer, -1)
	chanThreadsPrivateMatch := regexpURLChannelThreadsArchivedPrivate.FindAllString(buffer, -1)
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
			for i, runeVal := range suffix {
				if runeVal == '/' {
					until = i
					break
				}
			}
			newSuffix := "{emoji}" + suffix[until:]
			buffer = buffer[:len(buffer)-len(suffix)] + newSuffix
		}
	} else if guildThreadMatch != nil {
		buffer += "threads/active"
	} else if currentUserThreadsMatch != nil {
		buffer += "users/@me/threads/archived/private"
	} else if chanThreadsPublicMatch != nil {
		buffer += "threads/archived/public"
	} else if chanThreadsPrivateMatch != nil {
		buffer += "threads/archived/private"
	}

	if strings.HasSuffix(buffer, "/") {
		buffer = buffer[:len(buffer)-1]
	}
	return r.Method + ":" + buffer
}
