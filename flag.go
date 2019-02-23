package disgord

//go:generate stringer -type=Flag

type Flag uint32

func (f Flag) Ignorecache() bool {
	return (f & DisableCache) > 0
}

const (
	DisableCache Flag = 1 << iota
)

func mergeFlags(flags []Flag) (f Flag) {
	for i := range flags {
		f |= flags[i]
	}

	return f
}
