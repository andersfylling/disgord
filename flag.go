package disgord

//go:generate stringer -type=Flag
type Flag uint32

// IgnoreCache is to signify that we don't want to extract a certain
// data point from the cache - if we interact with it.
func (f Flag) Ignorecache() bool {
	return (f & IgnoreCache) > 0
}

func (f Flag) IgnoreEmptyParams() bool {
	return (f & IgnoreEmptyParams) > 0
}

func (f Flag) Sort() bool {
	flags := SortByID | SortByName
	flags |= OrderAscending | OrderDescending

	return (f & flags) > 0
}

const (
	IgnoreCache Flag = 1 << iota
	IgnoreEmptyParams

	// sort options
	SortByID
	SortByName
	SortByHoist
	SortByGuildID
	SortByChannelID

	// ordering
	OrderAscending // default when sorting
	OrderDescending
)

func mergeFlags(flags []Flag) (f Flag) {
	for i := range flags {
		f |= flags[i]
	}

	return f
}
