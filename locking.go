// +build !removeDiscordMutex

package disgord

import "sync"

type Lockable struct {
	sync.RWMutex
}
