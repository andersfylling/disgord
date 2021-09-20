package disgord

type Flag uint32

func (f Flag) Ignorecache() bool {
	return (f & IgnoreCache) > 0
}

func (f Flag) IgnoreEmptyParams() bool {
	return (f & IgnoreEmptyParams) > 0
}

// Deprecated
func (f Flag) Sort() bool {
	flags := SortByID | SortByName
	flags |= OrderAscending | OrderDescending

	return (f & flags) > 0
}

const (
	IgnoreCache Flag = 1 << iota
	IgnoreEmptyParams

	// Deprecated: use disgordutil.SortByID instead
	SortByID

	// Deprecated: use disgordutil.SortByName instead
	SortByName

	// Deprecated: use disgordutil.SortByHoist instead
	SortByHoist

	// Deprecated: use disgordutil.SortByGuildID instead
	SortByGuildID

	// Deprecated: use disgordutil.SortByChannelID instead
	SortByChannelID


	// Deprecated: use disgordutil.OrderAscending instead
	OrderAscending // default when sorting

	// Deprecated: use disgordutil.OrderDescending instead
	OrderDescending
)

func mergeFlags(flags []Flag) (f Flag) {
	for i := range flags {
		f |= flags[i]
	}

	return f
}

func ignoreCache(flags ...Flag) bool {
	return (mergeFlags(flags) & IgnoreCache) > 0
}
