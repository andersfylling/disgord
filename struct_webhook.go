package disgord

import "sync"

// Webhook Used to represent a webhook
// https://discordapp.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	sync.RWMutex `json:"-"`

	ID        Snowflake `json:"id"`                 //  |
	GuildID   Snowflake `json:"guild_id,omitempty"` //  |?
	ChannelID Snowflake `json:"channel_id"`         //  |
	User      *User     `json:"user,omitempty"`     // ?|
	Name      string    `json:"name"`               //  |?
	Avatar    string    `json:"avatar"`             //  |?
	Token     string    `json:"token"`              //  |
}

// DeepCopy see interface at struct.go#DeepCopier
func (w *Webhook) DeepCopy() (copy interface{}) {
	copy = &Webhook{}
	w.CopyOverTo(copy)

	return
}

// CopyOverTo see interface at struct.go#Copier
func (w *Webhook) CopyOverTo(other interface{}) (err error) {
	var ok bool
	var hook *Webhook
	if hook, ok = other.(*Webhook); !ok {
		err = newErrorUnsupportedType("given interface{} was not of type *Webhook")
		return
	}

	w.RLock()
	hook.Lock()

	hook.ID = w.ID
	hook.GuildID = w.GuildID
	hook.ChannelID = w.ChannelID
	hook.User = w.User.DeepCopy().(*User)
	hook.Name = w.Name
	hook.Avatar = w.Avatar
	hook.Token = w.Token

	w.RUnlock()
	hook.Unlock()
	return
}
