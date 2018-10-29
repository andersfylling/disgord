// +build removeDiscordMutex

package disgord

// Lockable is removed on compile time since it holds no content. This allows the removal of mutexes if desired by the
// developer. Might improve memory usage for larger bots.
type Lockable struct{}

func (l *Lockable) RLock()   {}
func (l *Lockable) RUnlock() {}
func (l *Lockable) Lock()    {}
func (l *Lockable) Unlock()  {}
