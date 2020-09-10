package ratelimit

import (
	"net/url"
	"strings"
)

func MajorRoute(path string) bool {
	for _, prefix := range []string{"guilds", "channels", "webhooks"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func HashURL(method string, u *url.URL) string {
	matches := regexpURLSnowflakes.FindAllString(u.Path, -1)
	isMajor := MajorRoute(u.Path)
	buffer := u.Path
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
			for i, r := range suffix {
				if r == '/' {
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
	return method + ":" + buffer
}
