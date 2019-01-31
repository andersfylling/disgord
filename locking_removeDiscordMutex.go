// +build disgord_removeDiscordMutex

package disgord

// Lockable is removed on compile time since it holds no content. This allows the removal of mutexes if desired by the
// developer.
type Lockable struct{}

func (l *Lockable) RLock()   {}
func (l *Lockable) RUnlock() {}
func (l *Lockable) Lock()    {}
func (l *Lockable) Unlock()  {}
