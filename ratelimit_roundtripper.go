package disgord
//
// import (
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"sync"
//
// 	"golang.org/x/time/rate"
// )
//
// type DiscordURL url.URL
//
// func (d *DiscordURL) IsMajor() bool {
// 	for _, prefix := range []string{"guilds", "channels", "webhooks"} {
// 		if strings.HasPrefix(d.Path, prefix) {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func (d *DiscordURL) Hash(method string) string {
// 	matches := regexpURLSnowflakes.FindAllString(d.Path, -1)
// 	isMajor := d.IsMajor()
// 	buffer := d.Path
// 	for i := range matches {
// 		if i == 0 && isMajor {
// 			continue
// 		}
//
// 		buffer = strings.ReplaceAll(buffer, matches[i], "/{id}/")
// 	}
//
// 	// check for reaction endpoints, convert emoji identifier to {emoji}
// 	reactionPrefixMatch := regexpURLReactionPrefix.FindAllString(buffer, -1)
// 	if reactionPrefixMatch != nil {
// 		if regexpURLReactionEmoji.FindAllString(buffer, -1) != nil {
// 			reactionEmojis := regexpURLReactionEmojiSegment.FindAllString(buffer, -1)
// 			for i := range reactionEmojis {
// 				buffer = strings.ReplaceAll(buffer, reactionEmojis[i], "/reactions/{emoji}")
// 			}
// 		} else {
// 			// corner case for urls with emojis
// 			suffix := buffer[len(reactionPrefixMatch[0]):]
// 			until := len(suffix)
// 			for i, r := range suffix {
// 				if r == '/' {
// 					until = i
// 					break
// 				}
// 			}
// 			newSuffix := "{emoji}" + suffix[until:]
// 			buffer = buffer[:len(buffer)-len(suffix)] + newSuffix
// 		}
// 	}
//
// 	if strings.HasSuffix(buffer, "/") {
// 		buffer = buffer[:len(buffer)-1]
// 	}
// 	return method + ":" + buffer
// }
//
// type reqWaiterChan chan *rate.Limiter
//
// type RateLimit struct {
// 	Next http.RoundTripper
//
// 	vtableMu sync.RWMutex
// 	vtable   map[string]string
//
// 	waitersMu sync.RWMutex
// 	waiters   map[string]reqWaiterChan
//
// 	bucketsMy sync.RWMutex
// 	buckets   map[string]*rate.Limiter
// }
//
// var _ http.RoundTripper = &RateLimit{}
//
// func (r *RateLimit) RoundTrip(req *http.Request) (resp *http.Response, err error) {
// 	if !r.isDiscordAPIRequest(req) {
// 		return r.Next.RoundTrip(req)
// 	}
//
// 	durl := DiscordURL(*req.URL)
// 	localHash := durl.Hash(req.Method)
// 	discordHash := r.discordHash(localHash)
// 	if discordHash == "" {
// 		waitChan := r.waiter(localHash, func() {
// 			resp, err := r.Next.RoundTrip(req)
// 			if err != nil {
// 				return nil, err
// 			}
// 		})
// 	}
//
// 	return r.rateLimit(req, func() (*http.Response, error) {
// 		return r.Next.RoundTrip(req)
// 	})
// }
//
// func (r *RateLimit) isDiscordAPIRequest(req *http.Request) bool {
// 	const DiscordAPIURLPrefix = "https://discord.com/api/v"
// 	return strings.HasPrefix(req.URL.String(), DiscordAPIURLPrefix)
// }
//
// func (r *RateLimit) discordHash(localHash string) string {
// 	r.vtableMu.RLock()
// 	defer r.vtableMu.RUnlock()
//
// 	if discordHash, ok := r.vtable[localHash]; ok {
// 		return discordHash
// 	} else {
// 		return ""
// 	}
// }
//
// func (r *RateLimit) waiter(localHash string) reqWaiterChan {
// 	r.waitersMu.RLock()
// 	defer r.waitersMu.RUnlock()
//
// 	if waiter, ok := r.waiters[localHash]; ok {
// 		return waiter
// 	} else {
// 		return nil
// 	}
// }
//
// func (r *RateLimit) bucket(hash string, setupBucket func() (discordHash string, bucket *rate.Limiter)) string {
// 	r.vtableMu.RLock()
// 	if discordHash, ok := r.vtable[hash]; ok {
// 		r.vtableMu.RUnlock()
// 		return discordHash
// 	}
// 	r.vtableMu.RUnlock()
//
// 	r.vtableMu.Lock()
// 	// check if another request got here faster
// 	if discordHash, ok := r.vtable[hash]; ok {
// 		r.vtableMu.Unlock()
// 		return discordHash
// 	}
//
// 	discordHash, bucket := setupBucket()
// 	r.vtable[hash] = discordHash
//
// 	// this should be updated once bucket information is delivered by discord
// 	r.vtable[hash] = hash
// 	return hash
// }
//
// func (r *RateLimit) bucket(hash string) *rate.Limiter {
//
// }
