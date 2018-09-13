package disgord

// common functionality/types used by struct_*.go files goes here

// Copier holds the CopyOverTo method which copies all it's content from one
// struct to another. Note that this requires a deep copy.
// useful when overwriting already existing content in the cache to reduce GC.
type Copier interface {
	CopyOverTo(other interface{}) error
}

func NewErrorUnsupportedType(message string) *ErrorUnsupportedType {
	return &ErrorUnsupportedType{
		info: message,
	}
}

type ErrorUnsupportedType struct {
	info string
}

func (eut *ErrorUnsupportedType) Error() string {
	return eut.info
}

// DiscordUpdater holds the Update method for updating any given Discord struct
// (fetch the latest content). If you only want to keep up to date with the
// cache use the UpdateFromCache method.
// TODO: change param type for UpdateFromCache once caching is implemented
//type DiscordUpdater interface {
//	Update(session Session)
//	UpdateFromCache(session Session)
//}

// DiscordSaver holds the SaveToDiscord method for sending changes to the
// Discord API over REST.
// If you change any of the values and want to notify Discord about your change,
// use the Save method to send a REST request (assuming that the struct values
// can be updated).
//
// NOTE! if the struct has an snowflake/ID, it will update content. But if the
// snowflake is missing/not set, it will create content (if possible,
// otherwise you will get an error)
type discordSaver interface {
	saveToDiscord(session Session) error
}

// DiscordDeleter holds the DeleteFromDiscord method which deletes a given
// object from the Discord servers.
type discordDeleter interface {
	deleteFromDiscord(session Session) error
}

// DeepCopier holds the DeepCopy method which creates and returns a deep copy of
// any struct.
type DeepCopier interface {
	DeepCopy() interface{}
}
