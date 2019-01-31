// +build !disgord_removeDiscordMutex

package disgord

import "sync"

type Lockable struct {
	sync.RWMutex
}
