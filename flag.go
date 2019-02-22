package disgord

//go:generate stringer -type=Flag

type Flag uint32

func (f Flag) Ignorecache() bool {
	return (f & DisableCache) > 0
}

const (
	DisableCache Flag = iota
)
